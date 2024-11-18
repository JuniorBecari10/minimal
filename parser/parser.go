package parser

import (
	"vm-go/ast"
	"vm-go/token"
)

const (
	PrecPrimary = iota	// literals, identifiers
	PrecCall			// ()
	PrecUnary			// not -
	PrecFactor			// * /
	PrecTerm			// + -
	PrecComparison			// < > <= >=
	PrecEqual			// == !=
	PrecAnd				// and
	PrecOr				// or
	PrecAssign          // =
)

type Parser struct {
	tokens []token.Token
	current int

	prefixMap map[token.TokenKind] func() ast.Expression
	infixMap map[token.TokenKind] func(ast.Expression, token.Position) ast.Expression
	precedenceMap map[token.TokenKind] int

	hadError bool
	panicMode bool
}

func NewParser(tokens []token.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		current: 0,

		hadError: false,
		panicMode: false,
	}

	p.prefixMap = map[token.TokenKind] func() ast.Expression {
		token.TokenNumber: p.parseNumber,
		token.TokenIdentifier: p.parseIdentifier,

		token.TokenLeftParen: p.parseGroup,
		token.TokenNotKw: func() ast.Expression { return p.parseUnary(token.TokenNotKw) },
	}

	p.infixMap = map[token.TokenKind]func(ast.Expression, token.Position) ast.Expression{
		token.TokenPlus:         func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenPlus) },
		token.TokenMinus:        func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenMinus) },
		
		token.TokenStar:         func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenStar) },
		token.TokenSlash:        func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenSlash) },
		
		token.TokenAndKw:        func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenAndKw) },
		token.TokenOrKw:         func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenOrKw) },
		token.TokenXorKw:         func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenXorKw) },
		
		token.TokenGreater:      func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenGreater) },
		token.TokenGreaterEqual: func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenGreaterEqual) },
		
		token.TokenLess:         func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenLess) },
		token.TokenLessEqual:    func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenLessEqual) },

		token.TokenDoubleEqual:   func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenDoubleEqual) },
		token.TokenBangEqual:     func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenBangEqual) },
	}

	p.precedenceMap = map[token.TokenKind] int {
		token.TokenPlus: PrecTerm,
		token.TokenMinus: PrecTerm,

		token.TokenStar: PrecFactor,
		token.TokenSlash: PrecFactor,

		token.TokenAndKw: PrecAnd,
		token.TokenOrKw: PrecOr,

		token.TokenGreater: PrecComparison,
		token.TokenGreaterEqual: PrecComparison,

		token.TokenLess: PrecComparison,
		token.TokenLessEqual: PrecComparison,

		token.TokenDoubleEqual: PrecEqual,
		token.TokenBangEqual: PrecEqual,
	}

	return p
}

func (p *Parser) Parse() ([]ast.Statement, bool) {
	stmts := []ast.Statement {}

	for !p.isAtEnd() {
		stmts = append(stmts, p.statement())
	}

	return stmts, p.hadError
}
