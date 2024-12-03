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

	switch e := expr.(type) {
		case ast.NumberExpression: {
			index := c.addConstant(value.ValueNumber{ Value: e.Literal })

			c.writeBytePos(OP_PUSH_CONST, e.Pos)
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.StringExpression: {
			index := c.addConstant(value.ValueString{ Value: e.Literal })

			c.writeBytePos(OP_PUSH_CONST, e.Pos)
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.BoolExpression: {
			index := c.addConstant(value.ValueBool{ Value: e.Literal })

			c.writeBytePos(OP_PUSH_CONST, e.Pos)
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.NilExpression: {
			index := c.addConstant(value.ValueNil{})

			c.writeBytePos(OP_PUSH_CONST, e.Pos)
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.IdentifierExpression: {
			index := c.resolveVariable(e.Ident)

			if index < 0 {
				return
			}

			c.writeBytePos(OP_GET_VAR, e.Pos)
			c.writeBytes(util.IntToBytes(index))
		}

		case ast.BinaryExpression: {
			c.expression(e.Left)
			c.expression(e.Right)

			c.chunk.Positions = append(c.chunk.Positions, e.Operator.Pos)
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
				case token.TokenXorKw:
					c.writeByte(OP_XOR)

				default:
					panic(fmt.Sprintf("Unknown binary operator: '%s'", e.Operator.Kind))
			}
		}

		case ast.UnaryExpression: {
			c.expression(e.Operand)

			c.chunk.Positions = append(c.chunk.Positions, e.Operator.Pos)
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
			c.expression(e.Callee)

			for _, arg := range e.Arguments {
				c.expression(arg)
			}

			c.writeBytePos(OP_CALL, e.Pos)
			c.writeBytes(util.IntToBytes(len(e.Arguments)))
		}

		case ast.GroupExpression:
			c.expression(e.Expr)
		
		case ast.IdentifierAssignmentExpression: {
			index := c.resolveVariable(e.Name)

			if index < 0 {
				return
			}

			c.expression(e.Expr)

			c.writeBytePos(OP_SET_VAR, e.Pos)
			c.writeBytes(util.IntToBytes(index))
		}
	}
}
