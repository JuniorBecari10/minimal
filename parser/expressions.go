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
		p.error(fmt.Sprintf("Unexpected token: '%s'", p.peek(0).Lexeme))
		return nil
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

	return ast.NumberExpression{
		AstBase: ast.AstBase{
			Pos: tok.Pos,
		},

		Literal: value,
	}
}

func (p *Parser) parseString() ast.Expression {
	tok := p.advance()

	return ast.StringExpression{
		AstBase: ast.AstBase{
			Pos: tok.Pos,
		},

		Literal: tok.Lexeme,
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

func (p *Parser) parseBool() ast.Expression {
	tok := p.advance()

	return ast.BoolExpression{
		AstBase: ast.AstBase{
			Pos: tok.Pos,
		},

		Literal: tok.Kind == token.TokenTrueKw,
	}
}

func (p *Parser) parseNil() ast.Expression {
	tok := p.advance()

	return ast.NilExpression{
		AstBase: ast.AstBase{
			Pos: tok.Pos,
		},
	}
}

func (p *Parser) parseVoid() ast.Expression {
	tok := p.advance()

	return ast.VoidExpression{
		AstBase: ast.AstBase{
			Pos: tok.Pos,
		},
	}
}

func (p *Parser) parseGroup() ast.Expression {
	pos := p.peek(0).Pos
	p.expect(token.TokenLeftParen)

	expr := p.expression(PrecLowest)
	p.expect(token.TokenRightParen)

	return ast.GroupExpression{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Expr:    expr,
	}
}

func (p *Parser) parseUnary(op token.TokenKind) ast.Expression {
	pos := p.peek(0).Pos

	operator := p.expectToken(op)
	operand := p.expression(PrecUnary)

	return ast.UnaryExpression{
		AstBase:  ast.AstBase{
			Pos: pos,
		},

		Operand:  operand,
		Operator: operator,
	}
}

func (p *Parser) parseBinary(left ast.Expression, pos token.Position, op token.TokenKind) ast.Expression {
	precedence := p.precedenceMap[op]

	operator := p.expectToken(op)
	right := p.expression(precedence)

	return ast.BinaryExpression{
		AstBase:  ast.AstBase{
			Pos: pos,
		},

		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (p *Parser) parseLogical(left ast.Expression, pos token.Position, op token.TokenKind) ast.Expression {
	precedence := p.precedenceMap[op]

	operator := p.expectToken(op)
	shortCircuit := !p.match(token.TokenStar)

	right := p.expression(precedence)

	return ast.LogicalExpression{
		AstBase:  ast.AstBase{
			Pos: pos,
		},

		Left:     left,
		Right:    right,
		Operator: operator,
		ShortCircuit: shortCircuit,
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

	return ast.CallExpression{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Callee: left,
		Arguments: arguments,
	}
}

func (p *Parser) parseAssignment(left ast.Expression, pos token.Position) ast.Expression {
	p.expectToken(token.TokenEqual)
	right := p.expression(PrecLowest) // accept one level higher because assignment is right-associative

	name, ok := left.(ast.IdentifierExpression)

	if !ok {
		p.error(fmt.Sprintf("Invalid assignment target: '%v'", left))
	}

	return ast.IdentifierAssignmentExpression{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Name: name.Ident,
		Expr: right,
	}
}
