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
		inst := d.nextByte()
		d.PrintInstruction(inst)
	}
}

func (d *Disassembler) PrintInstruction(inst byte) {
	ip := d.ip

	switch inst {
		// inst index value
		case compiler.OP_PUSH_CONST: {
			index, _ := util.BytesToInt([]byte(d.code[d.ip : d.ip+4]))
			d.ip += 4

			fmt.Printf(
				"%s %s | %s (%s)\n",
				util.PadLeft(strconv.Itoa(ip), 4, " "),
				util.PadRight(getInstructionName(inst), MAX_INSTRUCTION_LENGTH, " "),

				util.PadRight(strconv.Itoa(index), 4, " "),
				d.constants[index].String(),
			)
		}

		// inst [int]
		case compiler.OP_POPN_VAR, compiler.OP_GET_VAR, compiler.OP_SET_VAR: {
			count, _ := util.BytesToInt([]byte(d.code[d.ip : d.ip+4]))
			d.ip += 4

			fmt.Printf(
				"%s %s | %d\n",
				util.PadLeft(strconv.Itoa(ip), 4, " "),
				util.PadRight(getInstructionName(inst), MAX_INSTRUCTION_LENGTH, " "),

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
				util.PadRight(getInstructionName(inst), MAX_INSTRUCTION_LENGTH, " "),

				util.PadRight(strconv.Itoa(count), 4, " "),
				d.ip + count,
			)
		}

		// inst
		default: {
			fmt.Printf(
				"%s %s |\n",
				util.PadLeft(strconv.Itoa(ip), 4, " "),
				util.PadRight(getInstructionName(inst), MAX_INSTRUCTION_LENGTH, " "),
			)
		}
	}
}
