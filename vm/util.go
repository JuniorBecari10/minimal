package vm

import (
	"reflect"
	"vm-go/util"
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

func (v *VM) getByte() byte {
	res := v.currentChunk.Code[v.ip]
	v.ip += 1

	return res
}

func (v *VM) getInt() int {
	res, _ := util.BytesToInt(v.currentChunk.Code[v.ip:v.ip + 4])
	v.ip += 4

	return res
}

func isNumber(v value.Value) bool {
	_, ok := v.(value.ValueNumber)
	return ok
}

func isBool(v value.Value) bool {
	_, ok := v.(value.ValueBool)
	return ok
}

func isString(v value.Value) bool {
	_, ok := v.(value.ValueString)
	return ok
}

func isClosure(v value.Value) bool {
	_, ok := v.(value.ValueClosure)
	return ok
}

func isNativeFunction(v value.Value) bool {
	_, ok := v.(value.ValueNativeFn)
	return ok
}

// ---

func (v *VM) captureUpvalue(slot *value.Value) *value.Upvalue {
	// Search for an existing upvalue for that variable.
	for i, upvalue := range v.upvalues {
		// If an upvalue to this location already exists, return it.
		if upvalue.Location == slot {
			return &v.upvalues[i]
		}
	}

	// Otherwise, create a new upvalue, and return a reference to it.
	v.upvalues = append(v.upvalues, value.Upvalue{
		Location: slot,
		IsClosed: false,
	})

	return &v.upvalues[len(v.upvalues)-1]
}

func (v *VM) closeUpvalue(slot *value.Value) {
	for i := range v.upvalues {
		if v.upvalues[i].Location == slot {
			upvalue := value.Upvalue{
				ClosedValue: *slot,
				IsClosed: true,
			}
	
			upvalue.Location = &upvalue.ClosedValue
			v.upvalues[i] = upvalue

			return
		}
	}
}
