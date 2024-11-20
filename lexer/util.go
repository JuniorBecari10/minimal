package lexer

import (
	"vm-go/token"
	"vm-go/util"
)

func (l *Lexer) match(c byte) bool {
	if l.isAtEnd() {
		return false
	}

	if l.source[l.current] == c {
		l.advance()
		return true
	}

	return false
}

func (l *Lexer) advance() byte {
	peek := l.peek()
	l.current += 1

	if peek == '\n' {
		l.increaseLine()
	} else {
		l.currentPos.Col += 1
	}

	return peek
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}

	return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.isAtEnd() {
		return 0
	}

	return l.source[l.current+1]
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) increaseLine() {
	l.currentPos.Line += 1
	l.currentPos.Col = 0
}

func (l *Lexer) error(message string) {
	util.Error(l.startPos, message)
	l.hadError = true
}

// ---

func (l *Lexer) addToken(kind token.TokenKind) {
	l.tokens = append(l.tokens, token.Token{
		Kind:   kind,
		Lexeme: l.source[l.start:l.current],
		Pos:    l.startPos,
	})
}

func (l *Lexer) addTokenLexeme(kind token.TokenKind, lexeme string) {
	l.tokens = append(l.tokens, token.Token{
		Kind:   kind,
		Lexeme: lexeme,
		Pos:    l.startPos,
	})
}
