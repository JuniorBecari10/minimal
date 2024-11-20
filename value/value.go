package value

import "fmt"

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

// ---

func (x ValueNumber) String() string { return fmt.Sprintf("%.2f", x.Value) }
func (x ValueString) String() string { return x.Value }
func (x ValueBool) String() string {
	if x.Value {
		return "true"
	} else {
		return "false"
	}
}
