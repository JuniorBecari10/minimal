package lexer

import (
	"unicode"
	"vm-go/token"
)

func (l *Lexer) number() {
	for unicode.IsDigit(rune(l.peek())) {
		l.advance()
	}

	if l.match('.') && unicode.IsDigit(rune(l.peekNext())) {
		for unicode.IsDigit(rune(l.peek())) {
			l.advance()
		}
	}

	l.addToken(token.TokenNumber)
}

func (l *Lexer) string() {
	for l.peek() != '"' && !l.isAtEnd() {
		l.advance()
	}

	if l.isAtEnd() {
		l.error("Unterminated string")
		return
	}

	l.advance() // the closing '"'
	l.addTokenLexeme(token.TokenString, l.source[l.start + 1 : l.current - 1])
}

func (l *Lexer) identifier() {
	for unicode.IsLetter(rune(l.peek())) || l.peek() == '_' {
		l.advance()
	}

	l.addToken(l.checkKeyword())
}

func (l *Lexer) checkKeyword() token.TokenKind {
	switch l.source[l.start:l.current] {
		case "if": return token.TokenIfKw
		case "else": return token.TokenElseKw
		case "while": return token.TokenWhileKw
		case "var": return token.TokenVarKw
		case "print": return token.TokenPrintKw

		case "and": return token.TokenAndKw
		case "or": return token.TokenOrKw
		case "xor": return token.TokenXorKw
		case "not": return token.TokenNotKw

		case "true": return token.TokenTrueKw
		case "false": return token.TokenFalseKw
		case "nil": return token.TokenNilKw

		default: return token.TokenIdentifier
	}
}
