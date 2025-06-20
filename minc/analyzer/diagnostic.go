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

func (a *Analyzer) makeWarnNotModified(name token.Token) diagnostic.Diagnostic {
	return a.makeWarningHelpDiagnostic(
		fmt.Sprintf("'%s' is mutable but is not modified.", name.Lexeme),
		diagnostic.WARN_UNREACHABLE,
		[]string{ "Declare it with 'let' if you don't plan to modify it." },
		name,
	)
}

func (a *Analyzer) makeWarnNotUsed(name token.Token) diagnostic.Diagnostic {
	return a.makeWarningDiagnostic(
		fmt.Sprintf("'%s' is declared but is not used.", name.Lexeme),
		diagnostic.WARN_UNREACHABLE,
		name,
	)
}

func (a *Analyzer) makeWarnUnusedValueExprStmt(gotType types.TypeData, tok token.Token) diagnostic.Diagnostic {
	return a.makeWarningHelpDiagnostic(
		fmt.Sprintf("Unused value of expression of type '%s'.", gotType.String()),
		diagnostic.WARN_UNREACHABLE,
		[]string{ "Wrap it with 'void()' to intentionally discard it." },
		tok,
	)
}

func (a *Analyzer) makeExpectedConcreteType(t types.Type) diagnostic.Diagnostic {
	return a.makeHelpDiagnostic(
		fmt.Sprintf("Expected concrete type, but '%s' is abstract.", t.Token.FormatError()),
		[]string{ "Please annotate this with a concrete type." },
		t.Token,
	)
}

func (a *Analyzer) makeTypeAnnotationsNeeded(token token.Token) diagnostic.Diagnostic {
	return a.makeHelpDiagnostic(
		"Type annotations needed here.",
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

func (a *Analyzer) makeWarningHelpDiagnostic(
	message string,
	warnType diagnostic.WarningType,
	help []string,
	token token.Token,
) diagnostic.WarningHelpDiagnostic {
	return diagnostic.WarningHelpDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: message,
			Span: diagnostic.Span{
				Pos: token.Pos,
				Length: len(token.Lexeme),
			},
			FileData: a.fileData,
		},
		WarnType: warnType,
		Help: help,
	}
}
