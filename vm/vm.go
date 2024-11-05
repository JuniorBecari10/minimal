package vm

import (
	"fmt"
	"vm-go/compiler"
	"vm-go/util"
)

type InterpretResult int

const (
	STATUS_OK = iota
	STATUS_STACK_EMPTY
	STATUS_DIV_ZERO
)

type VM struct {
	code      string
	constants []util.Value

	stack     []util.Value
	variables []util.Value

	ip int
	hadError bool
}

func NewVM(code string, constants []util.Value) *VM {
	return &VM{
		code:      code,
		constants: constants,

		stack:     []util.Value{},
		variables: []util.Value{},

		ip:        0,
		hadError:  false,
	}
}

func (v *VM) Run() InterpretResult {
	for !v.isAtEnd() {
		i := v.nextByte()

		switch i {
			case compiler.OP_CONSTANT: {
				index, _ := util.BytesToInt([]byte(v.code[v.ip:v.ip + 4]))
				v.ip += 4

				v.Push(v.constants[index])
			}

			case compiler.OP_ADD, compiler.OP_SUB, compiler.OP_MUL, compiler.OP_DIV: {
				status := v.binary(i)

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

			case compiler.OP_POP_VAR: {
				v.PopVar()
			}

			case compiler.OP_POPN_VAR: {
				amount, _ := util.BytesToInt([]byte(v.code[v.ip:v.ip + 4]))
				v.ip += 4

				v.PopnVar(amount)
			}

			case compiler.OP_PRINT: fmt.Printf("%.2f\n", v.Pop())
		}
	}

	return STATUS_OK
}

// ---

func (v *VM) binary(operator byte) InterpretResult {
	right := v.Pop()

	if v.stackIsEmpty() {
		v.error("Not enough stack items to perform a binary operation")
		return STATUS_STACK_EMPTY
	}

	left := v.Pop()

	switch operator {
		case compiler.OP_ADD: v.Push(left + right)
		case compiler.OP_SUB: v.Push(left - right)
		case compiler.OP_MUL: v.Push(left * right)
		case compiler.OP_DIV: {
			if right == 0 {
				v.error("Cannot divide by zero")
				return STATUS_DIV_ZERO
			}

			v.Push(left / right)
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

func (v *VM) Push(f util.Value) {
	v.stack = append(v.stack, f)
}

// can have errors
func (v *VM) Pop() util.Value {
	if v.stackIsEmpty() {
		v.error("Performed a pop operation on an empty stack")
		return 0
	}

	lastIndex := len(v.stack) - 1
	topElement := v.stack[lastIndex]

	v.stack = v.stack[:lastIndex]
	return topElement
}

func (v *VM) PopVar() util.Value {
	return v.PopnVar(1)
}

func (v *VM) PopnVar(n int) util.Value {
	lastIndex := len(v.variables) - n
	topElement := v.variables[lastIndex]

	v.variables = v.variables[:lastIndex]
	return topElement
}

func (v *VM) error(message string) {
	fmt.Printf("Error at VM: %s\n", message)
	v.hadError = true
}
