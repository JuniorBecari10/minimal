package parser

import (
	"fmt"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) expect(kind token.TokenKind) bool {
	_, diag := p.expectToken(kind)
	return diag == nil
}

func (p *Parser) expectToken(kind token.TokenKind) (token.Token, diagnostic.Diagnostic) {
	fmt.Println(p.current, kind)
	if !p.check(kind) {
		return token.Token{}, p.makeExpectedTokenDiagnostic(kind)
	}

	return p.advance()
}

func (p *Parser) expectSemicolon() diagnostic.Diagnostic {
	if !p.match(token.TokenSemicolon) {
		// TODO: maybe make this error message more specific
		return p.makeExpectedTokenDiagnostic(token.TokenSemicolon)
	}

	return nil
}

func (p *Parser) advance() (token.Token, diagnostic.Diagnostic) {
	p.current = p.next
	next, diag := p.lexer.Lex()

	if diag != nil {
		p.hadLexerError = true

		return p.current, diag
	}

	p.previous = p.current
	p.next = next
	return p.current, nil
}

func (p *Parser) check(kind token.TokenKind) bool {
	return p.current.Kind == kind
}

// This advances.
func (p *Parser) match(kind token.TokenKind) bool {
	if p.check(kind) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) synchronize() {
	for !p.current.IsEnd() {
		if p.previous.Kind == token.TokenSemicolon {
			p.advance()
			break
		}

		switch p.current.Kind {
			case token.TokenLetKw, token.TokenVarKw, token.TokenFnKw, token.TokenRecordKw,
				 token.TokenIfKw, token.TokenWhileKw, token.TokenForKw:
				break
		}

		p.advance()
	}
}

