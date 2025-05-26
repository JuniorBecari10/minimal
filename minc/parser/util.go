package parser

import (
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) expect(kind token.TokenKind) bool {
	_, diag := p.expectToken(kind)
	return diag == nil
}

func (p *Parser) expectToken(kind token.TokenKind) (token.Token, diagnostic.Diagnostic) {
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
		kind := p.current.Kind

		switch kind {
			case
				token.TokenRightBrace,
				token.TokenSemicolon:
					return
		}

		p.advance()
	}
	p.advance()
}
