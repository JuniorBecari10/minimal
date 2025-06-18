package analyzer

import (
	"fmt"
	"minc/diagnostic"
	"minc/types"
	"minlib/token"
)

func (a *Analyzer) makeWarnUnreachable(tok token.Token) diagnostic.Diagnostic {
	return a.makeWarningDiagnostic(
		"Unreachable code.",
		diagnostic.WARN_UNREACHABLE,
		tok,
	)
}

func (a *Analyzer) makeExpectedConcreteType(t types.Type) diagnostic.Diagnostic {
	return a.makeHelpDiagnostic(
		fmt.Sprintf("Expected concrete type, but '%s' is abstract.", t),
		[]string{ "Please annotate this with a concrete type." },
		t.Token,
	)
}

func (a *Analyzer) makeExpectedTypeAnnotation(token token.Token) diagnostic.Diagnostic {
	return a.makeHelpDiagnostic(
		"Expected type annotation here.",
		[]string{ "Please annotate this with an actual type." },
		token,
	)
}

func (a *Analyzer) makeExpectedIterableType(t types.Type) diagnostic.Diagnostic {
	return a.makeHelpDiagnostic(
		fmt.Sprintf("Expected iterable type, but '%s' is not.", t.Token.FormatError()),
		[]string{ "Please annotate this with an actual type." },
		t.Token,
	)
}

func (a *Analyzer) makeExpectedType(expected, got types.TypeData, tok token.Token) diagnostic.Diagnostic {
	return a.makeDiagnostic(
		fmt.Sprintf("Expected type '%s' (or one that can be coerced to it), but got '%s', which cannot.", expected.String(), got.String()),
		tok,
	)
}

func (a *Analyzer) makeBreakContinueOutsideLoop(tok token.Token) diagnostic.Diagnostic {
	return a.makeDiagnostic(
		fmt.Sprintf("'%s' statement is outside of a loop.", tok.Lexeme),
		tok,
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

func (a *Analyzer) makeWarningDiagnostic(message string, warnType diagnostic.WarningType, token token.Token) diagnostic.WarningDiagnostic {
	return diagnostic.WarningDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: message,
			Span: diagnostic.Span{
				Pos: token.Pos,
				Length: len(token.Lexeme),
			},
			FileData: a.fileData,
		},
		WarnType: warnType,
	}
}
