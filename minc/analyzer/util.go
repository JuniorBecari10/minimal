package analyzer

import (
	"minc/types"
	"minlib/token"
)

func newNative(name string, globalType types.Type) Global {
	return Global{
		name: token.Token{ Lexeme: name },
		globalType: globalType,
		immutable: true, // all native types should be immutable.
		initialized: true,
	}
}

func typeIsConcrete(t types.Type) bool {
	switch t.(type) {
		// the only two abstract types
		case types.TypeUntypedNil, types.TypeUnknown:
			return false

		default:
			return true
	}
}
