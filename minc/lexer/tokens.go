package lexer

import (
	"minlib/token"
	"unicode"
)

func (l *Lexer) number() {
	for unicode.IsDigit(rune(l.peek(0))) {
		l.advance()
	}

	if l.peek(0) == '.' && unicode.IsDigit(rune(l.peek(1))) {
		l.advance()

		for unicode.IsDigit(rune(l.peek(0))) {
			l.advance()
		}
	}

	l.addToken(token.TokenNumber)
}

func (l *Lexer) string() {
	for l.peek(0) != '"' && !l.isAtEnd(0) {
		l.advance()
	}

	if l.isAtEnd(0) {
		l.error("Unterminated string")
		return
	}

	l.advance() // the closing '"'
	l.addTokenLexeme(token.TokenString, l.source[l.start + 1 : l.current - 1])
}

func (l *Lexer) identifier() {
	for unicode.IsLetter(rune(l.peek(0))) || unicode.IsDigit(rune(l.peek(0))) || l.peek(0) == '_' {
		l.advance()
	}

	l.addToken(l.checkKeyword())
}

func (l *Lexer) checkKeyword() token.TokenKind {
	switch l.source[l.start:l.current] {
		case "if": return token.TokenIfKw
		case "else": return token.TokenElseKw
		case "while": return token.TokenWhileKw
		case "for": return token.TokenForKw
		case "loop": return token.TokenLoopKw
		case "var": return token.TokenVarKw
		case "fn": return token.TokenFnKw
		case "break": return token.TokenBreakKw
		case "continue": return token.TokenContinueKw
		case "in": return token.TokenInKw
		case "self": return token.TokenSelfKw
		case "record": return token.TokenRecordKw
		case "return": return token.TokenReturnKw

		case "and": return token.TokenAndKw
		case "or": return token.TokenOrKw
		case "not": return token.TokenNotKw

		case "true": return token.TokenTrueKw
		case "false": return token.TokenFalseKw
		case "nil": return token.TokenNilKw
		case "void": return token.TokenVoidKw

		default: return token.TokenIdentifier
	}
}
