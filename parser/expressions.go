package parser

import (
	"fmt"
	"strconv"
	"vm-go/ast"
	"vm-go/token"
)

func (p *Parser) expression(precedence int) ast.Expression {
	pos := p.peek(0).Pos
	prefixFn, ok := p.prefixMap[p.peek(0).Kind]

	if !ok {
		p.error(fmt.Sprintf("Unexpected token: '%s'.", p.peek(0).Lexeme))
		return ast.Expression{}
	}

	left := prefixFn()

	for p.precedenceMap[p.peek(0).Kind] > precedence {
		infixFn, ok := p.infixMap[p.peek(0).Kind]

		if !ok {
			break
		}

		left = infixFn(left, pos)
	}

	return left
}

// ---

func (p *Parser) parseNumber() ast.Expression {
	tok := p.advance()
	value, _ := strconv.ParseFloat(tok.Lexeme, 64)

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme),
		},
		Data: ast.NumberExpression{
			Literal: value,
		},
	}
}

func (p *Parser) parseString() ast.Expression {
	tok := p.advance()

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme) + 2, // the quotes
		},
		Data: ast.StringExpression{
			Literal: tok.Lexeme,
		},
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := p.expectToken(token.TokenIdentifier)

	return ast.Expression{
		Base: ast.AstBase{
			Pos: ident.Pos,
			Length: len(ident.Lexeme),
		},
		Data: ast.IdentifierExpression{
			Token: ident,
		},
	}
}

func (p *Parser) parseSelf() ast.Expression {
	self := p.expectToken(token.TokenSelfKw)

	return ast.Expression{
		Base: ast.AstBase{
			Pos: self.Pos,
			Length: len(self.Lexeme),
		},
		Data: ast.SelfExpression{
			Token: self,
		},
	}
}

func (p *Parser) parseBool() ast.Expression {
	tok := p.advance()

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme),
		},
		Data: ast.BoolExpression{
			Literal: tok.Kind == token.TokenTrueKw,
		},
	}
}

func (p *Parser) parseNil() ast.Expression {
	tok := p.advance()

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme),
		},
		Data: ast.NilExpression{},
	}
}

func (p *Parser) parseVoid() ast.Expression {
	tok := p.advance()

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme),
		},
		Data: ast.VoidExpression{},
	}
}

func (p *Parser) lParen() ast.Expression {
	// lambda: ')' | ( ident ',' | ')' )
	if p.peek(1).Kind == token.TokenRightParen ||
		(p.peek(1).Kind == token.TokenIdentifier && (p.peek(2).Kind == token.TokenRightParen || p.peek(2).Kind == token.TokenComma)) {
		return p.parseLambda()
	} else {
		return p.parseGroup()
	}
}

func (p *Parser) parseLambda() ast.Expression {
	pos := p.peek(0).Pos
	parameters := p.parseParameters()

	p.expect(token.TokenArrow)

	var body ast.BlockStatement

	if p.peek(0).Kind == token.TokenLeftBrace {
		body = p.parseBlock()
	} else {
		peek := p.peek(0)
		expr := p.expression(PrecLowest)

		body = ast.BlockStatement{
			Stmts: []ast.Statement{
				{
					Base: ast.AstBase{
						Pos: peek.Pos,
						Length: len(peek.Lexeme),
					},
					Data: ast.ReturnStatement{
						Expression: &expr,
					},
				},
			},
		}
	}

	return ast.Expression{
		Base: ast.AstBase{
			Pos: pos,
			Length: 1, // TODO: change
		},
		Data: ast.FnExpression{
			Parameters: parameters,
			Body: body,
		},
	}
}

func (p *Parser) parseIfExpr() ast.Expression {
	if_ := p.advance()
	cond := p.expression(PrecLowest)

	p.expect(token.TokenColon)
	then := p.expression(PrecLowest)

	p.expect(token.TokenElseKw)
	var else_ ast.Expression

	// Check 'else if' chain and not require the colon if so.
	if p.check(token.TokenIfKw) {
		else_ = p.expression(PrecLowest)
	} else {
		p.expect(token.TokenColon)
		else_ = p.expression(PrecLowest)
	}

	return ast.Expression{
		Base: ast.AstBase{
			Pos: if_.Pos,
			Length: len(if_.Lexeme),
		},
		Data: ast.IfExpression{
			Condition: cond,
			Then: then,
			Else: else_,
		},
	}
}

func (p *Parser) parseGroup() ast.Expression {
	pos := p.peek(0).Pos
	p.expect(token.TokenLeftParen)

	expr := p.expression(PrecLowest)
	p.expect(token.TokenRightParen)

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    pos,
			Length: expr.Base.Length + 2, // The parentheses
		},
		Data: ast.GroupExpression{
			Expr: expr,
		},
	}
}

func (p *Parser) parseUnary(op token.TokenKind) ast.Expression {
	pos := p.peek(0).Pos

	operator := p.expectToken(op)
	operand := p.expression(PrecUnary)

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    pos,
			Length: len(operator.Lexeme) + operand.Base.Length,
		},
		Data: ast.UnaryExpression{
			Operand:  operand,
			Operator: operator,
		},
	}
}

// --- Infix ---

func (p *Parser) parseBinary(left ast.Expression, _ token.Position, op token.TokenKind) ast.Expression {
	precedence := p.precedenceMap[op]

	operator := p.expectToken(op)
	right := p.expression(precedence)

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    operator.Pos,
			Length: len(operator.Lexeme),
		},
		Data: ast.BinaryExpression{
			Left:     left,
			Right:    right,
			Operator: operator,
		},
	}
}

func (p *Parser) parseLogical(left ast.Expression, _ token.Position, op token.TokenKind) ast.Expression {
	precedence := p.precedenceMap[op]

	operator := p.expectToken(op)
	shortCircuit := !p.match(token.TokenStar)

	right := p.expression(precedence)

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    operator.Pos,
			Length: len(operator.Lexeme),
		},
		Data: ast.LogicalExpression{
			Left:         left,
			Right:        right,
			Operator:     operator,
			ShortCircuit: shortCircuit,
		},
	}
}

func (p *Parser) parseCall(left ast.Expression, pos token.Position) ast.Expression {
	p.expectToken(token.TokenLeftParen)
	arguments := []ast.Expression{}

	for !p.match(token.TokenRightParen) {
		arguments = append(arguments, p.expression(PrecLowest))

		if !p.check(token.TokenRightParen) {
			p.expect(token.TokenComma)
		}
	}

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    pos,
			Length: left.Base.Length, // the callee
		},
		Data: ast.CallExpression{
			Callee:    left,
			Arguments: arguments,
		},
	}
}

func (p *Parser) parseAssignment(left ast.Expression, pos token.Position) ast.Expression {
	equal := p.expectToken(token.TokenEqual)
	right := p.expression(PrecLowest) // accept one level higher because assignment is right-associative

	switch lValue := left.Data.(type) {
	case ast.IdentifierExpression:
		return ast.Expression{
			Base: ast.AstBase{
				Pos:    equal.Pos,
				Length: len(equal.Lexeme),
			},
			Data: ast.IdentifierAssignmentExpression{
				Name: lValue.Token,
				Expr: right,
			},
		}

	case ast.GetPropertyExpression:
		return ast.Expression{
			Base: ast.AstBase{
				Pos:  	lValue.Property.Pos,
				Length: len(lValue.Property.Lexeme),
			},
			Data: ast.SetPropertyExpression{
				Left:    lValue.Left,
				Property: lValue.Property,
				Value:   right,
			},
		}

	default:
		p.error(fmt.Sprintf("Invalid assignment target: '%v'.", left))
		return ast.Expression{}
	}
}

func (p *Parser) parseDot(left ast.Expression, pos token.Position) ast.Expression {
	p.advance() // Skip over the '.'
	property := p.expectToken(token.TokenIdentifier)

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    property.Pos,
			Length: len(property.Lexeme),
		},
		Data: ast.GetPropertyExpression{
			Left:    left,
			Property: property,
		},
	}
}
