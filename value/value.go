package value

import (
	"bytes"
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
	FieldNames []string
	Name string
}

type ValueInstance struct {
	Fields []Field
	Record *ValueRecord
}

func (in *ValueInstance) GetProperty(name string) (Value, bool) {
	for _, field := range in.Fields {
		if field.Name == name {
			return field.Value, true
		}
	}

	return ValueNil{}, false
}

func (in *ValueInstance) SetProperty(name string, value Value) bool {
	for i, field := range in.Fields {
		if field.Name == name {
			in.Fields[i].Value = value
			return true
		}
	}

	return false
}

// ---

func (x ValueNumber) String() string {
	// Check if it's an integer.
	// Is it cheaper to convert it to int64 and then back to float64, or
	// use math.Fmod to simulate 'x % 1 == 0' for floats?
	if x.Value == float64(int64(x.Value)) {
		return fmt.Sprintf("%.0f", x.Value)
	}
	
	return fmt.Sprintf("%.2f", x.Value)
}

func (x ValueString) String() string { return x.Value }

// '%t' is the built-in formatter for bools in Go.
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
func (x ValueRecord) String() string { return fmt.Sprintf("<record %s>", x.Name) }

func (x ValueInstance) String() string {
	res := bytes.Buffer{}
	res.WriteString(fmt.Sprintf("%s(", x.Record.Name))

	for i, field := range x.Fields {
		// Get the name from the referenced record, so there's no need to store the name of the field twice.
		res.WriteString(fmt.Sprintf("%s: %s", field.Name, field.Value.String()))

		// Add a comma and space if it isn't the last field.
		if i < len(x.Fields) - 1 {
			res.WriteString(", ")
		}
	}

	res.WriteString(")")
	return res.String()
}
