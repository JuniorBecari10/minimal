package parser

import (
	"fmt"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) makeExpectedTokenDiagnostic(expected token.TokenKind) diagnostic.SimpleDiagnostic {
	return p.makeDiagnostic(
		fmt.Sprintf("Expected '%s' after '%s', but got '%s'.",
			expected, p.previous.FormatError(), p.current.FormatError()))
}

func (p *Parser) makeStatementsNotAllowedDiagnostic() diagnostic.SimpleDiagnostic {
	return p.makeDiagnostic("Statements are not allowed at top-level.")
}

func (p *Parser) makeDiagnostic(message string) diagnostic.SimpleDiagnostic {
	return diagnostic.SimpleDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: message,
			Span: diagnostic.Span{
				Pos: p.current.Pos,
				Length: len(p.current.Lexeme),
			},
			FileData: p.fileData,
		},
	}
}
