package compiler

import (
	"vm-go/ast"
	"vm-go/util"
	"vm-go/value"
)

func (c *Compiler) statement(stmt ast.Statement) {
	switch s := stmt.(type) {
		case ast.FnStatement: {
			fnCompiler := NewCompiler(s.Body.Stmts, c.fileData)
			fnChunk, hadError := fnCompiler.compileBody()

			if hadError {
				c.hadError = true
				return
			}

			function := value.ValueFunction{
				Arity: len(s.Parameters),
				Chunk: fnChunk,
				Name: &s.Name.Lexeme,
			}

			index := c.addConstant(function)
			c.writeBytePos(OP_PUSH_CONST, s.Pos)
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.IfStatement: {
			c.expression(s.Condition)
			c.writeBytePos(OP_JUMP_FALSE, s.Pos)

			jumpFalseOffsetIndex := len(c.chunk.Code)
			c.writeBytes(util.IntToBytes(0)) // dummy
			c.writeBytePos(OP_POP, s.Pos)

			c.block(s.Then.Stmts, s.Pos)

			if s.Else != nil {
				c.writeBytePos(OP_JUMP, s.Pos)
				jumpOffsetIndex := len(c.chunk.Code)
				c.writeBytes(util.IntToBytes(0)) // dummy

				// insert the real offset into the instruction, right before OP_POP
				c.backpatch(jumpFalseOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpFalseOffsetIndex - 4)) // index

				c.writeBytePos(OP_POP, s.Pos)
				c.block(s.Else.Stmts, s.Pos)
 
				c.backpatch(jumpOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpOffsetIndex - 4)) // index
			} else {
				// insert the real offset into the instruction, if there's no else
				c.backpatch(jumpFalseOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpFalseOffsetIndex - 4)) // index
			}
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
			
			c.writeBytes(c.endScope(s.Pos))
		}

		case ast.VarStatement: {
			c.addVariable(s.Name, s.Pos)
			c.expression(s.Init)

			if c.hadError {
				return
			}

			// pop from stack and push to variable stack
			c.writeBytePos(OP_DEF_VAR, s.Pos)
		}

		case ast.BlockStatement:
			c.block(s.Stmts, s.Pos)

		case ast.PrintStatement: {
			c.expression(s.Expr)

			if c.hadError {
				return
			}

			c.writeBytePos(OP_PRINT, s.Pos)
		}

		case ast.ExprStatement: {
			c.expression(s.Expr)
			c.writeBytePos(OP_POP, s.Pos)
		}
	}
}
