package analyzer

import (
	"minc/types"
	"minlib/token"
)

// TODO: add flags canBeMutable (for functions and refine the error message)
type Local struct {
	name token.Token
	localType types.Type

	immutable bool
	depth uint32
	modified bool
	used bool
}

func (v *Local) Name() token.Token {
	return v.name
}

func (v *Local) Type() types.Type {
	return v.localType
}

func (v *Local) IsImmutable() bool {
	return v.immutable
}

func (v *Local) IsModified() bool {
	return v.modified
}

func (v *Local) Used() bool {
	return v.used
}

func (v *Local) MarkModified() {
	v.modified = true
}

func (v *Local) MarkUsed() {
	v.used = true
}

type Global struct {
	name token.Token
	globalType types.Type

	immutable bool
	initialized bool // to check if it has already been declared or just hoisted
	modified bool
	used bool
}

func (v *Global) Name() token.Token {
	return v.name
}

func (v *Global) Type() types.Type {
	return v.globalType
}

func (v *Global) IsImmutable() bool {
	return v.immutable
}

func (v *Global) IsModified() bool {
	return v.modified
}

func (v *Global) Used() bool {
	return v.used
}

func (v *Global) MarkModified() {
	v.modified = true
}

func (v *Global) MarkUsed() {
	v.used = true
}
