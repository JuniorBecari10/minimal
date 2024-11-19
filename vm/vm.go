package vm

import (
	"fmt"
	"vm-go/compiler"
	"vm-go/util"
	"vm-go/value"
)

type InterpretResult = int

const (
	STATUS_OK = iota
	STATUS_STACK_EMPTY
	STATUS_OUT_OF_BOUNDS
	STATUS_DIV_ZERO
	STATUS_TYPE_ERROR
)

type VM struct {
	code      []byte
	constants []value.Value

	stack     []value.Value
	variables []value.Value

	ip int
	hadError bool
}

func NewVM(code []byte, constants []value.Value) *VM {
	return &VM{
		code:      code,
		constants: constants,

		stack:     []value.Value{},
		variables: []value.Value{},

		ip:        0,
		hadError:  false,
	}
}

func (v *VM) Run() InterpretResult {
	for !v.isAtEnd() {
		i := v.nextByte()

		switch i {
			case compiler.OP_PUSH_CONST: {
				index, _ := util.BytesToInt([]byte(v.code[v.ip:v.ip + 4]))
				v.ip += 4

				v.Push(v.constants[index])
			}

			case compiler.OP_ADD, compiler.OP_SUB, compiler.OP_MUL, compiler.OP_DIV: {
				status := v.binaryNum(i)

				if status != STATUS_OK {
					return status
				}
			}

			case compiler.OP_DEF_VAR: {
				v.variables = append(v.variables, v.Pop())
			}

			case compiler.OP_GET_VAR: {
				index, _ := util.BytesToInt([]byte(v.code[v.ip:v.ip + 4]))
				v.ip += 4

				v.Push(v.variables[index])
			}

			case compiler.OP_SET_VAR: {
				index, _ := util.BytesToInt([]byte(v.code[v.ip:v.ip + 4]))
				v.ip += 4

				v.variables[index] = v.Pop()
			}

			case compiler.OP_POP: {
				v.Pop()
			}

			case compiler.OP_POP_VAR: {
				v.PopVar()
			}

			case compiler.OP_POPN_VAR: {
				amount, _ := util.BytesToInt([]byte(v.code[v.ip:v.ip + 4]))
				v.ip += 4

				v.PopnVar(amount)
			}

			case compiler.OP_JUMP: {
				amount, _ := util.BytesToInt([]byte(v.code[v.ip:v.ip + 4]))
				v.ip += amount
			}

			case compiler.OP_JUMP_FALSE: {
				amount, _ := util.BytesToInt([]byte(v.code[v.ip:v.ip + 4]))
				v.ip += 4
				
				// TODO: check for out of bounds by checking nil
				if b, ok := v.Peek(0).(value.ValueBool); ok {
					if v.hadError {
						return STATUS_OUT_OF_BOUNDS
					}

					if !b.Value {
						v.ip += amount
					}
				} else {
					v.error("Expression is not a boolean")
					return STATUS_TYPE_ERROR
				}
			}

			case compiler.OP_LOOP_FALSE: {
				amount, _ := util.BytesToInt([]byte(v.code[v.ip:v.ip + 4]))
				
				// TODO: check for out of bounds by checking nil
				if b, ok := v.Peek(0).(value.ValueBool); ok {
					if v.hadError {
						return STATUS_OUT_OF_BOUNDS
					}

					if !b.Value {
						v.ip -= amount
					}
				} else {
					v.error("Expression is not a boolean")
					return STATUS_TYPE_ERROR
				}
			}

			case compiler.OP_EQUAL: {
				b := v.Pop()
				a := v.Pop()

				if !checkTypes(a, b) {
					v.error("Types must be the same when comparing")
					return STATUS_TYPE_ERROR
				}

				switch val := a.(type) {
					case value.ValueNumber:
						v.Push(value.ValueBool{ Value: val.Value == b.(value.ValueNumber).Value })
					
					case value.ValueBool:
						v.Push(value.ValueBool{ Value: val.Value == b.(value.ValueBool).Value })
				}
			}

			case compiler.OP_PRINT: fmt.Printf("%.2f\n", v.Pop())

			default:
				panic(fmt.Sprintf("Unknown instruction: '%d'", i))
		}
	}

	return STATUS_OK
}

// ---

func (v *VM) binaryNum(operator byte) InterpretResult {
	right := v.Pop()

	if v.stackIsEmpty() {
		v.error("Not enough stack items to perform a binary operation")
		return STATUS_STACK_EMPTY
	}

	left := v.Pop()

	if !isNumber(left) || !isNumber(right) {
		v.error("Types must be numbers when performing arithmetic")
		return STATUS_TYPE_ERROR
	}

	leftNum := left.(value.ValueNumber)
	rightNum := right.(value.ValueNumber)

	switch operator {
		case compiler.OP_ADD: v.Push(value.ValueNumber{ Value: leftNum.Value + rightNum.Value })
		case compiler.OP_SUB: v.Push(value.ValueNumber{ Value: leftNum.Value - rightNum.Value })
		case compiler.OP_MUL: v.Push(value.ValueNumber{ Value: leftNum.Value * rightNum.Value })
		case compiler.OP_DIV: {
			if rightNum.Value == 0 {
				v.error("Cannot divide by zero")
				return STATUS_DIV_ZERO
			}

			v.Push(value.ValueNumber{ Value: leftNum.Value / rightNum.Value })
		}
	}

	return STATUS_OK
}

// ---

func (v *VM) nextByte() byte {
	ip := v.ip
	v.ip += 1

	return v.code[ip]
}

func (v *VM) stackIsEmpty() bool {
	return len(v.stack) == 0
}

func (v *VM) isAtEnd() bool {
	return v.ip >= len(v.code)
}

// ---

func (v *VM) Push(f value.Value) {
	v.stack = append(v.stack, f)
}

// can have errors
func (v *VM) Pop() value.Value {
	if v.stackIsEmpty() {
		v.error("Performed a pop operation on an empty stack")
		return nil
	}

	lastIndex := len(v.stack) - 1
	topElement := v.stack[lastIndex]

	v.stack = v.stack[:lastIndex]
	return topElement
}

func (v *VM) Peek(offset int) value.Value {
	pos := len(v.stack) - 1 - offset
	if pos < 0 || pos > len(v.stack) - 1 {
		v.error("Peek position out of bounds")
		return nil
	}

	return v.stack[pos]
}

func (v *VM) PopVar() value.Value {
	return v.PopnVar(1)
}

func (v *VM) PopnVar(n int) value.Value {
	lastIndex := len(v.variables) - n
	topElement := v.variables[lastIndex]

	v.variables = v.variables[:lastIndex]
	return topElement
}

func (v *VM) error(message string) {
	fmt.Printf("Error at VM: %s\n", message)
	v.hadError = true
}
