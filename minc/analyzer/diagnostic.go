package analyzer

import (
	"fmt"
	"minc/diagnostic"
	"minc/types"
	"minlib/token"
)

func (a *Analyzer) makeExpectedConcreteType(t types.Type, token token.Token) diagnostic.Diagnostic {
	return a.makeHelpDiagnostic(
		fmt.Sprintf("Expected concrete type, but '%s' is abstract.", t),
		[]string{ "Please annotate this with a concrete type." },
		token,
	)
}

func (a *Analyzer) makeExpectedTypeAnnotation(token token.Token) diagnostic.Diagnostic {
	return a.makeHelpDiagnostic(
		"Expected type annotation here.",
		[]string{ "Please annotate this with an actual type." },
		token,
	)
}

func (a *Analyzer) makeDiagnostic(message string, token token.Token) diagnostic.SimpleDiagnostic {
	return diagnostic.SimpleDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: message,
			Span: diagnostic.Span{
				Pos: token.Pos,
				Length: len(token.Lexeme),
			},
			FileData: a.fileData,
		},
	}
}

func (a *Analyzer) makeHelpDiagnostic(message string, help []string, token token.Token) diagnostic.HelpDiagnostic {
	return diagnostic.HelpDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: message,
			Span: diagnostic.Span{
				Pos: token.Pos,
				Length: len(token.Lexeme),
			},
			FileData: a.fileData,
		},
		Help: help,
	}
}
