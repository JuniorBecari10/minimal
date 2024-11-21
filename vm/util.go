package vm

import (
	"reflect"
	"vm-go/value"
)

func typesEqual(a, b value.Value) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

func valuesEqual(a, b value.Value) bool {
	switch val := a.(type) {
		case value.ValueNumber:
			return val.Value == b.(value.ValueNumber).Value
		
		case value.ValueBool:
			return val.Value == b.(value.ValueBool).Value
	}

	return false
}

func isNumber(v value.Value) bool {
	_, ok := v.(value.ValueNumber)
	return ok
}

func isBool(v value.Value) bool {
	_, ok := v.(value.ValueBool)
	return ok
}

func xor(a, b bool) bool {
    return (a || b) && !(a && b)
}