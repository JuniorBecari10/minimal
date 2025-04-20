package parser

import (
	"minc/ast"
	"minlib/token"
	"minlib/util"
)

const (
	PrecLowest = iota
	PrecAssignment          // = += -= *= /= %=
	PrecOr                  // or
	PrecAnd                 // and
	PrecEqual               // == !=
	PrecComparison          // < > <= >=
	PrecTerm                // + -
	PrecFactor              // * /
	PrecUnary               // not -
	PrecCall                // ()
	PrecGetProperty			// .
	PrecRange				// ..
	// PrecPrimary          // literals, identifiers (unused)
)


type Parser struct {
	tokens []token.Token
	current int

	prefixMap map[token.TokenKind]func() ast.Expression
	infixMap map[token.TokenKind]func(ast.Expression, token.Position) ast.Expression
	precedenceMap map[token.TokenKind]int // add another map if necessary

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
		token.TokenSelfKw: p.parseSelf,

		token.TokenTrueKw: p.parseBool,
		token.TokenFalseKw: p.parseBool,
		token.TokenNilKw: p.parseNil,
		token.TokenVoidKw: p.parseVoid,

		token.TokenLeftParen: p.lParen,
		token.TokenIfKw: p.parseIfExpr,

		token.TokenNotKw: func() ast.Expression { return p.parseUnary(token.TokenNotKw) },
		token.TokenMinus: func() ast.Expression { return p.parseUnary(token.TokenMinus) },
	}

	p.infixMap = map[token.TokenKind]func(ast.Expression, token.Position) ast.Expression{
		token.TokenPlus:         func(left ast.Expression, _ token.Position) ast.Expression { return p.parseBinary(left, token.TokenPlus) },
		token.TokenMinus:        func(left ast.Expression, _ token.Position) ast.Expression { return p.parseBinary(left, token.TokenMinus) },
		
		token.TokenStar:         func(left ast.Expression, _ token.Position) ast.Expression { return p.parseBinary(left, token.TokenStar) },
		token.TokenSlash:        func(left ast.Expression, _ token.Position) ast.Expression { return p.parseBinary(left, token.TokenSlash) },
		token.TokenPercent:      func(left ast.Expression, _ token.Position) ast.Expression { return p.parseBinary(left, token.TokenPercent) },

		token.TokenPlusEqual:    func(left ast.Expression, _ token.Position) ast.Expression { return p.parseOperatorAssignment(left, token.TokenPlus) },
		token.TokenMinusEqual:   func(left ast.Expression, _ token.Position) ast.Expression { return p.parseOperatorAssignment(left, token.TokenMinus) },
		
		token.TokenStarEqual:    func(left ast.Expression, _ token.Position) ast.Expression { return p.parseOperatorAssignment(left, token.TokenStar) },
		token.TokenSlashEqual:   func(left ast.Expression, _ token.Position) ast.Expression { return p.parseOperatorAssignment(left, token.TokenSlash) },
		token.TokenPercentEqual: func(left ast.Expression, _ token.Position) ast.Expression { return p.parseOperatorAssignment(left, token.TokenPercent) },
		
		token.TokenOrKw:         func(left ast.Expression, _ token.Position) ast.Expression { return p.parseLogical(left, token.TokenOrKw) },
		token.TokenAndKw:        func(left ast.Expression, _ token.Position) ast.Expression { return p.parseLogical(left, token.TokenAndKw) },
		
		token.TokenGreater:      func(left ast.Expression, _ token.Position) ast.Expression { return p.parseBinary(left, token.TokenGreater) },
		token.TokenGreaterEqual: func(left ast.Expression, _ token.Position) ast.Expression { return p.parseBinary(left, token.TokenGreaterEqual) },
		
		token.TokenLess:         func(left ast.Expression, _ token.Position) ast.Expression { return p.parseBinary(left, token.TokenLess) },
		token.TokenLessEqual:    func(left ast.Expression, _ token.Position) ast.Expression { return p.parseBinary(left, token.TokenLessEqual) },

		token.TokenDoubleEqual:  func(left ast.Expression, _ token.Position) ast.Expression { return p.parseBinary(left, token.TokenDoubleEqual) },
		token.TokenBangEqual:    func(left ast.Expression, _ token.Position) ast.Expression { return p.parseBinary(left, token.TokenBangEqual) },

		token.TokenEqual: p.parseAssignment,
		token.TokenLeftParen: p.parseCall,
		
		token.TokenDot: p.parseDot,
		token.TokenDoubleDot: p.parseRange,
	}

	p.precedenceMap = map[token.TokenKind] int {
		token.TokenPlus: PrecTerm,
		token.TokenMinus: PrecTerm,

		token.TokenStar: PrecFactor,
		token.TokenSlash: PrecFactor,
		token.TokenPercent: PrecFactor,

		token.TokenPlusEqual: PrecAssignment,
		token.TokenMinusEqual: PrecAssignment,

		token.TokenStarEqual: PrecAssignment,
		token.TokenSlashEqual: PrecAssignment,
		token.TokenPercentEqual: PrecAssignment,

		token.TokenAndKw: PrecAnd,
		token.TokenOrKw: PrecOr,

		token.TokenGreater: PrecComparison,
		token.TokenGreaterEqual: PrecComparison,

		token.TokenLess: PrecComparison,
		token.TokenLessEqual: PrecComparison,

		token.TokenDoubleEqual: PrecEqual,
		token.TokenBangEqual: PrecEqual,

		token.TokenEqual: PrecAssignment,
		token.TokenLeftParen: PrecCall,

		token.TokenDot: PrecGetProperty,
		token.TokenDoubleDot: PrecRange,
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
