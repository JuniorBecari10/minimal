package tast

import (
	"minc/types"
	"minlib/token"
)

type Variable interface {
	Name() token.Token
	Type() types.Type
	IsImmutable() bool
	IsModified() bool
	Used() bool

	MarkModified()
	MarkUsed()
}
