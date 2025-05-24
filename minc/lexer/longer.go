package lexer

import (
	"minc/diagnostic"
	"minlib/token"
	"unicode"
)

func (l *Lexer) number() token.Token {
	tokenType := token.TokenKind(token.TokenIntLiteral)

	for unicode.IsDigit(rune(l.peek(0))) {
		l.advance()
	}

	if l.peek(0) == '.' && unicode.IsDigit(rune(l.peek(1))) {
		tokenType = token.TokenFloatLiteral
		l.advance()

		for unicode.IsDigit(rune(l.peek(0))) {
			l.advance()
		}
	}

	return l.makeToken(tokenType)
}

func (l *Lexer) string() (token.Token, diagnostic.Diagnostic) {
	for l.peek(0) != '"' && !l.isAtEnd(0) {
		l.advance()
	}

	if l.isAtEnd(0) {
		return token.Token{}, l.makeUnterminatedStringLiteralDiagnostic()
	}

	l.advance() // the closing '"'
	return l.makeTokenLexeme(token.TokenStringLiteral, l.source[l.start + 1 : l.current - 1]), nil
}

func (l *Lexer) char() (token.Token, diagnostic.Diagnostic) {
	for l.peek(0) != '\'' && !l.isAtEnd(0) {
		l.advance()
	}

	if l.isAtEnd(0) {
		return token.Token{}, l.makeUnterminatedCharLiteralDiagnostic()
	}

	if l.current - l.start - 2 > 1 {
		return token.Token{}, l.makeCharLiteralTooLongDiagnostic()
	}

	l.advance() // the closing '\''
	return l.makeTokenLexeme(token.TokenCharLiteral, l.source[l.start + 1 : l.current - 1]), nil
}

func (l *Lexer) identifier() token.Token {
	for unicode.IsLetter(rune(l.peek(0))) || unicode.IsDigit(rune(l.peek(0))) || l.peek(0) == '_' {
		l.advance()
	}

	return l.makeToken(l.checkKeyword())
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
