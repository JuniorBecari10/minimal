package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) expectToken(kind token.TokenKind) (token.Token, diagnostic.Diagnostic) {
	if !p.check(kind) {
		return token.Token{}, p.makeExpectedTokenDiagnostic(kind)
	}

	return p.advance()
}

func (p *Parser) expectSemicolon() diagnostic.Diagnostic {
	if !p.match(token.TokenSemicolon) {
		// TODO: maybe make this error message more specific, like 'Expected semicolon after statement.'
		return p.makeExpectedTokenDiagnostic(token.TokenSemicolon)
	}

	return nil
}

func (p *Parser) exprCanSkipSemicolon(expr ast.ExprData) bool {
	switch expr.(type) {
		case ast.BlockExpression, ast.IfExpression:
			return true

		default:
			return false
	}
}

func (p *Parser) advance() (token.Token, diagnostic.Diagnostic) {
	p.previous = p.current
	p.current = p.next
	
	next, diag := p.lexer.Lex(p.next)

	if diag != nil {
		p.hadLexerError = true

		return p.current, diag
	}

	p.next = next
	return p.previous, nil
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
			break
		}

		switch p.current.Kind {
			case token.TokenLetKw, token.TokenVarKw, token.TokenFnKw, token.TokenRecordKw,
				 token.TokenIfKw, token.TokenWhileKw, token.TokenForKw:
				return
		}

		p.advance()
	}
}
