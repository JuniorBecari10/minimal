package compiler

import (
	"bytes"
	"fmt"
	"vm-go/ast"
	"vm-go/token"
	"vm-go/util"
	"vm-go/value"
)

func (c *Compiler) expression(expr ast.Expression) []byte {
	res := bytes.Buffer{}

	if c.hadError {
		return res.Bytes()
	}

	switch e := expr.(type) {
		case ast.NumberExpression: {
			index := c.AddConstant(value.ValueNumber{ Value: e.Literal })

			res.WriteByte(OP_PUSH_CONST)
			res.WriteString(string(util.IntToBytes(index)))
		}

		case ast.StringExpression: {
			index := c.AddConstant(value.ValueString{ Value: e.Literal })

			res.WriteByte(OP_PUSH_CONST)
			res.WriteString(string(util.IntToBytes(index)))
		}

		case ast.BoolExpression: {
			index := c.AddConstant(value.ValueBool{ Value: e.Literal })

			res.WriteByte(OP_PUSH_CONST)
			res.WriteString(string(util.IntToBytes(index)))
		}

		case ast.NilExpression: {
			index := c.AddConstant(value.ValueNil{})

			res.WriteByte(OP_PUSH_CONST)
			res.WriteString(string(util.IntToBytes(index)))
		}

		case ast.IdentifierExpression: {
			index := c.resolveVariable(e.Ident)

			if index < 0 {
				return res.Bytes()
			}

			res.WriteByte(OP_GET_VAR)
			res.WriteString(string(util.IntToBytes(index)))
		}

		case ast.BinaryExpression: {
			res.WriteString(string(c.expression(e.Left)))
			res.WriteString(string(c.expression(e.Right)))

			switch e.Operator.Kind {
				case token.TokenPlus:
					res.WriteByte(OP_ADD)
				case token.TokenMinus:
					res.WriteByte(OP_SUB)
				case token.TokenStar:
					res.WriteByte(OP_MUL)
				case token.TokenSlash:
					res.WriteByte(OP_DIV)
				
				case token.TokenDoubleEqual:
					res.WriteByte(OP_EQUAL)
				case token.TokenBangEqual:
					res.WriteByte(OP_NOT_EQUAL)
				
				case token.TokenGreater:
					res.WriteByte(OP_GREATER)
				case token.TokenGreaterEqual:
					res.WriteByte(OP_GREATER_EQUAL)
				
				case token.TokenLess:
					res.WriteByte(OP_LESS)
				case token.TokenLessEqual:
					res.WriteByte(OP_LESS_EQUAL)
				
				case token.TokenAndKw:
					res.WriteByte(OP_AND)
				case token.TokenOrKw:
					res.WriteByte(OP_OR)
				case token.TokenXorKw:
					res.WriteByte(OP_XOR)

				default:
					panic(fmt.Sprintf("Unknown binary operator: '%s'", e.Operator.Kind))
			}
		}

		case ast.UnaryExpression: {
			res.WriteString(string(c.expression(e.Operand)))

			switch e.Operator.Kind {
				case token.TokenNotKw:
					res.WriteByte(OP_NOT)
				case token.TokenMinus:
					res.WriteByte(OP_NEGATE)
				
				default:
					panic(fmt.Sprintf("Unknown unary operator: '%s'", e.Operator.Kind))
			}
		}

		case ast.GroupExpression:
			res.WriteString(string(c.expression(e.Expr)))
		
		case ast.IdentifierAssignmentExpression: {
			index := c.resolveVariable(e.Name)

			if index < 0 {
				return res.Bytes()
			}

			res.WriteString(string(c.expression(e.Expr)))

			res.WriteByte(OP_SET_VAR)
			res.WriteString(string(util.IntToBytes(index)))
		}
	}

	return res.Bytes()
}
