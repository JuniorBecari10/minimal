package analyzer

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
	"minc/types"
	"minlib/token"
	"reflect"
)

// if this outputs successfully, both expressions should have the same type, either coerced or not.
func (a *Analyzer) analyzeBinary(
	left, right ast.Expression,
	operator token.Token,
	shallow bool,
) (tast.Expression, tast.Expression, diagnostic.Diagnostic) {
	errReturn := func(diag diagnostic.Diagnostic) (tast.Expression, tast.Expression, diagnostic.Diagnostic) {
		return tast.Expression{}, tast.Expression{}, diag
	}

	leftTyped, diag := a.analyzeExpression(left, shallow, nil); if diag != nil {
		return errReturn(diag)
	}

	rightTyped, diag := a.analyzeExpression(right, shallow, nil); if diag != nil {
		return errReturn(diag)
	}

	merged := mergeTypes(leftTyped.Data.Type(), rightTyped.Data.Type())

	// try again, coercing to the other side.
    if merged == nil {
		merged = mergeTypes(rightTyped.Data.Type(), leftTyped.Data.Type())
	}

	// if they still can't coerce, throw an error.
	if merged == nil {
		return errReturn(a.makeIncompatibleTypes(leftTyped.Data.Type(), rightTyped.Data.Type(), operator))
	}

	return leftTyped, rightTyped, nil
}

func (a *Analyzer) analyzeLoopBody(body ast.BlockExpression) (tast.BlockExpression, diagnostic.Diagnostic) {
	block, res := a.analyzeBlockAlone(body, MODE_LOOP); if res == RES_ERROR {
		return tast.BlockExpression{}, diagnostic.HandledDiagnostic{}
	}

	if !canCoerce(block.BlockType, types.TypeVoid{}) {
		return tast.BlockExpression{}, a.makeExpectedType(types.TypeVoid{}, block.BlockType, block.Token)
	}

	return block, nil
}

// TODO: check the other types that can have abstract types inside them
func typeIsConcrete(t types.TypeData) bool {
	switch type_ := t.(type) {
		// the abstract types (nil and unknown).
		case types.TypeUntypedNil, types.TypeUnknown:
			return false

		// if the type inside the optional is abstract, the entire optional type is also abstract.
		case types.TypeOptional:
			return typeIsConcrete(type_.Inside.Data)

		// if any type inside a function type is abstract, the entire function type is also abstract.
		case types.TypeFunction: {
			// check parameters
			for _, param := range type_.Parameters {
				if !typeIsConcrete(param.Data) {
					return false
				}
			}

			return typeIsConcrete(type_.Return.Data)
		}

		case types.TypeRange: 
			return typeIsConcrete(type_.Inside.Data)

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

func canCoerce(from, to types.TypeData) bool {
	t := mergeTypes(from, to)
	return t != nil
}

// merge both types so that the type returned emcompasses all types from one and from the other,
// like T -> any returns 'any'.
// returns 'nil' if it cannot be done.
// TODO: support function types and range types.
func mergeTypes(from, to types.TypeData) types.TypeData {
	// if types are equal, they can be merged.
	// T <-> T
	if reflect.DeepEqual(from, to) {
		return from
	}

	// nil -> T?
	if _, isNil := from.(types.TypeUntypedNil); isNil {
		if _, ok := to.(types.TypeOptional); ok {
			return to
		}
	}
	// T? -> nil is not allowed.

	// unknown -> T
	if _, isUnknown := from.(types.TypeUnknown); isUnknown {
		return to
	}
	// T -> unknown is not allowed.

	// T? -> U? if T -> U, or
	// T? <-> U? if T <-> U
	fromOpt, isFromOpt := from.(types.TypeOptional)
	toOpt, isToOpt := to.(types.TypeOptional)

	if isFromOpt && isToOpt {
		inner := mergeTypes(fromOpt.Inside.Data, toOpt.Inside.Data)
		
		// cannot merge.
		if inner == nil {
			return nil
		}

		return types.TypeOptional{
			Inside: types.DummyType(inner),
		}
	}

	// T -> any
	if _, isAny := from.(types.TypeAny); isAny {
		return from
	}
	// any -> T (downgrade) is not allowed.

	// cannot merge.
	return nil
}

// returns if 
func typeNeedsUnusedWarning(t types.TypeData) bool {
	switch t.(type) {
		case types.TypeVoid, types.TypeNever:
			return false

		default:
			return true
	}
}

func (a *Analyzer) newScope() {
	a.scopeDepth++
}

func (a *Analyzer) endScope() {
	a.scopeDepth--

	if a.scopeDepth < 0 {
		a.scopeDepth = 0
	}

	// remove all variables from the current scope.
	// for this, we'd need to traverse the locals array backwards.
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

		if !global.immutable && !global.modified {
			a.makeWarnNotModified(global.name).PrintDiagnostic()
		} else if !global.used {
			if !a.canSayNotUsed(global.name.Lexeme) {
				continue
			}

			a.makeWarnNotUsed(global.name).PrintDiagnostic()
		}
	}

	if !foundMain {
		a.makeExpectedMain().PrintDiagnostic()
		res = RES_ERROR
	}

	return res
}

// assumes the name is from a global variable.
func (a *Analyzer) canSayNotUsed(name string) bool {
	// the main function, of course, won't be called by the user normally.
	if name == "main" {
		return false
	}

	// natives should also be suppressed.
	for _, native := range a.natives {
		if native.name.Lexeme == name {
			return false
		}
	}

	return true
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

func (a *Analyzer) resolveVariable(name string) (tast.Variable, AnalyzerResult) {
	// search in locals from back to forth.
	for i := len(a.locals) - 1; i >= 0; i-- {
		// get a pointer directly to not make a copy of it.
		local := &a.locals[i]

		if local.name.Lexeme == name {
			return local, RES_OK
		}
	}

	// didn't find. search in globals.
	for _, global := range a.globals {
		if global.name.Lexeme == name {
			return &global, RES_OK
		}
	}

	return nil, RES_ERROR
}
