package value

import (
	"fmt"
	"vm-go/token"
)

// Another alias for Chunk, because of import cycle
type Chunk = struct {
	Code      []byte
	Constants []Value

	Positions []token.Position
}

type Value interface {
	String() string
}

// ---

type ValueNumber struct {
	Value float64
}

type ValueString struct {
	Value string
}

type ValueBool struct {
	Value bool
}

type ValueNil struct {}

type ValueFunction struct {
	Arity int
	Chunk Chunk
	Name string
}

// ---

func (x ValueNumber) String() string { return fmt.Sprintf("%.2f", x.Value) }
func (x ValueString) String() string { return x.Value }
func (x ValueBool) String() string { return fmt.Sprintf("%t", x.Value) }
func (x ValueNil) String() string { return "nil" }
func (x ValueFunction) String() string { return fmt.Sprintf("<fn %s>", x.Name) }
