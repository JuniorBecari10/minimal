package analyzer

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
	"minc/types"
	"minlib/token"
)

func (a *Analyzer) analyzeBinary(
	left, right ast.Expression,
	shallow bool,
	expectedType types.TypeData,
) (tast.Expression, tast.Expression, diagnostic.Diagnostic) {
	errReturn := func(diag diagnostic.Diagnostic) (tast.Expression, tast.Expression, diagnostic.Diagnostic) {
		return tast.Expression{}, tast.Expression{}, diag
	}

	leftTyped, diag := a.analyzeExpression(left, shallow, &expectedType); if diag != nil {
		return errReturn(diag)
	}

	rightTyped, diag := a.analyzeExpression(right, shallow, &expectedType); if diag != nil {
		return errReturn(diag)
	}

	leftCoerced, ok := coerceExpr(leftTyped, expectedType); if !ok {
		return errReturn(a.makeExpectedType(expectedType, leftTyped.Data.Type(), left.Base.Token))
	}

	rightCoerced, ok := coerceExpr(rightTyped, expectedType); if !ok {
		return errReturn(a.makeExpectedType(expectedType, rightTyped.Data.Type(), right.Base.Token))
	}

	leftTyped = leftCoerced
	rightTyped = rightCoerced

	return leftTyped, rightTyped, nil
}

func newNative(name string, globalType types.TypeData) Global {
	return Global{
		name: token.Token{ Lexeme: name },
		globalType: types.DummyType(globalType),
		immutable: true, // all native types should be immutable.
		initialized: true,
		modified: false,
		used: false,
	}
}

func typeIsConcrete(t types.TypeData) bool {
	switch t.(type) {
		// the abstract types (nil and unknown).
		case types.TypeUntypedNil, types.TypeUnknown:
			return false

		default:
			return true
	}
}

func typeIsIterable(t types.TypeData) bool {
	switch t.(type) {
		case types.TypeStr, types.TypeRange:
			return true
		default:
			return false
	}
}

func typeIsNumeric(t types.TypeData) bool {
	ok := canCoerce(t, types.TypeFloat{}) // int and float will succeed.
	return ok
}

/*
func (a *Analyzer) assertType(t, expected types.TypeData, tok token.Token) diagnostic.Diagnostic {
	if !canCoerce(t, expected) {
		return a.makeExpectedType(expected, t, tok)
	}

	return nil
}
*/

// assumes that 'iterable' has an iterable type.
func getIteratorType(iterable tast.Expression) types.TypeData {
	switch t := iterable.Data.Type().(type) {
		case types.TypeStr:
			return types.TypeChar{}

		case types.TypeRange:
			return t.Inside.Data

		default:
			panic("Internal: Type is not iterable!")
	}
}

// wraps the given expression in a CoerceExpression, if the types can be coerced, but not equal.
func coerceExpr(expr tast.Expression, convertType types.TypeData) (tast.Expression, bool) {
	// if the types are equal, there's no need to coerce.
	if expr.Data.Type() == convertType {
		return expr, true
	}

	ok := canCoerce(expr.Data.Type(), convertType)

	if !ok {
		return tast.Expression{}, ok
	}

	return tast.Expression{
		Base: expr.Base,
		Data: tast.CoerceExpression{
			Operand: expr,
			ConvertType: convertType,
		},
	}, true
}

func canCoerce(from, to types.TypeData) bool {
	t := mergeTypes(from, to)
	return t != nil
}

func mergeTypes(from, to types.TypeData) types.TypeData {
	// if types are equal, they can be merged.
	if from == to {
        return from
	}

	switch from.(type) {
		// int -> float
		case types.TypeInt: {
			if _, ok := to.(types.TypeFloat); ok {
				return to
			}
		}

		// nil -> any?
		case types.TypeUntypedNil: {
			if _, ok := to.(types.TypeOptional); ok {
				return to
			}
		}

		// unknown -> any
		case types.TypeUnknown: {
			return to
		}
	}

	return nil
}

func (a *Analyzer) newScope() {
	a.scopeDepth++
}

func (a *Analyzer) endScope() {
	a.scopeDepth--

	// remove all variables from the current scope
	for i := len(a.locals) - 1; i >= 0; i-- {
		if a.locals[i].depth <= a.scopeDepth {
			// this means all above this variable belongs to the removed scope.
			a.locals = a.locals[:i + 1]
			return
		}

		// show warnings.

		// var + not modified
		if !a.locals[i].immutable && !a.locals[i].modified {
			a.makeWarnNotModified(a.locals[i].name).PrintDiagnostic()
		} else if !a.locals[i].used {
			// ignore main function
			if a.locals[i].depth == 1 && a.locals[i].name.Lexeme == "main" {
				continue
			}

			a.makeWarnNotUsed(a.locals[i].name).PrintDiagnostic()
		}
	}

	a.locals = []Local{}
}

func (a *Analyzer) endTopLevel() AnalyzerResult {
	res := RES_OK
	foundMain := false

	for _, global := range a.globals {
		if global.name.Lexeme == "main" {
			foundMain = true

			mainType := types.TypeFunction{
				Parameters: []types.Type{}, // no parameters
				Return: types.DummyType(types.TypeVoid{}), // void
			}

			if !canCoerce(global.globalType.Data, mainType) {
				a.makeExpectedType(mainType, global.globalType.Data, global.name).PrintDiagnostic()
				res = RES_ERROR
			}
		}
	}

	if !foundMain {
		a.makeExpectedMain().PrintDiagnostic()
		res = RES_ERROR
	}

	return res
}

func (a *Analyzer) addVariable(name token.Token, varType types.Type, immutable bool) {
	if a.scopeDepth == 0 {
		// don't add; mark the variable as initialized.
		for i, global := range a.globals {
			if global.name.Lexeme == name.Lexeme {
				a.globals[i].initialized = true
				return
			}
		}
	} else {
		// add.
		a.locals = append(a.locals, Local{
			name: name,
			localType: varType,
			immutable: immutable,
			depth: a.scopeDepth,
			modified: false,
			used: false,
		})
	}
}
