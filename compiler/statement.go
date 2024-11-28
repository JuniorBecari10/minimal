package compiler

import (
	"fmt"
	"vm-go/ast"
	"vm-go/util"
)

func (c *Compiler) statement(stmt ast.Statement) {
	switch s := stmt.(type) {
		case ast.IfStatement: {
			c.expression(s.Condition)
			c.writeBytePos(OP_JUMP_FALSE, s.Pos)
			c.writeBytes(util.IntToBytes(0)) // dummy

			jumpOffsetIndex := len(c.chunk.Code)
			c.statements(s.Then.Stmts)
			
			offset := 1 // OP_POP (next instruction)

			if s.Else != nil {
				offset = 6 // OP_POP + OP_JUMP (amount: 4 bytes)
			}

			// insert the real position into the instruction
			amount := len(c.chunk.Code) - jumpOffsetIndex + 4 + offset
			c.backpatch(jumpOffsetIndex, util.IntToBytes(amount))

			c.writeBytePos(OP_POP, s.Pos)

			if s.Else != nil {
				c.writeBytePos(OP_JUMP, s.Pos)
				c.writeBytes(util.IntToBytes(0)) // dummy
				jumpOffsetIndex := len(c.chunk.Code)

				c.writeBytePos(OP_POP, s.Pos)
				c.statements(s.Else.Stmts)

				amount := len(c.chunk.Code) - jumpOffsetIndex + 5 // index + OP_POP
				c.backpatch(jumpOffsetIndex, util.IntToBytes(amount))
			}
		}

		case ast.WhileStatement: {
			c.expression(s.Condition)

			c.writeBytePos(OP_JUMP_FALSE, s.Pos)
			c.writeBytes(util.IntToBytes(0)) // dummy
			c.writeBytePos(OP_POP, s.Pos)

			c.statements(s.Block.Stmts)
			res.WriteString(string(util.IntToBytes(len(block) + 6))) // OP_POP + OP_LOOP_FALSE (amount: 4 bytes)

			c.positions = append(c.positions, s.Pos)
			res.WriteByte(OP_LOOP)
			res.WriteString(string(util.IntToBytes(len(block) + len(condition) + 11))) // own (LOOP_FALSE) (amount: 4 bytes) + block + OP_POP + JUMP_FALSE + condition

			c.positions = append(c.positions, s.Pos)
			res.WriteByte(OP_POP)
		}

		case ast.VarStatement: {
			// Find the variable to check if it already exists or not, in this scope
			index := -1
			for i := len(c.variables) - 1; i >= 0; i-- {
				if c.variables[i].name.Lexeme == s.Name.Lexeme {
					index = i
					break
				}

				if c.variables[i].depth < c.scopeDepth {
					break
				}
			}

			// didn't find it
			if index == -1 {
				c.variables = append(c.variables, Variable{
					name:        s.Name,
					depth:       c.scopeDepth,
					initialized: true,
				})
				// found it
			} else {
				// if it's a global and we've reached its declaration. also make sure that it isn't a redeclaration
				// we do that by checking if the initialized field is true
				if c.variables[index].depth == 0 && c.scopeDepth == 0 && !c.variables[index].initialized {
					c.variables[index].initialized = true
				} else {
					// Found a variable with the same name, it can be in the same scope or not
					existing := c.variables[index]

					if existing.depth == 0 && c.scopeDepth == 0 && !existing.initialized {
						// If it's a global variable and uninitialized, allow redeclaration
						c.variables[index].initialized = true
					} else if existing.depth == c.scopeDepth {
						// Redeclaration in the same scope is not allowed
						c.error(s.Pos, fmt.Sprintf("'%s' has already been declared in this scope", s.Name.Lexeme))
						return
					} else {
						// the variable is in an enclosing scope, we'll shadow it
						c.variables = append(c.variables, Variable{
							name:        s.Name,
							depth:       c.scopeDepth,
							initialized: true,
						})
					}
				}
			}

			c.expression(s.Init)

			if c.hadError {
				return
			}

			// pop from stack and push to variable stack
			c.writeBytePos(OP_DEF_VAR, s.Pos)
		}

		case ast.BlockStatement: {
			c.beginScope()
			c.statements(s.Stmts)
			c.writeBytes(c.endScope(s.Pos))
		}

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
