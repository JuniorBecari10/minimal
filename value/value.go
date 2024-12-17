package value

import (
	"fmt"
	"vm-go/token"
)

type NativeFn = func(args []Value) Value

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
type ValueVoid struct {}

type ValueFunction struct {
	Arity int
	Chunk Chunk
	Name *string // optional
}

type ValueClosure struct {
	Fn *ValueFunction
	Upvalues []*Upvalue
}

type ValueNativeFn struct {
	Arity int
	Fn NativeFn
}

type ValueRecord struct {
	Arity int
}

// ---

func (x ValueNumber) String() string {
	// if it's an integer
	if x.Value == float64(int64(x.Value)) {
		return fmt.Sprintf("%.0f", x.Value)
	}
	
	return fmt.Sprintf("%.2f", x.Value)
}

func (x ValueString) String() string { return x.Value }
func (x ValueBool) String() string { return fmt.Sprintf("%t", x.Value) }
func (x ValueNil) String() string { return "nil" }
func (x ValueVoid) String() string { return "void" }

func (x ValueFunction) String() string {
	if x.Name == nil {
		return "<fn>"
	} else {
		return fmt.Sprintf("<fn %s>", *x.Name)
	}
}

func (x ValueNativeFn) String() string { return "<native fn>" }
func (x ValueClosure) String() string { return x.Fn.String() }
func (x ValueRecord) String() string { return "<record>" }
