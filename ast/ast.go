package ast

import "vm-go/token"

type AstBase struct {
	Pos token.Position
}

// later we'll add types
type Parameter struct {
	Name token.Token
}
