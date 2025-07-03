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
		tok,
	)
}

func (a *Analyzer) makeWarnRedundantSemicolon(semicolon token.Token) diagnostic.Diagnostic {
	return a.makeWarningHelpDiagnostic(
		"This semicolon is redundant.",
		[]string{
			"Regardless of its presence, the expression will be returned",
			"from the block, since there is no other statement there.",
			"",
			"If you don't want to return the expression, wrap it in 'void()'",
			"and don't put a semicolon after it.",
		},
		semicolon,
	)
}

func (a *Analyzer) makeExpectedSemicolon(tok token.Token) diagnostic.Diagnostic {
	return a.makeHelpDiagnostic(
		"Expected a semicolon after this expression.",
		[]string{
			"The expression won't be returned regardless of the semicolon,",
			"since there is more than one statement in this block;",
			"but you need to put it there to terminate this last statement.",
			"",
			"If you want to return this expression from this block, use the 'out' statement.",
		},
		tok,
	)
}

func (a *Analyzer) makeWarnNotModified(name token.Token) diagnostic.Diagnostic {
	return a.makeWarningHelpDiagnostic(
		fmt.Sprintf("'%s' is mutable but is not modified.", name.Lexeme),
		[]string{ "Declare it with 'let' if you don't plan to modify it." },
		name,
	)
}

func (a *Analyzer) makeWarnNotUsed(name token.Token) diagnostic.Diagnostic {
	return a.makeWarningDiagnostic(
		fmt.Sprintf("'%s' is declared but is not used.", name.Lexeme),
		name,
	)
}

func (a *Analyzer) makeWarnUnusedValueExprStmt(gotType types.TypeData, tok token.Token) diagnostic.Diagnostic {
	return a.makeWarningHelpDiagnostic(
		fmt.Sprintf("Unused value of expression of type '%s'.", gotType.String()),
		[]string{ "Wrap it with 'void()' to intentionally discard it." },
		tok,
	)
}

func (a *Analyzer) makeNameNotDefined(name token.Token) diagnostic.Diagnostic {
	return a.makeDiagnostic(
		fmt.Sprintf("'%s' is not defined in scope.", name.Lexeme),
		name,
	)
}

func (a *Analyzer) makeCannotModifyImmutable(name token.Token) diagnostic.Diagnostic {
	return a.makeHelpDiagnostic(
		fmt.Sprintf("Cannot modify the value of '%s', as it's immutable.", name.Lexeme),
		[]string { "If you plan to modify its value, declare it using 'var' instead." },
		name,
	)
}

func (a *Analyzer) makeExpectedConcreteType(t types.Type) diagnostic.Diagnostic {
	return a.makeHelpDiagnostic(
		fmt.Sprintf("Expected concrete type, but '%s' is abstract.", t.Data),
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
	return a.makeDiagnostic(
		fmt.Sprintf("Expected iterable type, but '%s' is not.", t.Token.FormatError()),
		t.Token,
	)
}

func (a *Analyzer) makeExpectedCallableType(t types.TypeData, tok token.Token) diagnostic.Diagnostic {
	return a.makeDiagnostic(
		fmt.Sprintf("Expected callable type, but '%s' is not.", t.String()),
		tok,
	)
}

func (a *Analyzer) makeIncompatibleTypes(left, right types.TypeData, tok token.Token) diagnostic.Diagnostic {
	return a.makeDiagnostic(
		fmt.Sprintf("Incompatible types: '%s' and '%s'.", left.String(), right.String()),
		tok,
	)
}

func (a *Analyzer) makeExpectedType(expected, got types.TypeData, tok token.Token) diagnostic.Diagnostic {
	return a.makeDiagnostic(
		fmt.Sprintf("Expected type '%s', but got '%s', which cannot be coerced to it.", expected.String(), got.String()),
		tok,
	)
}

func (a *Analyzer) makeExpectedReturnType(expected, got types.TypeData, tok token.Token, annotated bool) diagnostic.Diagnostic {
	if annotated {
		return a.makeDiagnostic(
			fmt.Sprintf("Expected return type to be '%s', but got '%s', which cannot be coerced to it.", expected.String(), got.String()),
			tok,
		)
	} else {
		return a.makeHelpDiagnostic(
			fmt.Sprintf("Expected return type to be '%s', but got '%s', which cannot be coerced to it.", expected.String(), got.String()),
			[]string { "The returned type is not annotated, so it's set to be 'void' by default." },
			tok,
		)
	}
}

func (a *Analyzer) makeExpectedArity(expected, got int, tok token.Token) diagnostic.Diagnostic {
	return a.makeDiagnostic(
		fmt.Sprintf("Expected %d argument(s), but got %d instead.", expected, got),
		tok,
	)
}

func (a *Analyzer) makeExpectedTypeArity(expected, got int, tok token.Token) diagnostic.Diagnostic {
	return a.makeDiagnostic(
		fmt.Sprintf("Expected %d type argument(s), but got %d instead.", expected, got),
		tok,
	)
}

func (a *Analyzer) makeExpectedMain() diagnostic.Diagnostic {
	return a.makeHeadDiagnostic(
		"File must have a main function.",
		[]string { "Declare it using the following signature: 'fn main()'." },
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
				Length: token.Length(),
			},
			FileData: a.fileData,
		},
	}
}

func (a *Analyzer) makeHeadDiagnostic(message string, help []string) diagnostic.HeadHelpDiagnostic {
	return diagnostic.HeadHelpDiagnostic{
		Message: message,
		FileData: a.fileData,
		Help: help,
	}
}

func (a *Analyzer) makeHelpDiagnostic(message string, help []string, token token.Token) diagnostic.HelpDiagnostic {
	return diagnostic.HelpDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: message,
			Span: diagnostic.Span{
				Pos: token.Pos,
				Length: token.Length(),
			},
			FileData: a.fileData,
		},
		Help: help,
	}
}

func (a *Analyzer) makeWarningDiagnostic(message string, token token.Token) diagnostic.WarningDiagnostic {
	return diagnostic.WarningDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: message,
			Span: diagnostic.Span{
				Pos: token.Pos,
				Length: token.Length(),
			},
			FileData: a.fileData,
		},
	}
}

func (a *Analyzer) makeWarningHelpDiagnostic(
	message string,
	help []string,
	token token.Token,
) diagnostic.WarningHelpDiagnostic {
	return diagnostic.WarningHelpDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: message,
			Span: diagnostic.Span{
				Pos: token.Pos,
				Length: token.Length(),
			},
			FileData: a.fileData,
		},
		Help: help,
	}
}
