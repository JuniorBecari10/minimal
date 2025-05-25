package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
)

type PrefixFn = func() ast.Expression
type InfixFn = func(ast.Expression, token.Position) ast.Expression
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
	prefixFn := getPrefixFn(p.current.Kind)

	if prefixFn == nil {
		return ast.Expression{}, p.makeExpectedExpressionDiagnostic()
	}

	left := prefixFn()

	for getPrecedenceFor(p.current.Kind) > prec {
		infixFn := getInfixFn(p.current.Kind)

		if infixFn == nil {
			break
		}

		left = infixFn(left, leftPos)
	}

	return left, nil
}

func getPrefixFn(kind token.TokenKind) PrefixFn {

}

func getInfixFn(kind token.TokenKind) InfixFn {

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


