package parser

import (
	"fmt"
	"vm-go/ast"
	"vm-go/token"
	"vm-go/util"
)

func (p *Parser) parseBlock() []ast.Statement {
	stmts := []ast.Statement {}

	for !p.isAtEnd() && !p.check(token.TokenRightBrace) && !p.hadError {
		stmts = append(stmts, p.statement())
	}

	p.expect(token.TokenRightBrace)
	return stmts
}

func (p *Parser) expect(kind token.TokenKind) bool {
	return !p.expectToken(kind).IsAbsent()
}

func (p *Parser) expectToken(kind token.TokenKind) token.Token {
	if !p.check(kind) {
		if p.isAtEnd() {
			p.error(fmt.Sprintf("Expected '%s', reached end", kind))
		} else {
			p.error(fmt.Sprintf("Expected '%s', got '%s'", kind, p.peek().Kind))
		}
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

func (p *Parser) match(kind token.TokenKind) bool {
	if p.check(kind) {
		p.advance()
		return true
	}

	return false
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

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) synchronize() {
	p.panicMode = false

	for !p.isAtEnd() {
		switch p.peek().Kind {
		case token.TokenVarKw, token.TokenLeftBrace, token.TokenIfKw:
			return
		}

		if p.peek().Kind == token.TokenSemicolon {
			return
		}

		p.advance()
	}
}

func (p *Parser) error(message string) {
	if p.panicMode {
		return
	}

	// TODO: if reached end, get the position of the last token
	util.Error(p.peek().Pos, message)

	p.hadError = true
	p.panicMode = true
}
