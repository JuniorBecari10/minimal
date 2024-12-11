package vm

import (
	"fmt"
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
