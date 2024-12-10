package parser

import (
	"vm-go/ast"
	"vm-go/token"
	"vm-go/util"
)

const (
	PrecLowest = iota
	PrecAssignment          // =
	PrecOr                  // or
	PrecXor                 // xor
	PrecAnd                 // and
	PrecEqual               // == !=
	PrecComparison          // < > <= >=
	PrecTerm                // + -
	PrecFactor              // * /
	PrecUnary               // not -
	PrecCall                // ()
	PrecPrimary             // literals, identifiers
)


type Parser struct {
	tokens []token.Token
	current int

	prefixMap map[token.TokenKind] func() ast.Expression
	infixMap map[token.TokenKind] func(ast.Expression, token.Position) ast.Expression
	precedenceMap map[token.TokenKind] int // add another map if necessary

	hadError bool
	panicMode bool

	fileData *util.FileData
}

func NewParser(tokens []token.Token, fileData *util.FileData) *Parser {
	p := &Parser{
		tokens: tokens,
		current: 0,

		hadError: false,
		panicMode: false,

		fileData: fileData,
	}

	p.prefixMap = map[token.TokenKind] func() ast.Expression {
		token.TokenNumber: p.parseNumber,
		token.TokenString: p.parseString,
		token.TokenIdentifier: p.parseIdentifier,

		token.TokenTrueKw: p.parseBool,
		token.TokenFalseKw: p.parseBool,
		token.TokenNilKw: p.parseNil,
		token.TokenVoidKw: p.parseVoid,

		token.TokenLeftParen: p.parseGroup,

		token.TokenNotKw: func() ast.Expression { return p.parseUnary(token.TokenNotKw) },
		token.TokenMinus: func() ast.Expression { return p.parseUnary(token.TokenMinus) },
	}

	p.infixMap = map[token.TokenKind]func(ast.Expression, token.Position) ast.Expression{
		token.TokenPlus:         func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenPlus) },
		token.TokenMinus:        func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenMinus) },
		
		token.TokenStar:         func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenStar) },
		token.TokenSlash:        func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenSlash) },
		token.TokenPercent:      func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenPercent) },
		
		token.TokenAndKw:        func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenAndKw) },
		token.TokenOrKw:         func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenOrKw) },
		token.TokenXorKw:        func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenXorKw) },
		
		token.TokenGreater:      func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenGreater) },
		token.TokenGreaterEqual: func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenGreaterEqual) },
		
		token.TokenLess:         func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenLess) },
		token.TokenLessEqual:    func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenLessEqual) },

		token.TokenDoubleEqual:  func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenDoubleEqual) },
		token.TokenBangEqual:    func(left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenBangEqual) },

		token.TokenEqual: p.parseAssignment,
		token.TokenLeftParen: p.parseCall,
	}

	p.precedenceMap = map[token.TokenKind] int {
		token.TokenPlus: PrecTerm,
		token.TokenMinus: PrecTerm,

		token.TokenStar: PrecFactor,
		token.TokenSlash: PrecFactor,
		token.TokenPercent: PrecFactor,

		token.TokenAndKw: PrecAnd,
		token.TokenOrKw: PrecOr,
		token.TokenXorKw: PrecXor,

		token.TokenGreater: PrecComparison,
		token.TokenGreaterEqual: PrecComparison,

		token.TokenLess: PrecComparison,
		token.TokenLessEqual: PrecComparison,

		token.TokenDoubleEqual: PrecEqual,
		token.TokenBangEqual: PrecEqual,

		token.TokenEqual: PrecAssignment,
		token.TokenLeftParen: PrecCall,
	}

	return p
}

func (p *Parser) Parse() ([]ast.Statement, bool) {
	stmts := []ast.Statement{}

	for !p.isAtEnd(0) {
		stmts = append(stmts, p.declaration(false))
	}

	return stmts, p.hadError
}
