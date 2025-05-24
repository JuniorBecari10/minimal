package lexer

import (
	"fmt"
	"minlib/token"
	"minc/diagnostic"
)

func (l *Lexer) match(c byte) bool {
	if l.isAtEnd(0) {
		return false
	}

	if l.source[l.current] == c {
		l.advance()
		return true
	}

	return false
}

func (l *Lexer) advance() byte {
	peek := l.peek(0)
	l.current += 1

	if peek == '\n' {
		l.increaseLine()
	} else {
		l.currentPos.Col += 1
	}

	return peek
}

func (l *Lexer) peek(offset int) byte {
	if l.isAtEnd(offset) {
		return 0
	}

	return l.source[l.current + offset]
}

func (l *Lexer) isAtEnd(offset int) bool {
	return l.current + offset >= len(l.source)
}

func (l *Lexer) increaseLine() {
	l.currentPos.Line += 1
	l.currentPos.Col = 0
}

// ---
func (l *Lexer) makeUnknownTokenDiagnostic(c byte) diagnostic.SimpleDiagnostic {
	return l.makeUnterminatedDiagnostic(fmt.Sprintf("Unknown token: '%c' (char '%d').", c, c))
}

func (l *Lexer) makeUnterminatedStringLiteralDiagnostic() diagnostic.SimpleDiagnostic {
	return l.makeUnterminatedDiagnostic("Unterminated string literal.")
}

func (l *Lexer) makeUnterminatedCharLiteralDiagnostic() diagnostic.SimpleDiagnostic {
	return l.makeUnterminatedDiagnostic("Unterminated char literal.")
}

func (l *Lexer) makeCharLiteralTooLongDiagnostic() diagnostic.SimpleDiagnostic {
	return l.makeUnterminatedDiagnostic("Char literal too long.")
}

func (l *Lexer) makeUnterminatedDiagnostic(message string) diagnostic.SimpleDiagnostic {
	return diagnostic.SimpleDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: message,
			Span: diagnostic.Span{
				Pos: l.startPos,
				Length: int(l.currentPos.Col - l.startPos.Col),
			},
			FileData: l.fileData,
		},
	}
}

// ---

func (l *Lexer) makeToken(kind token.TokenKind) token.Token {
	return token.Token{
		Kind:   kind,
		Lexeme: l.source[l.start:l.current],
		Pos:    l.startPos,
	}
}

func (l *Lexer) makeTokenLexeme(kind token.TokenKind, lexeme string) token.Token {
	return token.Token{
		Kind:   kind,
		Lexeme: lexeme,
		Pos:    l.startPos,
	}
}
