package compiler

import (
	"vm-go/ast"
	"vm-go/util"
	"vm-go/value"
)

func (c *Compiler) statement(stmt ast.Statement) {
	switch s := stmt.(type) {
		case ast.RecordStatement: {
			fieldNames := []string{}

			for _, field := range s.Fields {
				fieldNames = append(fieldNames, field.Name.Lexeme)
			}

			index := c.addConstant(value.ValueRecord{
				FieldNames: fieldNames,
				Name: s.Name.Lexeme,
			})

			c.writeBytePos(OP_PUSH_CONST, s.Pos)
			c.writeBytes(util.IntToBytes(index))

			c.addVariable(s.Name, s.Name.Pos)
			c.addDeclarationInstruction(s.Pos)
		}

		case ast.FnStatement: {
			c.compileFunction(s.Parameters, s.Body, &s.Name.Lexeme, s.Pos)

			c.addVariable(s.Name, s.Name.Pos)
			c.addDeclarationInstruction(s.Pos)
		}

		case ast.VarStatement: {
			c.expression(s.Init)

			if c.hadError {
				return
			}

			c.addVariable(s.Name, s.Name.Pos)
			c.addDeclarationInstruction(s.Pos)
		}

		case ast.IfStatement: {
			var else_ *func() = nil

			if s.Else != nil {
				elseFn := func() {
					c.block(s.Else.Stmts, s.Pos)
				}

				else_ = &elseFn
			}

			c.compileIf(s.Condition, func() { c.block(s.Then.Stmts, s.Pos) }, else_, s.Pos)
		}

		case ast.WhileStatement: {
			conditionPos := len(c.chunk.Code)
			c.expression(s.Condition)

			c.writeBytePos(OP_JUMP_FALSE, s.Pos)
			jumpOffsetIndex := len(c.chunk.Code)
			c.writeBytes(util.IntToBytes(0)) // dummy

			c.writeBytePos(OP_POP, s.Pos)
			c.block(s.Block.Stmts, s.Pos)

			c.writeBytePos(OP_LOOP, s.Pos)
			c.writeBytes(util.IntToBytes(len(c.chunk.Code) - conditionPos + 4)) // index

			c.backpatch(jumpOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpOffsetIndex - 4)) // index
			c.writeBytePos(OP_POP, s.Pos)
		}

		case ast.ForVarStatement: {
			// the declaration stays inside a new scope
			c.beginScope()
			c.statement(s.Declaration)

			conditionPos := len(c.chunk.Code)
			c.expression(s.Condition)

			c.writeBytePos(OP_JUMP_FALSE, s.Pos)
			jumpFalseOffsetIndex := len(c.chunk.Code)
			c.writeBytes(util.IntToBytes(0)) // dummy

			// and the block inside another
			c.writeBytePos(OP_POP, s.Pos)
			c.block(s.Block.Stmts, s.Pos)
			
			if s.Increment != nil {
				c.expression(*s.Increment)
				c.writeBytePos(OP_POP, s.Pos)
			}

			c.writeBytePos(OP_LOOP, s.Pos)
			c.writeBytes(util.IntToBytes(len(c.chunk.Code) - conditionPos + 4)) // index

			c.backpatch(jumpFalseOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpFalseOffsetIndex - 4)) // index
			c.writeBytePos(OP_POP, s.Pos)
			
			c.endScope(s.Pos)
		}

		case ast.ReturnStatement: {
			if s.Expression != nil {
				c.expression(*s.Expression)
			} else {
				c.writeBytePos(OP_VOID, s.Pos)
			}

			c.endScope(s.Pos)
			c.writeBytePos(OP_RETURN, s.Pos)
		}

		case ast.BlockStatement:
			c.block(s.Stmts, s.Pos)

		case ast.ExprStatement: {
			c.expression(s.Expr)
			c.writeBytePos(OP_POP, s.Pos)
		}
	}
}
