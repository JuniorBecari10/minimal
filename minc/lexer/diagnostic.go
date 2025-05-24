package lexer

import (
	"fmt"
	"minc/diagnostic"
)

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
