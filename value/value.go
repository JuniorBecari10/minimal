package value

import (
	"bytes"
	"fmt"
)

type NativeFn = func(args []Value) Value

type Value interface {
	String() string
	Type() string
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
	Name string
	FieldNames []string
	Methods []ValueClosure
}

type ValueInstance struct {
	Fields []Value
	Record *ValueRecord
}

func (in *ValueInstance) GetProperty(name string) (Value, bool) {
	for i, value := range in.Fields {
		if in.Record.FieldNames[i] == name {
			return value, true
		}
	}

	return ValueNil{}, false
}

func (in *ValueInstance) SetProperty(name string, value Value) bool {
	for i := range in.Fields {
		if in.Record.FieldNames[i] == name {
			in.Fields[i] = value
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
		res.WriteString(fmt.Sprintf("%s: %s", x.Record.FieldNames[i], field.String()))

		// Add a comma and space if it isn't the last field.
		if i < len(x.Fields) - 1 {
			res.WriteString(", ")
		}
	}

	res.WriteString(")")
	return res.String()
}

// ---

func (x ValueNumber) Type() string { return "num" }
func (x ValueString) Type() string { return "str" }
func (x ValueBool) Type() string { return "bool" }

func (x ValueNil) Type() string { return "nil" }
func (x ValueVoid) Type() string { return "void" }

func (x ValueFunction) Type() string { return "fn" }
func (x ValueNativeFn) Type() string { return "native fn" }
func (x ValueClosure) Type() string { return "fn" }

func (x ValueRecord) Type() string { return "record" }
func (x ValueInstance) Type() string { return x.Record.Name }
