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

		case ast.IdentifierExpression: {
			for i := len(c.variables) - 1; i >= 0; i-- {
				if c.variables[i].name.Lexeme == e.Ident.Lexeme {
					if !c.variables[i].initialized {
						// TODO check if the scopeDepth is not 0 and allow its use,
						// since functions are allowed in top level and the use of variables inside them is allowed,
						// because it is guaranteed to have them defined
						c.error(e.Pos, fmt.Sprintf("'%s' is not defined yet", e.Ident.Lexeme))
						return res.Bytes()
					}

					res.WriteByte(OP_GET_VAR)
					res.WriteString(string(util.IntToBytes(i)))

					return res.Bytes()
				}
			}

			c.error(e.Pos, fmt.Sprintf("'%s' doesn't exist", e.Ident.Lexeme))
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

				default:
					panic(fmt.Sprintf("Unknown binary operator: '%s'", e.Operator.Kind))
			}
		}

		case ast.GroupExpression:
			res.WriteString(string(c.expression(e.Expr)))
	}

	return res.Bytes()
}
