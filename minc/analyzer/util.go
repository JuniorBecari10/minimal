package analyzer

import (
	"minc/tast"
	"minc/types"
	"minlib/token"
)

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

// this also returns true if the types are equal, or the types inside them can coerce into the other too.
// this returns t twice in two expressions in order for the caller to be able to call this in an assignment if
// if t, ok := tryCoercing(.., ..); ok { .. }
func tryCoercing(from, to types.TypeData) (types.TypeData, bool) {
	t := mergeTypes(from, to)
	return t, t != nil
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

		// var + not modified
		if !a.locals[i].immutable && !a.locals[i].modified {
			a.makeWarnNotModified(a.locals[i].name).PrintDiagnostic()
		} else if !a.locals[i].used {
			a.makeWarnNotUsed(a.locals[i].name).PrintDiagnostic()
		}
	}

	a.locals = []Local{}
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
