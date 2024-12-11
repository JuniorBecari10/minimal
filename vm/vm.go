package vm

import (
	"fmt"
	"math"
	"strings"
	"vm-go/chunk"
	"vm-go/compiler"
	"vm-go/util"
	"vm-go/value"
)

type InterpretResult int

const (
	STATUS_OK = iota
	STATUS_STACK_EMPTY
	STATUS_OUT_OF_BOUNDS
	STATUS_DIV_ZERO
	STATUS_TYPE_ERROR
	STATUS_INCORRECT_ARITY
)

type VM struct {
	currentChunk *chunk.Chunk
	topLevel chunk.Chunk

	stack     []value.Value
	globals   []value.Value
	callStack []CallFrame

	ip int
	oldIp int

	hadError bool
	fileData *util.FileData
}

func NewVM(chunk chunk.Chunk, fileData *util.FileData) *VM {
	vm := VM{
		currentChunk: &chunk,
		topLevel: chunk,

		stack:     []value.Value{},
		globals:   []value.Value{},
		callStack: []CallFrame{},

		ip:        0,
		oldIp:     0,

		hadError:  false,
		fileData: fileData,
	}

	vm.includeNativeFns()
	return &vm
}

func (v *VM) Run() InterpretResult {
	for !v.isAtEnd() && !v.hadError {
		v.oldIp = v.ip
 		i := v.nextByte()

		switch i {
			case compiler.OP_PUSH_CONST:
				v.push(v.currentChunk.Constants[v.getInt()])

			// TODO: add a separated opcode for concatenating strings when typechecking is added
			case compiler.OP_ADD: {
				if !typesEqual(v.peek(0), v.peek(1)) {
					v.error(
						fmt.Sprintf(
							"Operands types must be equal when adding/concatenating. (left: '%s', right: '%s')",
							v.peek(1).String(),
							v.peek(0).String(),
						),
					)
					return STATUS_TYPE_ERROR
				}

				// if one operand is string, the other should be too
				if isString(v.peek(0)) {
					status := v.concatenateStrs()

					if status != STATUS_OK {
						return status
					}
				} else if isNumber(v.peek(0)) {
					status := v.binaryNum(i)

					if status != STATUS_OK {
						return status
					}
				} else {
					v.error(
						fmt.Sprintf(
							"Operands must be numbers or strings when adding/concatenating. (left: '%s', right: '%s')",
							v.peek(1).String(),
							v.peek(0).String(),
						),
					)
					return STATUS_TYPE_ERROR
				}
			}

			case compiler.OP_SUB, compiler.OP_MUL, compiler.OP_DIV, compiler.OP_MODULO: {
				status := v.binaryNum(i)

				if status != STATUS_OK {
					return status
				}
			}

			case compiler.OP_DEF_LOCAL:
				v.callStack[len(v.callStack)-1].locals = append(v.callStack[len(v.callStack)-1].locals, v.pop())

			case compiler.OP_GET_LOCAL:
				v.push(v.callStack[len(v.callStack)-1].locals[v.getInt()])

			case compiler.OP_SET_LOCAL:
				v.callStack[len(v.callStack)-1].locals[v.getInt()] = v.peek(0)

			case compiler.OP_DEF_GLOBAL:
				v.globals = append(v.globals, v.pop())

			case compiler.OP_GET_GLOBAL:
				v.push(v.globals[v.getInt()])

			case compiler.OP_SET_GLOBAL:
				v.globals[v.getInt()] = v.peek(0)

			case compiler.OP_POP:
				v.pop()

			case compiler.OP_POP_VAR:
				v.popVar()

			case compiler.OP_POPN_VAR:
				v.popnVar(v.getInt())

			case compiler.OP_JUMP:
				v.ip += v.getInt()

			case compiler.OP_JUMP_TRUE: {
				amount := v.getInt()

				// TODO: check for out of bounds by checking nil
				if b, ok := v.peek(0).(value.ValueBool); ok {
					if v.hadError {
						return STATUS_OUT_OF_BOUNDS
					}

					if b.Value {
						v.ip += amount
					}
				} else {
					v.error(fmt.Sprintf("Expression is not a boolean. (value: '%s')", v.peek(0).String()))
					return STATUS_TYPE_ERROR
				}
			}

			case compiler.OP_JUMP_FALSE: {
				amount := v.getInt()

				// TODO: check for out of bounds by checking nil
				if b, ok := v.peek(0).(value.ValueBool); ok {
					if v.hadError {
						return STATUS_OUT_OF_BOUNDS
					}

					if !b.Value {
						v.ip += amount
					}
				} else {
					v.error(fmt.Sprintf("Expression is not a boolean. (value: '%s')", v.peek(0).String()))
					return STATUS_TYPE_ERROR
				}
			}

			case compiler.OP_LOOP:
				v.ip -= v.getInt()

			case compiler.OP_EQUAL: {
				b := v.pop()
				a := v.pop()

				if !typesEqual(a, b) {
					v.error(
						fmt.Sprintf(
							"Types must be the same when comparing. (left: '%s', right: '%s')",
							a.String(),
							b.String(),
						),
					)
					return STATUS_TYPE_ERROR
				}

				v.push(value.ValueBool{ Value: valuesEqual(a, b) })
			}

			case compiler.OP_NOT_EQUAL: {
				b := v.pop()
				a := v.pop()

				if !typesEqual(a, b) {
					v.error(
						fmt.Sprintf(
							"Types must be the same when comparing. (left: '%s', right: '%s')",
							a.String(),
							b.String(),
						),
					)
					return STATUS_TYPE_ERROR
				}

				v.push(value.ValueBool{ Value: !valuesEqual(a, b) })
			}

			case compiler.OP_GREATER, compiler.OP_GREATER_EQUAL, compiler.OP_LESS, compiler.OP_LESS_EQUAL: {
				status := v.binaryComparison(i)

				if status != STATUS_OK {
					return status
				}
			}

			case compiler.OP_AND, compiler.OP_OR: {
				status := v.binaryBool(i)

				if status != STATUS_OK {
					return status
				}
			}

			case compiler.OP_NOT: {
				op := v.pop()

				if !isBool(op) {
					v.error(
						fmt.Sprintf(
							"Operand must be a boolean for performing a logical not. (value: '%s')",
							op.String(),
						),
					)
					return STATUS_TYPE_ERROR
				}

				opBool := op.(value.ValueBool)
				v.push(value.ValueBool{ Value: !opBool.Value })
			}

			case compiler.OP_NEGATE: {
				op := v.pop()

				if !isNumber(op) {
					v.error(
						fmt.Sprintf(
							"Operand must be a boolean for performing a number negation. (value: '%s')",
							op.String(),
						),
					)
					return STATUS_TYPE_ERROR
				}

				opNum := op.(value.ValueNumber)
				v.push(value.ValueNumber{ Value: -opNum.Value })
			}

			case compiler.OP_RETURN: {
				frame := v.popFrame()
				var chunk chunk.Chunk

				if len(v.callStack) == 0 {
					chunk = v.topLevel
				} else {
					chunk = v.callStack[len(v.callStack) - 1].function.Chunk
				}

				v.ip = frame.oldIp
				v.currentChunk = &chunk
			}

			case compiler.OP_TRUE: v.push(value.ValueBool{ Value: true })
			case compiler.OP_FALSE: v.push(value.ValueBool{ Value: false })

			case compiler.OP_NIL: v.push(value.ValueNil{})
			case compiler.OP_VOID: v.push(value.ValueVoid{})

			case compiler.OP_CALL: {
				arity := v.getInt()
				status := v.call(v.peek(arity), arity)

				if status != STATUS_OK {
					return status
				}
			}

			case compiler.OP_PRINT:
				fmt.Println(v.pop().String())

			default:
				panic(fmt.Sprintf("Unknown instruction: '%d'", i))
		}
	}

	return STATUS_OK
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
				"Operands must be strings when performing concatenation. (left: '%s', right: '%s')",
				left.String(),
				right.String(),
			),
		)
		return STATUS_TYPE_ERROR
	}

	leftStr := left.(value.ValueString)
	rightStr := right.(value.ValueString)

	v.push(value.ValueString{ Value: leftStr.Value + rightStr.Value })
	return STATUS_OK
}

func (v *VM) call(callee value.Value, arity int) InterpretResult {
	if !isFunction(callee) && !isNativeFunction(callee) {
		v.error(fmt.Sprintf("Can only call functions. (Called '%s')", callee.String()))
		return STATUS_TYPE_ERROR
	}

	switch function := callee.(type) {
		case value.ValueFunction: {
			if function.Arity != arity {
				v.error(fmt.Sprintf("Expected %d arguments, but got %d instead.", function.Arity, arity))
				return STATUS_INCORRECT_ARITY
			}
		
			v.callStack = append(v.callStack, CallFrame{
				function: &function,
				oldIp: v.ip,
			})
		
			vars := []value.Value{}
		
			for range arity {
				vars = append(vars, v.pop())
			}
		
			for i := len(vars) - 1; i >= 0; i-- {
				v.callStack[len(v.callStack) - 1].locals = append(v.callStack[len(v.callStack) - 1].locals, vars[i])
			}
		
			v.ip = 0
			v.currentChunk = &function.Chunk
		
			v.pop() // the function
		}

		case value.ValueNativeFn: {
			if function.Arity != arity {
				v.error(fmt.Sprintf("Expected %d arguments, but got %d instead.", function.Arity, arity))
				return STATUS_INCORRECT_ARITY
			}

			vars := []value.Value{}
		
			for range arity {
				vars = append(vars, v.pop())
			}

			util.Reverse(vars)
			result := function.Fn(vars)

			v.push(result)
		}
	}

	return STATUS_OK
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
				"Operands must be numbers when performing comparison. (left: '%s', right: '%s')",
				left.String(),
				right.String(),
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
		v.error("Not enough stack items to perform a binary operation")
		return STATUS_STACK_EMPTY
	}

	left := v.pop()

	if !isBool(left) || !isBool(right) {
		v.error(
			fmt.Sprintf(
				"Operands must be booleans when performing boolean operations. (left: '%s', right: '%s')",
				left.String(),
				right.String(),
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

	lastIndex := len(v.stack) - 1
	
	topElement := v.stack[lastIndex]
	v.stack = v.stack[:lastIndex]

	return topElement
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
	pos := v.currentChunk.Positions[v.oldIp] // TODO: save the last ip to use here

	fmt.Printf("[-] Runtime error at %s (%d, %d): %s\n", v.fileData.Name, pos.Line + 1, pos.Col + 1, message)
	fmt.Printf(" | %s\n", v.fileData.Lines[pos.Line])
	fmt.Printf(" | %s^\n", strings.Repeat(" ", pos.Col))
	fmt.Println("[-]")

	if len(v.callStack) > 0 {
		for i := len(v.callStack) - 1; i >= 0; i-- {
			posChunk := v.topLevel
			frame := v.callStack[i]

			if i > 0 {
				posChunk = v.callStack[i - 1].function.Chunk
			}

			pos := posChunk.Positions[frame.oldIp - 1]
			name := frame.function.Name

			if name == nil {
				fmt.Printf(" | in (%d, %d)\n", pos.Line + 1, pos.Col + 1)
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
