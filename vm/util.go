package vm

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"vm-go/compiler"
	"vm-go/util"
	"vm-go/value"
)

func typesEqual(a, b value.Value) bool {
	// Bypass the check if either side is nil, to allow for nil checking.
	if isNil(a) || isNil(b) {
		return true
	}

	// Check if the types of a and b are the same
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

func valuesEqual(a, b value.Value) bool {
	switch a.(type) {
		case value.ValueNil:
			return true;
		
		case value.ValueVoid:
			return true;
		
		default:
			return reflect.DeepEqual(a, b)
	}
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

func isNil(v value.Value) bool {
	_, ok := v.(value.ValueNil)
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

func isRecord(v value.Value) bool {
	_, ok := v.(value.ValueRecord)
	return ok
}

func isBoundMethod(v value.Value) bool {
	_, ok := v.(value.ValueBoundMethod)
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

func (v *VM) closeUpvalues(localsIndex int) {
	for i := range v.callStack[localsIndex].locals {
		v.closeUpvalue(localsIndex, i)
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

// ---

func (v *VM) concatenateStrs() InterpretResult {
	right := v.pop()

	if v.stackIsEmpty() {
		v.error("Not enough stack items to perform a binary operation")
		return STATUS_STACK_EMPTY
	}

	left := v.pop()

	if !isString(left) || !isString(right) {
		v.error(
			fmt.Sprintf(
				"Operands must be strings when concatenating. (left: '%s' (type '%s'), right: '%s' (type '%s'))",
				left.String(),
				left.Type(),

				right.String(),
				right.Type(),
			),
		)
		return STATUS_TYPE_ERROR
	}

	leftStr := left.(value.ValueString)
	rightStr := right.(value.ValueString)

	v.push(value.ValueString{ Value: leftStr.Value + rightStr.Value })
	return STATUS_OK
}

func (v *VM) getArguments(arity int) []value.Value {
	args := []value.Value{}
		
	for range arity {
		args = append(args, value.CopyValue(v.pop()))
	}

	return args
}

func (v *VM) call(callee value.Value, arity int) InterpretResult {
	if !isClosure(callee) && !isNativeFunction(callee) && !isRecord(callee) && !isBoundMethod(callee) {
		v.error(fmt.Sprintf("Can only call functions or records. (called '%s', of type '%s')", callee.String(), callee.Type()))
		return STATUS_TYPE_ERROR
	}

	switch function := callee.(type) {
		case value.ValueClosure: {
			if function.Fn.Arity != arity {
				v.error(fmt.Sprintf("Expected %d arguments, but got %d instead.", function.Fn.Arity, arity))
				return STATUS_INCORRECT_ARITY
			}
		
			v.callStack = append(v.callStack, CallFrame{
				function: &function,
				oldIp: v.ip,
			})
		
			args := v.getArguments(arity)
		
			// insert the arguments into the locals array
			for i := len(args) - 1; i >= 0; i-- {
				v.callStack[len(v.callStack) - 1].locals = append(v.callStack[len(v.callStack) - 1].locals, args[i])
			}
		
			v.ip = 0
			v.currentChunk = &function.Fn.Chunk
		
			v.pop() // The function.
		}

		case value.ValueNativeFn: {
			if function.Arity != arity {
				v.error(fmt.Sprintf("Expected %d arguments, but got %d instead.", function.Arity, arity))
				return STATUS_INCORRECT_ARITY
			}

			args := v.getArguments(arity)

			util.Reverse(args)
			v.pop() // The function.

			result := function.Fn(args)
			v.push(result)
		}

		case value.ValueRecord: {
			if len(function.FieldNames) != arity {
				v.error(fmt.Sprintf("Expected %d arguments, but got %d instead.", len(function.FieldNames), arity))
				return STATUS_INCORRECT_ARITY
			}

			args := v.getArguments(arity)

			util.Reverse(args)
			v.pop() // The record.

			// Create the object.
			instance := value.ValueInstance{
				Fields: args,
				Record: &function,
			}
			
			// Push it to the stack.
			v.push(instance)
		}

		case value.ValueBoundMethod: {
			if function.Method.Fn.Arity != arity {
				v.error(fmt.Sprintf("Expected %d arguments, but got %d instead.", function.Method.Fn.Arity, arity))
				return STATUS_INCORRECT_ARITY
			}
		
			v.callStack = append(v.callStack, CallFrame{
				function: &function.Method,
				oldIp: v.ip,
			})
		
			args := v.getArguments(arity)
			v.callStack[len(v.callStack) - 1].locals = append(v.callStack[len(v.callStack) - 1].locals, function.Receiver)

			// insert the arguments into the locals array
			for i := len(args) - 1; i >= 0; i-- {
				v.callStack[len(v.callStack) - 1].locals = append(v.callStack[len(v.callStack) - 1].locals, args[i])
			}
		
			v.ip = 0
			v.currentChunk = &function.Method.Fn.Chunk
		
			v.pop() // The function.
		}

		default:
			panic(fmt.Sprintf("Unknown called value: '%v'", callee))
	}

	return STATUS_OK
}

func (v *VM) getProperty(obj value.Value, index int) InterpretResult {
	nameValue := v.currentChunk.Constants[index]
	name := nameValue.(value.ValueString).Value
	
	switch instance := obj.(type) {
		case value.ValueInstance: {
			property, ok := instance.GetProperty(name)

			if !ok {
				v.error(fmt.Sprintf("Property '%s' doesn't exist.", name))
				return STATUS_PROPERTY_DOESNT_EXIST
			}

			v.push(property)
			return STATUS_OK
		}

		default: {
			// TODO: add methods to another types, defined by a table at runtime.
			v.error(fmt.Sprintf("This object ('%s') has no properties, because it isn't an instance. Its type is '%s'.", obj.String(), obj.Type()))
			return STATUS_PROPERTY_DOESNT_EXIST
		}
	}
}

func (v *VM) setProperty(obj value.Value, index int, val value.Value) InterpretResult {
	nameValue := v.currentChunk.Constants[index]
	name := nameValue.(value.ValueString).Value
	
	switch instance := obj.(type) {
		case value.ValueInstance: {
			ok := instance.SetProperty(name, val)

			if !ok {
				v.error(fmt.Sprintf("Property '%s' doesn't exist in the object '%s', of type '%s'.", name, obj.String(), obj.Type()))
				return STATUS_PROPERTY_DOESNT_EXIST
			}

			v.push(val)
			return STATUS_OK
		}

		default: {
			// TODO: add methods to another types, defined by a table at runtime.
			v.error(fmt.Sprintf("This object ('%s') has no properties, because it isn't an instance. Its type is '%s'.", obj.String(), obj.Type()))
			return STATUS_PROPERTY_DOESNT_EXIST
		}
	}
}

func (v *VM) binaryNum(operator byte) InterpretResult {
	right := v.pop()

	if v.stackIsEmpty() {
		v.error("Not enough stack items to perform a binary operation")
		return STATUS_STACK_EMPTY
	}

	left := v.pop()

	if !isNumber(left) || !isNumber(right) {
		v.error(
			fmt.Sprintf(
				"Operands must be numbers when performing arithmetic. (left: '%s', right: '%s')",
				left.String(),
				right.String(),
			),
		)
		return STATUS_TYPE_ERROR
	}

	leftNum := left.(value.ValueNumber)
	rightNum := right.(value.ValueNumber)

	switch operator {
		case compiler.OP_ADD: v.push(value.ValueNumber{ Value: leftNum.Value + rightNum.Value })
		case compiler.OP_SUB: v.push(value.ValueNumber{ Value: leftNum.Value - rightNum.Value })
		case compiler.OP_MUL: v.push(value.ValueNumber{ Value: leftNum.Value * rightNum.Value })
		case compiler.OP_DIV: {
			if rightNum.Value == 0 {
				v.error(
					fmt.Sprintf(
						"Cannot divide by zero. (left: '%s', right: '%s')",
						left.String(),
						right.String(),
					),
				)
				return STATUS_DIV_ZERO
			}

			v.push(value.ValueNumber{ Value: leftNum.Value / rightNum.Value })
		}
		case compiler.OP_MODULO: {
			if rightNum.Value == 0 {
				v.error(
					fmt.Sprintf(
						"Cannot divide by zero. (left: '%s', right: '%s')",
						left.String(),
						right.String(),
					),
				)
				return STATUS_DIV_ZERO
			}

			v.push(value.ValueNumber{ Value: math.Mod(leftNum.Value, rightNum.Value) })
		}
	}

	return STATUS_OK
}

func (v *VM) binaryComparison(operator byte) InterpretResult {
	right := v.pop()

	if v.stackIsEmpty() {
		v.error("Not enough stack items to perform a binary operation")
		return STATUS_STACK_EMPTY
	}

	left := v.pop()

	if !isNumber(left) || !isNumber(right) {
		v.error(
			fmt.Sprintf(
				"Operands must be numbers when comparing. (left: '%s' (type '%s'), right: '%s' (type '%s'))",
				left.String(),
				left.Type(),

				right.String(),
				right.Type(),
			),
		)
		return STATUS_TYPE_ERROR
	}

	leftNum := left.(value.ValueNumber)
	rightNum := right.(value.ValueNumber)

	switch operator {
		case compiler.OP_GREATER:
			v.push(value.ValueBool{ Value: leftNum.Value > rightNum.Value })
		case compiler.OP_GREATER_EQUAL:
			v.push(value.ValueBool{ Value: leftNum.Value >= rightNum.Value })
		case compiler.OP_LESS:
			v.push(value.ValueBool{ Value: leftNum.Value < rightNum.Value })
		case compiler.OP_LESS_EQUAL:
			v.push(value.ValueBool{ Value: leftNum.Value <= rightNum.Value })
	}

	return STATUS_OK
}

func (v *VM) binaryBool(operator byte) InterpretResult {
	right := v.pop()

	if v.stackIsEmpty() {
		v.error("Not enough stack items to perform a binary operation.")
		return STATUS_STACK_EMPTY
	}

	left := v.pop()

	if !isBool(left) || !isBool(right) {
		v.error(
			fmt.Sprintf(
				"Operands must be booleans to perform 'and' and 'or' operations. (left: '%s' (type '%s'), right: '%s' (type '%s'))",
				left.String(),
				left.Type(),

				right.String(),
				right.Type(),
			),
		)
		return STATUS_TYPE_ERROR
	}

	leftNum := left.(value.ValueBool)
	rightNum := right.(value.ValueBool)

	switch operator {
		case compiler.OP_AND: v.push(value.ValueBool{ Value: leftNum.Value && rightNum.Value })
		case compiler.OP_OR: v.push(value.ValueBool{ Value: leftNum.Value || rightNum.Value })
	}

	return STATUS_OK
}

// ---

func (v *VM) nextByte() byte {
	ip := v.ip
	v.ip += 1

	return v.currentChunk.Code[ip]
}

func (v *VM) stackIsEmpty() bool {
	return len(v.stack) == 0
}

func (v *VM) isAtEnd() bool {
	return v.ip >= len(v.currentChunk.Code)
}

// ---

func (v *VM) push(f value.Value) {
	v.stack = append(v.stack, f)
}

// can have errors
func (v *VM) pop() value.Value {
	if v.stackIsEmpty() {
		v.error("Performed a pop operation on an empty stack")
		return nil
	}

	return util.PopList(&v.stack)
}

func (v *VM) popFrame() CallFrame {
	if len(v.callStack) == 0 {
		v.error("Performed a pop operation on an empty call stack")
		return CallFrame{}
	}

	lastIndex := len(v.callStack) - 1

	topElement := v.callStack[lastIndex]
	v.callStack = v.callStack[:lastIndex]

	return topElement
}

func (v *VM) peek(offset int) value.Value {
	pos := len(v.stack) - 1 - offset
	if pos < 0 || pos > len(v.stack) - 1 {
		v.error("Peek position out of bounds")
		return nil
	}

	return v.stack[pos]
}

func (v *VM) popVar() value.Value {
	return v.popnVar(1)
}

func (v *VM) popnVar(amount int) value.Value {
	lastIndex := len(v.callStack[len(v.callStack)-1].locals) - amount
	topElement := v.callStack[len(v.callStack)-1].locals[lastIndex]

	v.callStack[len(v.callStack)-1].locals = v.callStack[len(v.callStack)-1].locals[:lastIndex]
	return topElement
}

func (v *VM) error(message string) {
	metadata := v.currentChunk.Metadata[v.oldIp]

	fmt.Printf("[-] Runtime error at %s (%d, %d): %s\n", v.fileData.Name, metadata.Position.Line + 1, metadata.Position.Col + 1, message)
	fmt.Printf(" |  %d | %s\n", metadata.Position.Line + 1, v.fileData.Lines[metadata.Position.Line])
	fmt.Printf(" | %s    %s%s\n", strings.Repeat(" ", len(strconv.Itoa(metadata.Position.Line + 1))), strings.Repeat(" ", metadata.Position.Col), strings.Repeat("^", metadata.Length))
	fmt.Println("[-]")

	if len(v.callStack) > 0 {
		for i := len(v.callStack) - 1; i >= 0; i-- {
			posChunk := v.topLevel
			frame := v.callStack[i]

			if i > 0 {
				posChunk = v.callStack[i - 1].function.Fn.Chunk
			}

			pos := posChunk.Metadata[frame.oldIp - 1].Position
			name := frame.function.Fn.Name

			if name == nil {
				fmt.Printf(" | in <anonymous> (%d, %d)\n", pos.Line + 1, pos.Col + 1)
			} else {
				fmt.Printf(" | in %s (%d, %d)\n", *name, pos.Line + 1, pos.Col + 1)
			}
		}

		fmt.Print("[-]\n\n")
	} else {
		fmt.Println()
	}
	
	v.hadError = true
}
