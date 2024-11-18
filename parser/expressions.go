package parser

import (
	"fmt"
	"strconv"
	"vm-go/ast"
	"vm-go/token"
)

func (p *Parser) expression(precedence int) ast.Expression {
	pos := p.peek().Pos
	prefixFn, ok := p.prefixMap[p.peek().Kind]

	if !ok {
		p.error(fmt.Sprintf("Unexpected token: '%s'", p.peek().Lexeme))
		return nil
	}

	left := prefixFn()

	for p.precedenceMap[p.peek().Kind] > precedence {
		infixFn, ok := p.infixMap[p.peek().Kind]

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

	return ast.NumberExpression{
		AstBase: ast.AstBase{Pos: tok.Pos},
		Literal: value,
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := p.expectToken(token.TokenIdentifier)

	return ast.IdentifierExpression{
		AstBase: ast.AstBase{
			Pos: ident.Pos,
		},

		Ident: ident,
	}
}

func (p *Parser) parseGroup() ast.Expression {
	pos := p.peek().Pos
	p.expect(token.TokenLeftParen)

	expr := p.expression(0)

	p.expect(token.TokenRightParen)

	return ast.GroupExpression{
		AstBase: ast.AstBase{Pos: pos},
		Expr:    expr,
	}
}

func (p *Parser) parseUnary(op token.TokenKind) ast.Expression {
	pos := p.peek().Pos
	precedence := p.precedenceMap[op]

	operator := p.expectToken(op)
	operand := p.expression(precedence)

	return ast.UnaryExpression{
		AstBase:  ast.AstBase{Pos: pos},
		Operand:  operand,
		Operator: operator,
	}
}

func (p *Parser) parseBinary(left ast.Expression, pos token.Position, op token.TokenKind) ast.Expression {
	precedence := p.precedenceMap[op]

	operator := p.expectToken(op)
	right := p.expression(precedence)

	return ast.BinaryExpression{
		AstBase:  ast.AstBase{Pos: pos},
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}
