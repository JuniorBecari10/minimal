package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
	"strconv"
)

type PrefixFn = func() (ast.Expression, diagnostic.Diagnostic)
type InfixFn = func(ast.Expression, token.Position) (ast.Expression, diagnostic.Diagnostic)
type Precedence int

const (
	PREC_LOWEST Precedence = iota
	PREC_ASSIGNMENT    // = += -= *= /= %=
	PREC_OR            // or
	PREC_AND           // and
	PREC_EQUAL         // == !=
	PREC_COMPARISON    // < > <= >=
	PREC_TERM          // + -
	PREC_FACTOR        // * /
	PREC_UNARY         // not -
	PREC_CALL          // ()
	PREC_GET_PROPERTY  // .
	PREC_RANGE         // ..
	// PREC_PRIMARY    // literals, identifiers (unused)
)

func (p *Parser) parseExpression() (ast.Expression, diagnostic.Diagnostic) {
	return p.expression(PREC_LOWEST)
}

func (p *Parser) expression(prec Precedence) (ast.Expression, diagnostic.Diagnostic) {
	leftPos := p.current.Pos
	prefixFn := p.getPrefixFn(p.current.Kind)

	if prefixFn == nil {
		return ast.Expression{}, p.makeExpectedExpressionDiagnostic()
	}

	left, diag := prefixFn(); if diag != nil {
		return ast.Expression{}, diag
	}

	for getPrecedenceFor(p.current.Kind) > prec {
		infixFn := p.getInfixFn(p.current.Kind)

		if infixFn == nil {
			break
		}

		left, diag = infixFn(left, leftPos); if diag != nil {
			return ast.Expression{}, diag
		}
	}

	return left, nil
}

func (p *Parser) getPrefixFn(kind token.TokenKind) PrefixFn {
	switch kind {
		case token.TokenIntLiteral: return p.parseInt
		case token.TokenFloatLiteral: return p.parseFloat
		case token.TokenStringLiteral: return p.parseStr
		case token.TokenCharLiteral: return p.parseChar
		case token.TokenIdentifier: return p.parseIdentifier
		case token.TokenSelfKw: return p.parseSelf

		case token.TokenTrueKw: return p.parseBool
		case token.TokenFalseKw: return p.parseBool
		case token.TokenNilKw: return p.parseNil
		case token.TokenVoidKw: return p.parseVoid

		case token.TokenLeftParen: return p.parseGroup
		case token.TokenLeftBrace: return p.parseBlockExpr

		case token.TokenIfKw: return p.parseIf
		case token.TokenFnKw: return p.parseFnExpr

		case token.TokenNotKw: return func() (ast.Expression, diagnostic.Diagnostic) { return p.parseUnary(token.TokenNotKw) }
		case token.TokenMinus: return func() (ast.Expression, diagnostic.Diagnostic) { return p.parseUnary(token.TokenMinus) }

		default: return nil
	}
}

func (p *Parser) getInfixFn(kind token.TokenKind) InfixFn {
	switch kind {
		case token.TokenPlus: return p.infixBinary(token.TokenPlus)
		case token.TokenMinus: return p.infixBinary(token.TokenMinus)

		case token.TokenStar: return p.infixBinary(token.TokenStar)
		case token.TokenSlash: return p.infixBinary(token.TokenSlash)
		case token.TokenPercent: return p.infixBinary(token.TokenPercent)

		case token.TokenPlusEqual: return p.infixOpAssign(token.TokenPlus)
		case token.TokenMinusEqual: return p.infixOpAssign(token.TokenMinus)
		case token.TokenStarEqual: return p.infixOpAssign(token.TokenStar)
		case token.TokenSlashEqual: return p.infixOpAssign(token.TokenSlash)
		case token.TokenPercentEqual: return p.infixOpAssign(token.TokenPercent)

		case token.TokenOrKw: return p.infixLogical(token.TokenOrKw)
		case token.TokenAndKw: return p.infixLogical(token.TokenAndKw)

		case token.TokenGreater: return p.infixBinary(token.TokenGreater)
		case token.TokenGreaterEqual: return p.infixBinary(token.TokenGreaterEqual)
		case token.TokenLess: return p.infixBinary(token.TokenLess)
		case token.TokenLessEqual: return p.infixBinary(token.TokenLessEqual)

		case token.TokenDoubleEqual: return p.infixBinary(token.TokenDoubleEqual)
		case token.TokenBangEqual: return p.infixBinary(token.TokenBangEqual)

		case token.TokenEqual: return p.parseAssignment
		case token.TokenLeftParen: return p.parseCall
		case token.TokenDot: return p.parseDot
		case token.TokenDoubleDot: return p.parseRange

		default: return nil
	}
}

func getPrecedenceFor(kind token.TokenKind) Precedence {
	switch kind {
		case token.TokenPlus, token.TokenMinus:
			return PREC_TERM

		case token.TokenStar, token.TokenSlash, token.TokenPercent:
			return PREC_FACTOR

		case token.TokenPlusEqual, token.TokenMinusEqual,
			token.TokenStarEqual, token.TokenSlashEqual, token.TokenPercentEqual,
			token.TokenEqual:
			return PREC_ASSIGNMENT

		case token.TokenAndKw:
			return PREC_AND

		case token.TokenOrKw:
			return PREC_OR

		case token.TokenGreater, token.TokenGreaterEqual,
			token.TokenLess, token.TokenLessEqual:
			return PREC_COMPARISON

		case token.TokenDoubleEqual, token.TokenBangEqual:
			return PREC_EQUAL

		case token.TokenLeftParen:
			return PREC_CALL

		case token.TokenDot:
			return PREC_GET_PROPERTY

		case token.TokenDoubleDot:
			return PREC_RANGE

		default:
			return PREC_LOWEST
	}
}

// ---

func (p *Parser) infixBinary(op token.TokenKind) InfixFn {
	return func(left ast.Expression, _ token.Position) (ast.Expression, diagnostic.Diagnostic) {
		return p.parseBinary(left, op)
	}
}

func (p *Parser) infixOpAssign(op token.TokenKind) InfixFn {
	return func(left ast.Expression, _ token.Position) (ast.Expression, diagnostic.Diagnostic) {
		return p.parseOperatorAssignment(left, op)
	}
}

func (p *Parser) infixLogical(op token.TokenKind) InfixFn {
	return func(left ast.Expression, _ token.Position) (ast.Expression, diagnostic.Diagnostic) {
		return p.parseLogical(left, op)
	}
}


// ---

func (p *Parser) parseInt() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()
	value, _ := strconv.Atoi(tok.Lexeme)

	return newExpr(tok, ast.IntExpression{
		Literal: int32(value),
	}), nil
}

func (p *Parser) parseFloat() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()
	value, _ := strconv.ParseFloat(tok.Lexeme, 64)

	return newExpr(tok, ast.FloatExpression{
		Literal: value,
	}), nil
}

func (p *Parser) parseStr() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()

	return newExprLength(tok.Pos, len(tok.Lexeme) + 2, ast.StringExpression{
		Literal: tok.Lexeme,
	}), nil
}

func (p *Parser) parseChar() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()

	return newExprLength(tok.Pos, len(tok.Lexeme) + 2, ast.CharExpression{
		Literal: uint8(tok.Lexeme[0]), // guaranteed to be one character long
	}), nil
}

func (p *Parser) parseIdentifier() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()

	return newExpr(tok, ast.IdentifierExpression{
		Token: tok,
	}), nil
}

func (p *Parser) parseSelf() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()
	return newExpr(tok, ast.SelfExpression{}), nil
}

func (p *Parser) parseBool() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()

	return newExpr(tok, ast.BoolExpression{
		Literal: tok.Kind == token.TokenTrueKw,
	}), nil
}

func (p *Parser) parseNil() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()
	return newExpr(tok, ast.NilExpression{}), nil
}

func (p *Parser) parseVoid() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()

	var expr *ast.Expression = nil

	if p.match(token.TokenLeftParen) {
		e, diag := p.parseExpression(); if diag != nil {
			return ast.Expression{}, diag
		}

		p.expect(token.TokenRightParen)
		*expr = e
	}

	return newExpr(tok, ast.VoidExpression{
		Expr: expr,
	}), nil
}

func (p *Parser) parseGroup() (ast.Expression, diagnostic.Diagnostic) {
	pos := p.current.Pos

	_, diag := p.expectToken(token.TokenLeftParen); if diag != nil {
		return ast.Expression{}, diag
	}

	expr, diag := p.parseExpression(); if diag != nil {
		return ast.Expression{}, diag
	}

	_, diag = p.expectToken(token.TokenRightParen); if diag != nil {
		return ast.Expression{}, diag
	}

	return newExprLength(pos, expr.Base.Length + 2, ast.GroupExpression{
		Expr: expr,
	}), nil
}

func (p *Parser) parseBlockExpr() (ast.Expression, diagnostic.Diagnostic) {

}

func (p *Parser) parseIf() (ast.Expression, diagnostic.Diagnostic) {

}

func (p *Parser) parseFnExpr() (ast.Expression, diagnostic.Diagnostic) {

}

// ---

func (p *Parser) parseAssignment(left ast.Expression, pos token.Position) (ast.Expression, diagnostic.Diagnostic) {

}

func (p *Parser) parseCall(left ast.Expression, pos token.Position) (ast.Expression, diagnostic.Diagnostic) {

}

func (p *Parser) parseDot(left ast.Expression, pos token.Position) (ast.Expression, diagnostic.Diagnostic) {

}

func (p *Parser) parseRange(left ast.Expression, pos token.Position) (ast.Expression, diagnostic.Diagnostic) {

}

// ---

func (p *Parser) parseUnary(kind token.TokenKind) (ast.Expression, diagnostic.Diagnostic) {

}

func (p *Parser) parseBinary(left ast.Expression, kind token.TokenKind) (ast.Expression, diagnostic.Diagnostic) {

}

func (p *Parser) parseOperatorAssignment(left ast.Expression, kind token.TokenKind) (ast.Expression, diagnostic.Diagnostic) {

}

func (p *Parser) parseLogical(left ast.Expression, kind token.TokenKind) (ast.Expression, diagnostic.Diagnostic) {

}

// ---

func newExpr(tok token.Token, data ast.ExprData) ast.Expression {
	return ast.Expression{
		Base: ast.AstBase{
			Pos:    tok.Pos,
			Length: len(tok.Lexeme),
		},
		Data: data,
	}
}

func newExprLength(pos token.Position, length int, data ast.ExprData) ast.Expression {
	return ast.Expression{
		Base: ast.AstBase{
			Pos:    pos,
			Length: length,
		},
		Data: data,
	}
}
