package ast

import (
	"minc/types"
	"minlib/token"
)

type Ast = []Statement

type AstBase struct {
	Token token.Token
}

type Parameter struct {
	Name token.Token
	Type *types.Type // optional
}

type Field struct {
	Name token.Token
	Type types.Type
}
