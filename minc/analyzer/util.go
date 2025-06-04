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
