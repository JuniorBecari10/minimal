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
	}
}

func typeIsConcrete(t types.TypeData) bool {
	switch t.(type) {
		case types.TypeUntypedNil, types.TypeUnknown, types.TypeNever:
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
// this returns t twice in two expressions in order for the caller to be able to call this in an assignment switch
// if t, ok := tryCoercing(.., ..); ok { .. }
func tryCoercing(from, to types.TypeData) (types.TypeData, bool) {
	t := mergeTypes(from, to)
	return t, t != nil
}

func mergeTypes(from, to types.TypeData) types.TypeData {
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
