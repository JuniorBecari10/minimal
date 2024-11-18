package compiler

import (
	"bytes"
	"fmt"
	"vm-go/ast"
	"vm-go/util"
)

func (c *Compiler) statement(stmt ast.Statement) []byte {
	res := bytes.Buffer{}

	switch s := stmt.(type) {
		case ast.IfStatement: {
			res.WriteString(string(c.expression(s.Condition)))
			then := c.statements(s.Then.Stmts)

			res.WriteByte(OP_JUMP_FALSE)
			res.WriteString(string(util.IntToBytes(len(then))))
			res.WriteByte(OP_POP)

			res.WriteString(string(then))

			if s.Else != nil {
				else_ := c.statements(s.Else.Stmts)

				res.WriteByte(OP_JUMP)
				res.WriteString(string(util.IntToBytes(len(else_))))
				res.WriteByte(OP_POP)

				res.WriteString(string(else_))
			}
		}

		case ast.WhileStatement: {
			res.WriteString(string(c.expression(s.Condition)))
			block := c.statements(s.Block.Stmts)

			res.WriteByte(OP_LOOP_FALSE)
			res.WriteByte(OP_POP)
			res.WriteString(string(util.IntToBytes(len(block))))
		}

		case ast.VarStatement: {
			// Find the variable to check if it already exists or not
			index := -1
			for i := len(c.variables) - 1; i >= 0; i-- {
				if c.variables[i].name.Lexeme == s.Name.Lexeme {
					index = i
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
					// it's a redeclaration
					c.error(s.Pos, fmt.Sprintf("'%s' has already been declared in this scope", s.Name.Lexeme))
					return res.Bytes()
				}
			}

			res.WriteString(string(c.expression(s.Init)))

			if c.hadError {
				return res.Bytes()
			}

			// pop from stack and push to variable stack
			res.WriteByte(OP_DEF_VAR)
		}

		case ast.BlockStatement: {
			c.beginScope()
			res.WriteString(string(c.statements(s.Stmts)))
			res.WriteString(string(c.endScope()))
		}

		case ast.PrintStatement: {
			res.WriteString(string(c.expression(s.Expr)))

			if c.hadError {
				return res.Bytes()
			}

			res.WriteByte(OP_PRINT)
		}

		case ast.ExprStatement:
			res.WriteString(string(c.expression(s.Expr)))
	}

	return res.Bytes()
}
