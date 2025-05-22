package types

import (
	"fmt"
	"strings"
)

type Type interface {
	String() string
}

type TypeInt struct { }
type TypeFloat struct { }
type TypeStr struct { }
type TypeChar struct { }
type TypeBool struct { }

// used internally to represent nil in the AST, but it should never go to a variable.
type TypeUntypedNil struct { }

// for future use in generics and types intentionally omitted for inference.
type TypeUnknown struct { }

type TypeVoid struct { }

type TypeOptional struct {
	Inside Type
}

type TypeFunction struct {
	Parameters []Type
	Return Type
}

type TypeRange struct {
	Inside Type
}

// For records and others, like enums
type TypeUserDefined struct {
	Name string
}

// ---

func (x TypeInt) String() string { return "int" }
func (x TypeFloat) String() string { return "float" }
func (x TypeStr) String() string { return "str" }
func (x TypeChar) String() string { return "char" }
func (x TypeBool) String() string { return "bool" }
func (x TypeUntypedNil) String() string { return "untyped nil" }
func (x TypeUnknown) String() string { return "unknown" }
func (x TypeVoid) String() string { return "void" }
func (x TypeOptional) String() string { return fmt.Sprintf("%s?", x.Inside.String()) }

// fn (int, int): int
func (x TypeFunction) String() string {
	builder := strings.Builder{}
	builder.WriteString("fn (")

	for i, ty := range x.Parameters {
		builder.WriteString(ty.String())

		if i < len(x.Parameters) - 1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString("): ")
	builder.WriteString(x.Return.String())

	return builder.String()
}

func (x TypeRange) String() string { return fmt.Sprintf("range<%s>", x.Inside.String()) }
func (x TypeUserDefined) String() string { return x.Name }
