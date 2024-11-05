package parser

import (
	"fmt"
	"strconv"
	"vm-go/ast"
	"vm-go/token"
	"vm-go/util"
)

type Parser struct {
	tokens []token.Token
	current int

	prefixMap map[token.TokenKind] func() ast.Expression
	infixMap map[token.TokenKind] func(ast.Expression, token.Position) ast.Expression
	precedenceMap map[token.TokenKind] int

	hadError bool
}

func NewParser(tokens []token.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		current: 0,

		hadError: false,
	}

	p.prefixMap = map[token.TokenKind] func() ast.Expression {
		token.TokenNumber: p.parseNumber,
		token.TokenLeftParen: p.parseGroup,
		token.TokenIdentifier: p.parseIdentifier,
	}

	p.infixMap = map[token.TokenKind] func(ast.Expression, token.Position) ast.Expression {
		token.TokenPlus: func (left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenPlus) },
		token.TokenMinus: func (left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenMinus) },
		token.TokenStar: func (left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenStar) },
		token.TokenSlash: func (left ast.Expression, pos token.Position) ast.Expression { return p.parseBinary(left, pos, token.TokenSlash) },
	}

	p.precedenceMap = map[token.TokenKind] int {
		token.TokenPlus: 1,
		token.TokenMinus: 1,

		token.TokenStar: 2,
		token.TokenSlash: 2,
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

// ---

func (p *Parser) statement() ast.Statement {
	t := p.peek()
	
	switch t.Kind {
		case token.TokenVarKw: return p.varStatement()
		case token.TokenPrintKw: return p.printStatement()
		case token.TokenLeftBrace: return p.blockStatement()
		
		default: return p.exprStatement()
	}
}

// ---

func (p *Parser) varStatement() ast.Statement {
	pos := p.advance().Pos // 'var' keyword
	name := p.expectToken(token.TokenIdentifier)
	p.expect(token.TokenEqual)

	expr := p.expression(0)
	p.requireSemicolon()

	return ast.VarStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Name: name,
		Init: expr,
	}
}

func (p *Parser) printStatement() ast.Statement {
	pos := p.advance().Pos // 'print' keyword
	expr := p.expression(0)

	p.requireSemicolon()

	return ast.PrintStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Expr: expr,
	}
}

func (p *Parser) blockStatement() ast.Statement {
	pos := p.advance().Pos // '{'
	stmts := []ast.Statement {}

	for !p.isAtEnd() && !p.check(token.TokenRightBrace) {
		stmts = append(stmts, p.statement())
	}

	p.advance()
	return ast.BlockStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Stmts: stmts,
	}
}

func (p *Parser) exprStatement() ast.Statement {
	pos := p.peek().Pos
	expr := p.expression(0)
	p.requireSemicolon()

	return ast.ExprStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},
		Expr: expr,
	}
}

// ---

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
		AstBase: ast.AstBase{ Pos: tok.Pos },
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
		AstBase: ast.AstBase{ Pos: pos },
		Expr: expr,
	}
}

func (p *Parser) parseBinary(left ast.Expression, pos token.Position, op token.TokenKind) ast.Expression {
	precedence := p.precedenceMap[op]

	operator := p.expectToken(op)
	right := p.expression(precedence)

	return ast.BinaryExpression{
		AstBase: ast.AstBase{ Pos: pos },
		Left: left,
		Right: right,
		Operator: operator,
	}
}

// ---

func (p *Parser) expect(kind token.TokenKind) bool {
	return !p.expectToken(kind).IsAbsent()
}

func (p *Parser) expectToken(kind token.TokenKind) token.Token {
	if !p.check(kind) {
		p.error(fmt.Sprintf("Expected '%s' token", kind))
		return token.AbsentToken()
	}

	return p.advance()
}

func (p *Parser) requireSemicolon() {
	p.expect(token.TokenSemicolon)
}

func (p *Parser) check(kind token.TokenKind) bool {
	return p.peek().Kind == kind
}

func (p *Parser) advance() token.Token {
	peek := p.peek()
	p.current += 1

	return peek
}

func (p *Parser) peek() token.Token {
	if p.isAtEnd() {
		return token.AbsentToken()
	}

	return p.tokens[p.current]
}

func (p * Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) error(message string) {
	if p.hadError {
		return
	}
	
	util.Error(p.peek().Pos, message)
	p.hadError = true
}
