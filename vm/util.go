package vm

import (
	"reflect"
	"vm-go/value"
)

func checkTypes(a, b value.Value) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}