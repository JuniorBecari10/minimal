package value

import (
	"bytes"
	"fmt"
)

type RangeSetStatus int
const (
    RANGE_OK RangeSetStatus = iota
    RANGE_PROPERTY_DOESNT_EXIST
    RANGE_TYPE_ERROR
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
	Upvalues []Upvalue
}

type ValueNativeFn struct {
	Arity int
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

func (x ValueNumber) String() string {
	return fmt.Sprintf("%.10g", x.Value)
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

func (x ValueNumber) Type() string { return "num" }
func (x ValueString) Type() string { return "str" }
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

