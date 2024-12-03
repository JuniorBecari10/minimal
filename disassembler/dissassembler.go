package disassembler

import (
	"fmt"
	"strconv"
	"vm-go/chunk"
	"vm-go/compiler"
	"vm-go/util"
	"vm-go/value"
)

type Disassembler struct {
	chunk chunk.Chunk
	ip int

	fileData *util.FileData
}

func NewDisassembler(chunk chunk.Chunk, fileData *util.FileData) *Disassembler {
	return &Disassembler{
		chunk: chunk,
		ip: 0,

		fileData: fileData,
	}
}

const MAX_INSTRUCTION_LENGTH = 16

func (d *Disassembler) Disassemble() {
	d.disassemble("top-level")
}

func (d *Disassembler) disassemble(name string) {
	fmt.Printf("chunk: %s\n", name)
	fmt.Println(" offset | position  | instruction      | index  | result")
	fmt.Println("--------|-----------|------------------|--------|--------")

	i := 0
	for !d.isAtEnd() {
		ip := d.ip
		inst := d.nextByte()

		d.PrintInstruction(inst, ip, i)
		i++
	}

	fmt.Println()

	for _, c := range d.chunk.Constants {
		switch fn := c.(type) {
			case value.ValueFunction: {
				fnDiss := NewDisassembler(fn.Chunk, d.fileData)
				fnDiss.disassemble(fmt.Sprintf("function in %s's constant table", name))
			}
		}
	}
}

func (d *Disassembler) PrintInstruction(inst byte, ip int, i int) {
	fmt.Printf(
		" %s | %s %s | %s | ",
		util.PadLeft(strconv.Itoa(ip), 6, " "),

		util.PadRight(strconv.Itoa(d.chunk.Positions[i].Line + 1), 4, " "),
		util.PadRight(strconv.Itoa(d.chunk.Positions[i].Col + 1), 4, " "),

		util.PadRight(getInstructionName(inst), MAX_INSTRUCTION_LENGTH, " "),
	)

	switch inst {
		// inst index value
		case compiler.OP_PUSH_CONST: {
			index, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			// TODO: print the type as well
			fmt.Printf(
				"%s | %s\n",
				util.PadRight(strconv.Itoa(index), 6, " "),
				d.chunk.Constants[index].String(),
			)
		}

		// inst [int]
		case compiler.OP_POPN_VAR, compiler.OP_GET_VAR, compiler.OP_SET_VAR: {
			count, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			fmt.Printf(
				"%s |\n",
				util.PadRight(strconv.Itoa(count), 6, " "),
			)
		}

		// inst amount result
		case compiler.OP_JUMP_FALSE, compiler.OP_JUMP: {
			count, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			fmt.Printf(
				"%s | %d\n",
				util.PadRight(strconv.Itoa(count), 6, " "),
				d.ip + count,
			)
		}

		// inst amount result
		case compiler.OP_LOOP: {
			count, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			fmt.Printf(
				"%s | %d\n",
				util.PadRight(strconv.Itoa(count), 6, " "),
				d.ip - count,
			)
		}

		// inst
		default:
			// add the separator between index and constant columns
			fmt.Println("       |")
	}
}
