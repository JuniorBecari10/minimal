package compiler

import (
	"fmt"
	"vm-go/ast"
	"vm-go/token"
	"vm-go/util"
	"vm-go/value"
)

func (c *Compiler) expression(expr ast.Expression) {
	if c.hadError {
		return
	}

	switch e := expr.Data.(type) {
		case ast.NumberExpression: {
			index := c.addConstant(value.ValueNumber{ Value: e.Literal })

			c.writeBytePos(OP_PUSH_CONST, value.ChunkMetadata{
				Position: expr.Base.Pos,
				Length: expr.Base.Length,
			})
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.StringExpression: {
			index := c.addConstant(value.ValueString{ Value: e.Literal })

			c.writeBytePos(OP_PUSH_CONST, value.NewMetaLen1(expr.Base.Pos))
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.BoolExpression: {
			if e.Literal {
				c.writeBytePos(OP_PUSH_TRUE, value.ChunkMetadata{
					Position: expr.Base.Pos,
					Length: expr.Base.Length,
				})
			} else {
				c.writeBytePos(OP_PUSH_FALSE, value.ChunkMetadata{
					Position: expr.Base.Pos,
					Length: expr.Base.Length,
				})
			}
		}

		case ast.NilExpression: {
			c.writeBytePos(OP_PUSH_NIL, value.ChunkMetadata{
				Position: expr.Base.Pos,
				Length: expr.Base.Length,
			})
		}

		case ast.VoidExpression: {
			if e.Expr != nil {
				c.expression(*e.Expr)

				c.writeBytePos(OP_POP, value.ChunkMetadata{
					Position: expr.Base.Pos,
					Length: expr.Base.Length,
				})
			}

			c.writeBytePos(OP_PUSH_VOID, value.ChunkMetadata{
				Position: expr.Base.Pos,
				Length: expr.Base.Length,
			})
		}

		case ast.IdentifierExpression:
			c.identifier(e.Token, expr)

		case ast.SelfExpression:
			c.identifier(e.Token, expr)

		/*
            Logical Operators (short-circuit behavior)
            Control Flow

                [ left operand ]

            +-- OP_JUMP_FALSE (and) / OP_JUMP_TRUE (or)
            |   OP_POP
            |
            |   [ right operand ]
            |   OP_ASSERT_BOOL
            |
            +-> continues...
		*/
		case ast.LogicalExpression: {
			if e.ShortCircuit {
				var operation byte

				switch e.Operator.Kind {
					case token.TokenAndKw:
						operation = OP_JUMP_FALSE
					case token.TokenOrKw:
						operation = OP_JUMP_TRUE
					default:
						panic(fmt.Sprintf("Unknown logical operator: '%s'", e.Operator.Lexeme))
				}

				c.expression(e.Left)
				c.writeBytePos(operation, value.ChunkMetadata{
					Position: expr.Base.Pos,
					Length: expr.Base.Length,
				})

				jumpOffsetIndex := len(c.chunk.Code)
				c.writeBytes(util.IntToBytes(0)) // dummy
				c.writeBytePos(OP_POP, value.ChunkMetadata{
					Position: expr.Base.Pos,
					Length: expr.Base.Length,
				})

				c.expression(e.Right)
				c.writeBytePos(OP_ASSERT_BOOL, value.ChunkMetadata{
					Position: expr.Base.Pos,
					Length: expr.Base.Length,
				})

				c.backpatch(jumpOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpOffsetIndex - 4)) // index
			} else {
				c.expression(e.Left)
				c.expression(e.Right)

				c.chunk.Metadata = append(c.chunk.Metadata, value.ChunkMetadata{
					Position: e.Operator.Pos,
					Length: len(e.Operator.Lexeme),
				})

				switch e.Operator.Kind {
					case token.TokenAndKw:
						c.writeByte(OP_AND)
					case token.TokenOrKw:
						c.writeByte(OP_OR)

					default:
						panic(fmt.Sprintf("Unknown logical operator: '%s'", e.Operator.Kind))
				}
			}
		}

		case ast.BinaryExpression: {
			c.expression(e.Left)
			c.expression(e.Right)

			c.chunk.Metadata = append(c.chunk.Metadata, value.ChunkMetadata{
				Position: e.Operator.Pos,
				Length: len(e.Operator.Lexeme),
			})

			switch e.Operator.Kind {
				case token.TokenPlus:
					c.writeByte(OP_ADD)
				case token.TokenMinus:
					c.writeByte(OP_SUB)
				case token.TokenStar:
					c.writeByte(OP_MUL)
				case token.TokenSlash:
					c.writeByte(OP_DIV)
				
				case token.TokenPercent:
					c.writeByte(OP_MODULO)
				
				case token.TokenDoubleEqual:
					c.writeByte(OP_EQUAL)
				case token.TokenBangEqual:
					c.writeByte(OP_NOT_EQUAL)
				
				case token.TokenGreater:
					c.writeByte(OP_GREATER)
				case token.TokenGreaterEqual:
					c.writeByte(OP_GREATER_EQUAL)
				
				case token.TokenLess:
					c.writeByte(OP_LESS)
				case token.TokenLessEqual:
					c.writeByte(OP_LESS_EQUAL)
				
				case token.TokenAndKw:
					c.writeByte(OP_AND)
				case token.TokenOrKw:
					c.writeByte(OP_OR)

				default:
					panic(fmt.Sprintf("Unknown binary operator: '%s'", e.Operator.Kind))
			}
		}

		case ast.UnaryExpression: {
			c.expression(e.Operand)

			c.chunk.Metadata = append(c.chunk.Metadata, value.ChunkMetadata{
				Position: e.Operator.Pos,
				Length: len(e.Operator.Lexeme),
			})

			switch e.Operator.Kind {
				case token.TokenNotKw:
					c.writeByte(OP_NOT)
				case token.TokenMinus:
					c.writeByte(OP_NEGATE)
				
				default:
					panic(fmt.Sprintf("Unknown unary operator: '%s'", e.Operator.Kind))
			}
		}

		case ast.CallExpression: {
			// Generate an optimized call to properties.
			switch e.Callee.Data.(type) {
				case ast.GetPropertyExpression: {
					callee := e.Callee.Data.(ast.GetPropertyExpression)
					c.expression(callee.Left)

					for _, arg := range e.Arguments {
						c.expression(arg)
					}

					// Store the name as a string in the constant table and retrieve it later.
					index := c.addConstant(value.ValueString{ Value: callee.Property.Lexeme })

					c.writeBytePos(OP_CALL_PROPERTY, value.ChunkMetadata{
						Position: callee.Property.Pos,
						Length: len(callee.Property.Lexeme),
					})
					c.writeBytes(util.IntToBytes(index))
					c.writeBytes(util.IntToBytes(len(e.Arguments)))
				}

				default: {
					c.expression(e.Callee)

					for _, arg := range e.Arguments {
						c.expression(arg)
					}

					c.writeBytePos(OP_CALL, value.ChunkMetadata{
						Position: expr.Base.Pos,
						Length: expr.Base.Length,
					})
					c.writeBytes(util.IntToBytes(len(e.Arguments)))
				}
			}
		}

		case ast.GroupExpression:
			c.expression(e.Expr)
		
		case ast.IdentifierAssignmentExpression: {
			index, opcode := c.resolveVariable(e.Name, true)

			if index < 0 {
				return
			}

			c.expression(e.Expr)

			c.writeBytePos(byte(opcode), value.ChunkMetadata{
				Position: expr.Base.Pos,
				Length: expr.Base.Length,
			})
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.FnExpression:
			c.compileFunction(e.Parameters, e.Body, nil, expr.Base.Pos)
		
		case ast.IfExpression: {
			elseFn := func() {
				c.expression(e.Else)
			}
			
			else_ := &elseFn
			c.compileIf(e.Condition, func() { c.expression(e.Then) }, else_, expr.Base.Pos)
		}

        case ast.RangeExpression: {
            c.expression(e.Start)
            c.expression(e.End)

            if (e.Step != nil) {
                c.expression(*e.Step)
            } else {
                // Push the constant 'nil', which defers the step number decision to runtime.
                c.writeBytePos(OP_PUSH_NIL, value.NewMetaLen1(expr.Base.Pos))
            }

			c.writeBytePos(OP_MAKE_RANGE, value.ChunkMetadata{
                Position: expr.Base.Pos,
                Length: 2,
            })
        }

		case ast.GetPropertyExpression: {
			c.expression(e.Left)

			// Store the name as a string in the constant table and retrieve it later.
			index := c.addConstant(value.ValueString{ Value: e.Property.Lexeme })

			c.writeBytePos(OP_GET_PROPERTY, value.ChunkMetadata{
				Position: e.Property.Pos,
				Length: len(e.Property.Lexeme),
			})
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.SetPropertyExpression: {
			c.expression(e.Left)
			c.expression(e.Value)

			// Store the name as a string in the constant table and retrieve it later.
			// The value to be assigned will be on top of the object we'll assign it to.
			index := c.addConstant(value.ValueString{ Value: e.Property.Lexeme })

			c.writeBytePos(OP_SET_PROPERTY, value.ChunkMetadata{
				Position: e.Property.Pos,
				Length: len(e.Property.Lexeme),
			})
			c.writeBytes(util.IntToBytes(index))
		}
	}
}
