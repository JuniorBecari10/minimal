package parser

import "minlib/token"

func (p *Parser) advance() token.Token {
	p.current = p.next
	next, diag := p.lexer.Lex()

	if diag != nil {
		diag.PrintDiagnostic()

		p.hadLexerError = true
		p.panicMode = true
	}

	p.next = next
	return p.current
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
	p.panicMode = false

	for !p.current.IsEnd() {
		kind := p.current.Kind

		switch kind {
			case
				// Declarations and control flow
				token.TokenVarKw,
				token.TokenFnKw,
				token.TokenRecordKw,
				token.TokenIfKw,
				token.TokenElseKw,
				token.TokenWhileKw,
				token.TokenForKw,
				token.TokenLoopKw,
				token.TokenBreakKw,
				token.TokenContinueKw,
				token.TokenReturnKw,

				// Block start/end
				token.TokenLeftBrace,
				token.TokenRightBrace,

				// Expression starters
				token.TokenIdentifier,
				token.TokenIntLiteral,
				token.TokenFloatLiteral,
				token.TokenCharLiteral,
				token.TokenStringLiteral,
				token.TokenMinus,
				token.TokenNotKw,
				token.TokenLeftParen,

				// Soft sync point
				token.TokenSemicolon:
					return
		}

		p.advance()
	}
}
