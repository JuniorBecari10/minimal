package vm

import (
	"fmt"
	"strconv"
	"time"
	"vm-go/value"
)

func (v *VM) includeNativeFns() {
	v.globals = append(v.globals, value.ValueNativeFn{
		Arity: 1,
		Fn: nativePrint,
	})

	v.globals = append(v.globals, value.ValueNativeFn{
		Arity: 1,
		Fn: nativePrintln,
	})

	v.globals = append(v.globals, value.ValueNativeFn{
		Arity: 0,
		Fn: nativeTime,
	})

	v.globals = append(v.globals, value.ValueNativeFn{
		Arity: 1,
		Fn: nativeStr,
	})

	v.globals = append(v.globals, value.ValueNativeFn{
		Arity: 1,
		Fn: nativeNum,
	})
}

// ---

func nativePrint(args []value.Value) value.Value {
	fmt.Print(args[0])
	return value.ValueVoid{}
}

func nativePrintln(args []value.Value) value.Value {
	fmt.Println(args[0])
	return value.ValueVoid{}
}

func nativeTime(_ []value.Value) value.Value {
	return value.ValueNumber{ Value: float64(time.Now().UnixMilli()) }
}

func nativeStr(args []value.Value) value.Value {
	return value.ValueString{ Value: args[0].String() }
}

func nativeNum(args []value.Value) value.Value {
	argStr, ok := args[0].(value.ValueString)

	if !ok {
		return value.ValueNil{}
	}

	asNum, err := strconv.ParseFloat(argStr.Value, 64)

	if err != nil {
		return value.ValueNil{}
	}

	return value.ValueNumber{ Value: asNum }
}
