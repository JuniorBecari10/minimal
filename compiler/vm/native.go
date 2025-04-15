package vm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"vm-go/value"
)

func (v *VM) includeNativeFns() {
    appendNativeFn(&v.globals, 1, nativePrint)
    appendNativeFn(&v.globals, 1, nativePrintln)
    appendNativeFn(&v.globals, 1, nativeInput)
    appendNativeFn(&v.globals, 0, nativeTime)
    appendNativeFn(&v.globals, 1, nativeStr)
    appendNativeFn(&v.globals, 1, nativeNum)
    appendNativeFn(&v.globals, 1, nativeType)
}

func appendNativeFn(list *[]value.Value, arity int, fn value.NativeFn) {
    *list = append(*list, value.ValueNativeFn{
        Arity: arity,
        Fn: fn,
    })
}

// ---

// TODO: use format string "%.10g" without printing {}
func nativePrint(args []value.Value) value.Value {
	fmt.Print(args[0].String())
	return value.ValueVoid{}
}

func nativePrintln(args []value.Value) value.Value {
    fmt.Println(args[0].String())
	return value.ValueVoid{}
}

func nativeInput(args []value.Value) value.Value {
	prompt, ok := args[0].(value.ValueString)
    if !ok {
        return value.ValueNil{}
    }

    fmt.Print(prompt.Value)

    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')

    if err != nil {
        return value.ValueNil{}
    }

    // Trim the newline character from the input
    input = strings.TrimSpace(input)
    return value.ValueString{Value: input}
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

func nativeType(args []value.Value) value.Value {
	return value.ValueString{Value: args[0].Type()}
}
