package ast

import "vm-go/token"

type AstBase struct {
	Pos token.Position
	Length int
}

// later we'll add types
type Parameter struct {
	Name token.Token
}

type Field struct {
	Name token.Token
}
