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

type ValueInt struct {
	Value int32
}

type ValueFloat struct {
	Value float64
}

type ValueString struct {
	Value string
}

type ValueChar struct {
	Value uint8
}

type ValueBool struct {
	Value bool
}

type ValueNil struct {}
type ValueVoid struct {}

type ValueFunction struct {
	Arity uint32
	Chunk Chunk
	Name *string // optional
}

type ValueClosure struct {
	Fn *ValueFunction
	Upvalues []Upvalue
}

type ValueNativeFn struct {
	Arity uint32
	Fn NativeFn
}

type ValueRange struct {
    Start float64
    End   float64
    Step  float64

    Inclusive bool
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

type ValueBoundMethod struct {
	Receiver Value
	Method ValueClosure
}

// ---

func (x ValueInt) String() string {
	return fmt.Sprintf("%d", x.Value)
}

func (x ValueFloat) String() string {
	return fmt.Sprintf("%.10g", x.Value)
}

func (x ValueString) String() string { return x.Value }
func (x ValueChar) String() string { return fmt.Sprintf("%c", x.Value) }

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

func (x ValueRange) String() string {
    if x.Inclusive {
        return fmt.Sprintf("%.10g..=%.10g:%.10g", x.Start, x.End, x.Step)
    } else {
        return fmt.Sprintf("%.10g..%.10g:%.10g", x.Start, x.End, x.Step)
    }
}

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

func (x ValueBoundMethod) String() string {
	if x.Method.Fn.Name == nil {
		return "<bound fn>"
	} else {
		return fmt.Sprintf("<bound fn %s>", *x.Method.Fn.Name)
	}
}

// ---

func (x ValueInt) Type() string { return "int" }
func (x ValueFloat) Type() string { return "float" }
func (x ValueString) Type() string { return "str" }
func (x ValueChar) Type () string { return "char" }
func (x ValueBool) Type() string { return "bool" }

func (x ValueNil) Type() string { return "nil" }
func (x ValueVoid) Type() string { return "void" }

func (x ValueFunction) Type() string { return "fn" }
func (x ValueNativeFn) Type() string { return "native fn" }
func (x ValueClosure) Type() string { return "fn" }

func (x ValueRange) Type() string { return "range" }
func (x ValueRecord) Type() string { return "record" }
func (x ValueInstance) Type() string { return x.Record.Name }

func (x ValueBoundMethod) Type() string { return "bound fn" }
