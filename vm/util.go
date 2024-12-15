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

func (v *VM) captureUpvalue(localsIndex int, index int) *value.Upvalue {
	// Search for an existing upvalue for that variable.
	for i, upvalue := range v.openUpvalues {
		// If an upvalue to this location already exists, return it.
		if upvalue.LocalsIndex == localsIndex && upvalue.Index == index {
			return v.openUpvalues[i]
		}
	}

	// Otherwise, create a new upvalue, and return a reference to it.
	up := value.Upvalue{
		LocalsIndex: localsIndex,
		Index: index,
		IsClosed: false,
	}

	v.openUpvalues = append(v.openUpvalues, &up)
	return v.openUpvalues[len(v.openUpvalues)-1]
}

func (v *VM) closeUpvalue(localsIndex int, index int) {
	for i, upvalue := range v.openUpvalues {
		if upvalue.LocalsIndex == localsIndex && upvalue.Index == index {
			upvalue := value.Upvalue{
				ClosedValue: v.getUpvalueValue(upvalue),
				IsClosed: true,
			}
	
			*v.openUpvalues[i] = upvalue

			// remove the upvalue from the list, as it's not open anymore.
			// this isn't a concurrency problem because after changing the list,
			// we'll return from this function.
			util.Remove(v.openUpvalues, i)
			return
		}
	}
}

// ---

func (v *VM) getUpvalueValue(upvalue *value.Upvalue) value.Value {
	if upvalue.IsClosed {
		return upvalue.ClosedValue
	} else {
		return v.callStack[upvalue.LocalsIndex].locals[upvalue.Index]
	}
}

func (v *VM) setUpvalueValue(upvalue *value.Upvalue, val value.Value) {
	if upvalue.IsClosed {
		upvalue.ClosedValue = val
	} else {
		v.callStack[upvalue.LocalsIndex].locals[upvalue.Index] = val
	}
}
