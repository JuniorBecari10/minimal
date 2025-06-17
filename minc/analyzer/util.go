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
func typeCanCoerceTo(from, to types.TypeData) bool {

}
