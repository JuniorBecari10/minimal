package disassembler

import (
	"fmt"
	"strconv"
	"vm-go/compiler"
	"vm-go/util"
	"vm-go/value"
)

type Disassembler struct {
	code      []byte
	constants []value.Value
	ip int
}

func NewDisassembler(code []byte, constants []value.Value) *Disassembler {
	return &Disassembler{
		code:      code,
		constants: constants,
		ip: 0,
	}
}

const MAX_INSTRUCTION_LENGTH = 13

func (d *Disassembler) Disassemble() {
	for !d.isAtEnd() {
		ip := d.ip
		i := d.nextByte()

		switch i {
			// inst index value
			case compiler.OP_PUSH_CONST: {
				index, _ := util.BytesToInt([]byte(d.code[d.ip : d.ip+4]))
				d.ip += 4

				fmt.Printf(
					"%s %s | %s (%.2f)\n",
					util.PadLeft(strconv.Itoa(ip), 4, " "),
					util.PadRight(getInstructionName(i), MAX_INSTRUCTION_LENGTH, " "),

					util.PadRight(strconv.Itoa(index), 4, " "),
					d.constants[index],
				)
			}

			// inst [int]
			case compiler.OP_POPN_VAR, compiler.OP_GET_VAR, compiler.OP_SET_VAR: {
				count, _ := util.BytesToInt([]byte(d.code[d.ip : d.ip+4]))
				d.ip += 4

				fmt.Printf(
					"%s %s | %d\n",
					util.PadLeft(strconv.Itoa(ip), 4, " "),
					util.PadRight(getInstructionName(i), MAX_INSTRUCTION_LENGTH, " "),

					count,
				)
			}

			// inst amount result
			case compiler.OP_JUMP_FALSE, compiler.OP_JUMP, compiler.OP_LOOP: {
				count, _ := util.BytesToInt([]byte(d.code[d.ip : d.ip+4]))
				d.ip += 4

				fmt.Printf(
					"%s %s | %s (%d)\n",
					util.PadLeft(strconv.Itoa(ip), 4, " "),
					util.PadRight(getInstructionName(i), MAX_INSTRUCTION_LENGTH, " "),

					util.PadRight(strconv.Itoa(count), 4, " "),
					d.ip + count,
				)
			}

			// inst
			default: {
				fmt.Printf(
					"%s %s |\n",
					util.PadLeft(strconv.Itoa(ip), 4, " "),
					util.PadRight(getInstructionName(i), MAX_INSTRUCTION_LENGTH, " "),
				)
			}
		}
	}
}

// ---

func (d *Disassembler) nextByte() byte {
	ip := d.ip
	d.ip += 1

	return d.code[ip]
}

func (d *Disassembler) isAtEnd() bool {
	return d.ip >= len(d.code)
}

// ---

func getInstructionName(inst byte) string {
	switch inst {
		case compiler.OP_PUSH_CONST: return "OP_PUSH_CONST"

		case compiler.OP_ADD: return "OP_ADD"
		case compiler.OP_SUB: return "OP_SUB"
		case compiler.OP_MUL: return "OP_MUL"
		case compiler.OP_DIV: return "OP_DIV"

		case compiler.OP_DEF_VAR: return "OP_DEF_VAR"
		case compiler.OP_GET_VAR: return "OP_GET_VAR"
		case compiler.OP_SET_VAR: return "OP_SET_VAR"

		case compiler.OP_POP: return "OP_POP"
		case compiler.OP_POP_VAR: return "OP_POP_VAR"
		case compiler.OP_POPN_VAR: return "OP_POPN_VAR"

		case compiler.OP_JUMP: return "OP_JUMP"
		case compiler.OP_JUMP_FALSE: return "OP_JUMP_FALSE"
		case compiler.OP_LOOP: return "OP_LOOP_FALSE"

		case compiler.OP_EQUAL: return "OP_EQUAL"
		case compiler.OP_NOT_EQUAL: return "OP_NOT_EQUAL"

		case compiler.OP_GREATER: return "OP_GREATER"
		case compiler.OP_GREATER_EQUAL: return "OP_GREATER_EQUAL"

		case compiler.OP_LESS: return "OP_LESS"
		case compiler.OP_LESS_EQUAL: return "OP_LESS_EQUAL"

		case compiler.OP_PRINT: return "OP_PRINT"

		default:
			return "Unknown"
	}
}
