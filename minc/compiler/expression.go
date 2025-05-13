package compiler

import (
	"fmt"
	"minc/ast"
	"minlib/instructions"
	"minlib/token"
	"minlib/util"
	"minlib/value"
)

func (c *Compiler) expression(expr ast.Expression) {
	if c.hadError {
		return
	}

	switch e := expr.Data.(type) {
		case ast.IntExpression: {
			index := c.addConstant(value.ValueInt{ Value: e.Literal })

			c.writeBytePos(instructions.PUSH_CONST, value.Metadata{
				Position: expr.Base.Pos,
				Length: uint32(expr.Base.Length),
			})
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.FloatExpression: {
			index := c.addConstant(value.ValueFloat{ Value: e.Literal })

			c.writeBytePos(instructions.PUSH_CONST, value.Metadata{
				Position: expr.Base.Pos,
				Length: uint32(expr.Base.Length),
			})
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.StringExpression: {
			index := c.addConstant(value.ValueString{ Value: e.Literal })

			c.writeBytePos(instructions.PUSH_CONST, value.NewMetaLen1(expr.Base.Pos))
			c.writeBytes(util.IntToBytes(index))
		}
		
		case ast.CharExpression: {
			index := c.addConstant(value.ValueChar{ Value: e.Literal })

			c.writeBytePos(instructions.PUSH_CONST, value.NewMetaLen1(expr.Base.Pos))
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.BoolExpression: {
			if e.Literal {
				c.writeBytePos(instructions.PUSH_TRUE, value.Metadata{
					Position: expr.Base.Pos,
					Length: uint32(expr.Base.Length),
				})
			} else {
				c.writeBytePos(instructions.PUSH_FALSE, value.Metadata{
					Position: expr.Base.Pos,
					Length: uint32(expr.Base.Length),
				})
			}
		}

		case ast.NilExpression: {
			c.writeBytePos(instructions.PUSH_NIL, value.Metadata{
				Position: expr.Base.Pos,
				Length: uint32(expr.Base.Length),
			})
		}

		case ast.VoidExpression: {
			if e.Expr != nil {
				c.expression(*e.Expr)

				c.writeBytePos(instructions.POP, value.Metadata{
					Position: expr.Base.Pos,
					Length: uint32(expr.Base.Length),
				})
			}

			c.writeBytePos(instructions.PUSH_VOID, value.Metadata{
				Position: expr.Base.Pos,
				Length: uint32(expr.Base.Length),
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

            +-- instructions.JUMP_FALSE (and) / instructions.JUMP_TRUE (or)
            |   instructions.POP
            |
            |   [ right operand ]
            |   instructions.ASSERT_BOOL
            |
            +-> continues...
		*/
		case ast.LogicalExpression: {
			if e.ShortCircuit {
				var operation byte

				switch e.Operator.Kind {
					case token.TokenAndKw:
						operation = instructions.JUMP_FALSE
					case token.TokenOrKw:
						operation = instructions.JUMP_TRUE
					default:
						panic(fmt.Sprintf("Unknown logical operator: '%s'", e.Operator.Lexeme))
				}

				c.expression(e.Left)
				c.writeBytePos(operation, value.Metadata{
					Position: expr.Base.Pos,
					Length: uint32(expr.Base.Length),
				})

				jumpOffsetIndex := len(c.chunk.Code)
				c.writeBytes(util.IntToBytes(0)) // dummy
				c.writeBytePos(instructions.POP, value.Metadata{
					Position: expr.Base.Pos,
					Length: uint32(expr.Base.Length),
				})

				c.expression(e.Right)
				c.writeBytePos(instructions.ASSERT_BOOL, value.Metadata{
					Position: expr.Base.Pos,
					Length: uint32(expr.Base.Length),
				})

				c.backpatch(jumpOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpOffsetIndex - 4)) // index
			} else {
				c.expression(e.Left)
				c.expression(e.Right)

				c.chunk.Metadata = append(c.chunk.Metadata, value.Metadata{
					Position: e.Operator.Pos,
					Length: uint32(len(e.Operator.Lexeme)),
				})

				switch e.Operator.Kind {
					case token.TokenAndKw:
						c.writeByte(instructions.AND)
					case token.TokenOrKw:
						c.writeByte(instructions.OR)

					default:
						panic(fmt.Sprintf("Unknown logical operator: '%s'", e.Operator.Kind))
				}
			}
		}

		case ast.BinaryExpression: {
			c.expression(e.Left)
			c.expression(e.Right)

			c.chunk.Metadata = append(c.chunk.Metadata, value.Metadata{
				Position: e.Operator.Pos,
				Length: uint32(len(e.Operator.Lexeme)),
			})

			switch e.Operator.Kind {
				case token.TokenPlus:
					c.writeByte(instructions.ADD)
				case token.TokenMinus:
					c.writeByte(instructions.SUB)
				case token.TokenStar:
					c.writeByte(instructions.MUL)
				case token.TokenSlash:
					c.writeByte(instructions.DIV)
				
				case token.TokenPercent:
					c.writeByte(instructions.MOD)
				
				case token.TokenDoubleEqual:
					c.writeByte(instructions.EQUAL)
				case token.TokenBangEqual:
					c.writeByte(instructions.NOT_EQUAL)
				
				case token.TokenGreater:
					c.writeByte(instructions.GREATER)
				case token.TokenGreaterEqual:
					c.writeByte(instructions.GREATER_EQUAL)
				
				case token.TokenLess:
					c.writeByte(instructions.LESS)
				case token.TokenLessEqual:
					c.writeByte(instructions.LESS_EQUAL)
				
				case token.TokenAndKw:
					c.writeByte(instructions.AND)
				case token.TokenOrKw:
					c.writeByte(instructions.OR)

				default:
					panic(fmt.Sprintf("Unknown binary operator: '%s'", e.Operator.Kind))
			}
		}

		case ast.UnaryExpression: {
			c.expression(e.Operand)

			c.chunk.Metadata = append(c.chunk.Metadata, value.Metadata{
				Position: e.Operator.Pos,
				Length: uint32(len(e.Operator.Lexeme)),
			})

			switch e.Operator.Kind {
				case token.TokenNotKw:
					c.writeByte(instructions.NOT)
				case token.TokenMinus:
					c.writeByte(instructions.NEGATE)
				
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

					c.writeBytePos(instructions.CALL_PROPERTY, value.Metadata{
						Position: callee.Property.Pos,
						Length: uint32(len(callee.Property.Lexeme)),
					})
					c.writeBytes(util.IntToBytes(index))
					c.writeBytes(util.IntToBytes(len(e.Arguments)))
				}

				default: {
					c.expression(e.Callee)

					for _, arg := range e.Arguments {
						c.expression(arg)
					}

					c.writeBytePos(instructions.CALL, value.Metadata{
						Position: expr.Base.Pos,
						Length: uint32(expr.Base.Length),
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

			c.writeBytePos(byte(opcode), value.Metadata{
				Position: expr.Base.Pos,
				Length: uint32(expr.Base.Length),
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
                c.writeBytePos(instructions.PUSH_NIL, value.NewMetaLen1(expr.Base.Pos))
            }

            var opcode byte = instructions.MAKE_RANGE

            if e.Inclusive {
                opcode = instructions.MAKE_INCL_RANGE
            }

			c.writeBytePos(opcode, value.Metadata{
                Position: expr.Base.Pos,
                Length: 2,
            })
        }

		case ast.GetPropertyExpression: {
			c.expression(e.Left)

			// Store the name as a string in the constant table and retrieve it later.
			index := c.addConstant(value.ValueString{ Value: e.Property.Lexeme })

			c.writeBytePos(instructions.GET_PROPERTY, value.Metadata{
				Position: e.Property.Pos,
				Length: uint32(len(e.Property.Lexeme)),
			})
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.SetPropertyExpression: {
			c.expression(e.Left)
			c.expression(e.Value)

			// Store the name as a string in the constant table and retrieve it later.
			// The value to be assigned will be on top of the object we'll assign it to.
			index := c.addConstant(value.ValueString{ Value: e.Property.Lexeme })

			c.writeBytePos(instructions.SET_PROPERTY, value.Metadata{
				Position: e.Property.Pos,
				Length: uint32(len(e.Property.Lexeme)),
			})
			c.writeBytes(util.IntToBytes(index))
		}
	}
}
