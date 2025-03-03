package compiler

import (
	"vm-go/ast"
	"vm-go/util"
	"vm-go/value"
)

func (c *Compiler) statement(stmt ast.Statement) {
	switch s := stmt.Data.(type) {
		case ast.RecordStatement: {
			fieldNames := make([]string, 0, len(s.Fields))

			for _, field := range s.Fields {
				fieldNames = append(fieldNames, field.Name.Lexeme)
			}

			index := c.addConstant(value.ValueRecord{
				Name: s.Name.Lexeme,
				FieldNames: fieldNames,
				Methods: []value.ValueClosure{}, // empty for now
			})

			c.writeBytePos(OP_PUSH_CONST, value.NewMetaLen1(stmt.Base.Pos))
			c.writeBytes(util.IntToBytes(index))

			// stack the methods and add an APPEND_METHODS instruction to put them into the record
			for _, method := range s.Methods {
				c.compileMethod(method.Parameters, method.Body, &method.Name.Lexeme, stmt.Base.Pos) // TODO: add correct position
			}

			if len(s.Methods) > 0 {
				c.writeBytePos(OP_APPEND_METHODS, value.NewMetaLen1(stmt.Base.Pos))
				c.writeBytes(util.IntToBytes(len(s.Methods)))
			}

			c.addVariable(s.Name, s.Name.Pos)
			c.addDeclarationInstruction(stmt.Base.Pos)
		}

		case ast.FnStatement: {
			c.compileFunction(s.Parameters, s.Body, &s.Name.Lexeme, stmt.Base.Pos)

			c.addVariable(s.Name, s.Name.Pos)
			c.addDeclarationInstruction(stmt.Base.Pos)
		}

		case ast.VarStatement: {
			c.expression(s.Init)

			if c.hadError {
				return
			}

			c.addVariable(s.Name, s.Name.Pos)
			c.addDeclarationInstruction(stmt.Base.Pos)
		}

		// Control flow graph in the compileIf function.
		case ast.IfStatement: {
			var else_ *func() = nil

			if s.Else != nil {
				elseFn := func() {
					c.block(s.Else.Stmts, stmt.Base.Pos)
				}

				else_ = &elseFn
			}

			c.compileIf(s.Condition, func() { c.block(s.Then.Stmts, stmt.Base.Pos) }, else_, stmt.Base.Pos)
		}

		/*
            While Loop
            Control Flow:

                [ condition ] <-+
                                |
                                |
            +-- OP_JUMP_FALSE   | <- break/continue point
            |   OP_POP          |
            |                   |
            |   [ body ]        |
            |                   |
            |   OP_LOOP --------+
            +-> OP_POP

            continues...
		*/
		case ast.WhileStatement: {
			conditionPos := len(c.chunk.Code)
			c.expression(s.Condition)

			c.loopFlowPos = append(c.loopFlowPos, len(c.chunk.Code))
			c.writeBytePos(OP_JUMP_FALSE, value.NewMetaLen1(stmt.Base.Pos))
			jumpOffsetIndex := len(c.chunk.Code)
			c.writeBytes(util.IntToBytes(0)) // dummy

			c.writeBytePos(OP_POP, value.NewMetaLen1(stmt.Base.Pos))
			c.block(s.Block.Stmts, stmt.Base.Pos)

			util.PopList(&c.loopFlowPos)
			
			c.writeBytePos(OP_LOOP, value.NewMetaLen1(stmt.Base.Pos))
			c.writeBytes(util.IntToBytes(len(c.chunk.Code) - conditionPos + 4)) // index

			c.backpatch(jumpOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpOffsetIndex - 4)) // index
			c.writeBytePos(OP_POP, value.NewMetaLen1(stmt.Base.Pos))
		}
        
		/*
            For Loop
            Control Flow:

				- begin scope -

                [ iterable ]
                OP_MAKE_ITERATOR

                OP_DEF_LOCAL
                OP_GET_LOCAL <index> <-+
                                       |
            +-- OP_JUMP_HAS_NO_NEXT    |
            |   OP_POP                 |
            |                          |
            |   OP_JUMP ------+        |
            +-- OP_JUMP_FALSE |        | <- break/continue point
            |   OP_POP        |        |
            |   OP_JUMP ------+---+    |
            |                 |   |    |
            |   [ body ] <----+   |    |
            |                     V    |
            |   OP_GET_LOCAL <index>   |
            |   - end scope -          |
            |   - begin scope -        |
			|                          |
            |   OP_NEXT                |
            |   OP_DEF_LOCAL           |
            |                          |
            |   OP_LOOP ---------------+
            +-> OP_POP

				- end scope -

            continues...
		*/
        case ast.ForStatement: {
            c.beginScope()
            
            c.expression(s.Iterable)
			c.writeBytePos(OP_MAKE_ITERATOR, value.NewMetaLen1(stmt.Base.Pos))
			
            c.addVariable(s.Variable, s.Variable.Pos)
            c.addDeclarationInstruction(s.Variable.Pos)
            
            getLocalPos := len(c.chunk.Code)
            c.identifierNoExpr(s.Variable, value.ChunkMetadata{
                Position: s.Variable.Pos,
                Length: len(s.Variable.Lexeme),
            })

			c.writeBytePos(OP_JUMP_HAS_NO_NEXT, value.NewMetaLen1(stmt.Base.Pos))
			jumpNoNextOffsetIndex := len(c.chunk.Code)
			c.writeBytes(util.IntToBytes(0)) // dummy
			
            c.writeBytePos(OP_POP, value.NewMetaLen1(stmt.Base.Pos))
			
            c.writeBytePos(OP_JUMP, value.NewMetaLen1(stmt.Base.Pos))
			jump1OffsetIndex := len(c.chunk.Code)
			c.writeBytes(util.IntToBytes(0)) // dummy
			
            c.loopFlowPos = append(c.loopFlowPos, len(c.chunk.Code))

			c.writeBytePos(OP_JUMP_FALSE, value.NewMetaLen1(stmt.Base.Pos))
			jumpFalseOffsetIndex := len(c.chunk.Code)
			c.writeBytes(util.IntToBytes(0)) // dummy
            
            c.writeBytePos(OP_POP, value.NewMetaLen1(stmt.Base.Pos))
			
            c.writeBytePos(OP_JUMP, value.NewMetaLen1(stmt.Base.Pos))
			jump2OffsetIndex := len(c.chunk.Code)
            c.writeBytes(util.IntToBytes(0)) // dummy
			
            c.backpatch(jump1OffsetIndex, util.IntToBytes(len(c.chunk.Code) - jump1OffsetIndex - 4)) // index
            c.block(s.Block.Stmts, stmt.Base.Pos)
            
            util.PopList(&c.loopFlowPos)
            
			c.backpatch(jump2OffsetIndex, util.IntToBytes(len(c.chunk.Code) - jump2OffsetIndex - 4)) // index
            c.identifierNoExpr(s.Variable, value.ChunkMetadata{
                Position: s.Variable.Pos,
                Length: len(s.Variable.Lexeme),
            })
			
            // End the scope to discard the loop variable and close it in an upvalue if it's been captured.
			c.endScope(stmt.Base.Pos)
			c.beginScope()

            c.writeBytePos(OP_NEXT, value.NewMetaLen1(stmt.Base.Pos))

			// Create a new variable for the mutation to occur.
            c.addVariable(s.Variable, s.Variable.Pos)
            c.addDeclarationInstruction(s.Variable.Pos)

			c.writeBytePos(OP_LOOP, value.NewMetaLen1(stmt.Base.Pos))
			c.writeBytes(util.IntToBytes(len(c.chunk.Code) - getLocalPos + 4)) // index

			c.backpatch(jumpFalseOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpFalseOffsetIndex - 4)) // index
			c.backpatch(jumpNoNextOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpNoNextOffsetIndex - 4)) // index
			
            c.writeBytePos(OP_POP, value.NewMetaLen1(stmt.Base.Pos))
			c.endScope(stmt.Base.Pos)
        }

		/*
            For Var Loop
            Control Flow:

				- begin scope -

                [ initializer ]
                [ condition ] <------+
                                     |
            +-- OP_JUMP_FALSE        | <- break/continue point
            |   OP_POP               |
            |                        |
            |   [ body ]             |
            |                        |
            |   OP_GET_LOCAL <index> |
            |   - end scope -        |
            |   - begin scope -      |
			|                        |
			|   [ initializer* ]     | (this version of the initializer uses the saved value on the stack)
            |                        |
            | + [ increment ]        | (generated if increment is set)
            | + OP_POP               |
            |                        |
            |   OP_LOOP -------------+
            +-> OP_POP

				- end scope -

            continues...
		*/
		case ast.ForVarStatement: {
			// The declaration stays inside a new scope.
			c.beginScope()
			c.statement(s.Declaration)

			conditionPos := len(c.chunk.Code)
			c.expression(s.Condition)

			c.loopFlowPos = append(c.loopFlowPos, len(c.chunk.Code))

			c.writeBytePos(OP_JUMP_FALSE, value.NewMetaLen1(stmt.Base.Pos))
			jumpFalseOffsetIndex := len(c.chunk.Code)
			c.writeBytes(util.IntToBytes(0)) // dummy

			// And the block inside another.
			c.writeBytePos(OP_POP, value.NewMetaLen1(stmt.Base.Pos))
			c.block(s.Block.Stmts, stmt.Base.Pos)

			util.PopList(&c.loopFlowPos)

			// Push the old value to the stack to save it for the next iteration.
			c.identifier(s.Declaration.Data.(ast.VarStatement).Name, s.Declaration.Data.(ast.VarStatement).Init)

			// End the scope to discard the loop variable and close it in an upvalue if it's been captured.
			c.endScope(stmt.Base.Pos)
			c.beginScope()

			// Create a new variable for the mutation to occur.
			c.addVariable(s.Declaration.Data.(ast.VarStatement).Name, s.Declaration.Data.(ast.VarStatement).Name.Pos)
			c.addDeclarationInstruction(s.Declaration.Base.Pos)
			
			if s.Increment != nil {
				c.expression(*s.Increment)
				c.writeBytePos(OP_POP, value.NewMetaLen1(stmt.Base.Pos))
			}

			c.writeBytePos(OP_LOOP, value.NewMetaLen1(stmt.Base.Pos))
			c.writeBytes(util.IntToBytes(len(c.chunk.Code) - conditionPos + 4)) // index

			c.backpatch(jumpFalseOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpFalseOffsetIndex - 4)) // index
			c.writeBytePos(OP_POP, value.NewMetaLen1(stmt.Base.Pos))
			
			c.endScope(stmt.Base.Pos)
		}

		/*
            Indefinite Loop
            Control Flow:

            +-- OP_JUMP
            |   OP_JUMP_FALSE --+ <- break/continue point
            |   OP_POP          |
            |                   |
            +-> [ block ] <-+   |
                            |   |
                OP_LOOP ----+   |
                OP_POP <--------+

            continues...
		*/
		case ast.LoopStatement: {
			c.writeBytePos(OP_JUMP, value.NewMetaLen1(stmt.Base.Pos))
			jumpJumpOffsetIndex := len(c.chunk.Code)
			c.writeBytes(util.IntToBytes(0)) // dummy

			c.loopFlowPos = append(c.loopFlowPos, len(c.chunk.Code))
			c.writeBytePos(OP_JUMP_FALSE, value.NewMetaLen1(stmt.Base.Pos))
			jumpEndOffsetIndex := len(c.chunk.Code)
			c.writeBytes(util.IntToBytes(0)) // dummy
			c.writeBytePos(OP_POP, value.NewMetaLen1(stmt.Base.Pos))

			c.backpatch(jumpJumpOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpJumpOffsetIndex - 4)) // index

			loopPos := len(c.chunk.Code)
			c.block(s.Block.Stmts, stmt.Base.Pos)

			c.writeBytePos(OP_LOOP, value.NewMetaLen1(stmt.Base.Pos))
			c.writeBytes(util.IntToBytes(len(c.chunk.Code) - loopPos + 4)) // index

			util.PopList(&c.loopFlowPos)
			c.backpatch(jumpEndOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpEndOffsetIndex - 4)) // index
			c.writeBytePos(OP_POP, value.NewMetaLen1(stmt.Base.Pos))
		}

		case ast.BreakStatement: {
			// We'll jump to OP_JUMP_FALSE, which jumps to the end of the loop.
			// So, to do that, we'll push 'false' onto the stack and jump there,
			// which will cause the instruction to break out of the loop.

			if len(c.loopFlowPos) == 0 {
				c.error(stmt.Base.Pos, len(s.Token.Lexeme), "Cannot use 'break' outside of a loop.")
				return
			}

			c.writeBytePos(OP_PUSH_FALSE, value.NewMetaLen1(stmt.Base.Pos))

			c.writeBytePos(OP_LOOP, value.NewMetaLen1(stmt.Base.Pos))
			c.writeBytes(util.IntToBytes(len(c.chunk.Code) - c.loopFlowPos[len(c.loopFlowPos)-1] + 4)) // index
		}

		case ast.ContinueStatement: {
			// The same with continue, but we'll push 'true', because we want the loop to keep running.

			// this might be broken if multiple loops are nested
			if len(c.loopFlowPos) == 0 {
				c.error(stmt.Base.Pos, len(s.Token.Lexeme), "Cannot use 'continue' outside of a loop.")
				return
			}

			c.writeBytePos(OP_PUSH_TRUE, value.NewMetaLen1(stmt.Base.Pos))

			c.writeBytePos(OP_LOOP, value.NewMetaLen1(stmt.Base.Pos))
			c.writeBytes(util.IntToBytes(len(c.chunk.Code) - c.loopFlowPos[len(c.loopFlowPos)-1] + 4)) // index
		}

		case ast.ReturnStatement: {
			if s.Expression != nil {
				c.expression(*s.Expression)
			} else {
				c.writeBytePos(OP_PUSH_VOID, value.NewMetaLen1(stmt.Base.Pos))
			}

			c.writeBytePos(OP_RETURN, value.NewMetaLen1(stmt.Base.Pos))
		}

		case ast.BlockStatement:
			c.block(s.Stmts, stmt.Base.Pos)

		case ast.ExprStatement: {
			// Optimization to remove nodes that don't have side effects.
			reduced := reduceToSideEffect(s.Expr)

			if reduced != nil {
				c.expression(*reduced)
				c.writeBytePos(OP_POP, value.NewMetaLen1(stmt.Base.Pos))
			}
		}
	}
}
