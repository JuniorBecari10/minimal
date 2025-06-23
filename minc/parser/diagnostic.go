package parser

import (
	"fmt"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) makeWarnSemicolonNotNeeded() diagnostic.WarningHelpDiagnostic {
	return p.makeWarnHelpDiagnostic(
		"This semicolon is not necessary.",
		[]string{
			"",
		},
	)
}

func (p *Parser) makeExpectedTokenDiagnostic(expected token.TokenKind) diagnostic.SimpleDiagnostic {
	return p.makeDiagnostic(
		fmt.Sprintf("Expected %s after %s, but got %s.",
			expected, p.previous.FormatError(), p.current.FormatError()))
}

func (p *Parser) makeStatementsNotAllowedDiagnostic() diagnostic.SimpleDiagnostic {
	return p.makeDiagnostic("Statements are not allowed at top-level.")
}

func (p *Parser) makeInvalidTypeArgumentsLengthDiagnostic(expected, got int) diagnostic.SimpleDiagnostic {
	return p.makeDiagnostic(fmt.Sprintf("Expected %d type arguments, but got %d.", expected, got))
}

func (p *Parser) makeExpectedExpressionDiagnostic() diagnostic.SimpleDiagnostic {
	return p.makeDiagnostic(fmt.Sprintf("Expected expression after %s, but got %s.", p.previous.FormatError(), p.current.FormatError()))
}

func (p *Parser) makeExpectedTypeDiagnostic() diagnostic.SimpleDiagnostic {
	return p.makeDiagnostic(fmt.Sprintf("Expected type after %s, but got %s.", p.previous.FormatError(), p.current.FormatError()))
}

func (p *Parser) makeInvalidAssignmentTargetDiagnostic() diagnostic.SimpleDiagnostic {
	return p.makeDiagnostic("Invalid assignment target.")
}

func (p *Parser) makeDiagnostic(message string) diagnostic.SimpleDiagnostic {
	return diagnostic.SimpleDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: message,
			Span: diagnostic.Span{
				Pos: p.current.Pos,
				Length: p.current.Length(),
			},
			FileData: p.fileData,
		},
	}
}

func (p *Parser) makeWarnHelpDiagnostic(message string, help []string) diagnostic.WarningHelpDiagnostic {
	return diagnostic.WarningHelpDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: message,
			Span: diagnostic.Span{
				Pos: p.current.Pos,
				Length: p.current.Length(),
			},
			FileData: p.fileData,
		},
		Help: help,
	}
}
