package disassembler

import (
	"fmt"
	"minlib/instructions"
	"minlib/value"
	"minlib/util"
	"strconv"
	"strings"
)

type Disassembler struct {
	chunk value.Chunk
	ip int
}

func NewDisassembler(chunk value.Chunk) *Disassembler {
	return &Disassembler{
		chunk: chunk,
		ip: 0,
	}
}

const MAX_INSTRUCTION_LENGTH = 16

func (d *Disassembler) Disassemble() {
	d.disassemble("top-level")
}

func (d *Disassembler) disassemble(name string) {
	fmt.Println(util.Center(fmt.Sprintf("- %s -", name), 66, " ")) // 66 = len("--------|-----------|--------|------------------|--------|--------")

	if len(d.chunk.Code) == 0 {
		fmt.Println(util.Center("function is empty.", 66, " "))
		fmt.Println()
		return
	}

	fmt.Println("--------|-----------|--------|------------------|--------|--------")
	fmt.Println(" offset | position  | length | instruction      | index  | result")
	fmt.Println("--------|-----------|--------|------------------|--------|--------")

	for i := 0; !d.isAtEnd(); i++ {
		ip := d.ip
		inst := d.nextByte()

		d.PrintInstruction(inst, ip, i)
	}

	fmt.Println()

	for _, c := range d.chunk.Constants {
		switch fn := c.(type) {
			case value.ValueFunction: {
				fnDiss := NewDisassembler(fn.Chunk)
				fnName := fmt.Sprintf("anonymous function, in %s", name)

				if fn.Name != nil {
					fnName = fmt.Sprintf("function '%s', in %s", *fn.Name, name)
				}

				fnDiss.disassemble(fnName)
			}
		}
	}
}
func (d *Disassembler) PrintInstruction(inst byte, ip int, i int) {
	fmt.Printf(
		" %s | %s %s | %s | %s | ",
		util.PadLeft(strconv.Itoa(ip), 6, " "),

		util.PadRight(strconv.Itoa(int(d.chunk.Metadata[ip].Position.Line + 1)), 4, " "),
		util.PadRight(strconv.Itoa(int(d.chunk.Metadata[ip].Position.Col + 1)), 4, " "),

		util.PadRight(strconv.Itoa(int(d.chunk.Metadata[ip].Length)), 6, " "),

		util.PadRight(getInstructionName(inst), MAX_INSTRUCTION_LENGTH, " "),
	)

	switch inst {
		// inst index value
		case instructions.PUSH_CONST: {
			index, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			str := d.chunk.Constants[index].String()

			switch d.chunk.Constants[index].(type) {
				case value.ValueString:
					str = strconv.Quote(str)
			}

			// TODO: print the type as well
			fmt.Printf(
				"%s | %s: %s\n",
				util.PadRight(strconv.Itoa(index), 6, " "),
				str,
				d.chunk.Constants[index].Type(),
			)
		}

        case instructions.GET_PROPERTY, instructions.SET_PROPERTY: {
			index, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			str := d.chunk.Constants[index].String()

			fmt.Printf(
				"%s | %s\n",
				util.PadRight(strconv.Itoa(index), 6, " "),
				str,
			)
        }

		// inst index value count + metadata
		case instructions.PUSH_CLOSURE: {
			index, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			upvalueCount, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			// TODO: print the type as well
			fmt.Printf(
				"%s | %s\n",
				util.PadRight(strconv.Itoa(index), 6, " "),
				d.chunk.Constants[index].String(),
			)

			for i := range upvalueCount {
				isLocal := d.chunk.Code[d.ip] == 1
				d.ip += 1

				upvalueIndex, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
				d.ip += 4

				var text string
				
				if isLocal {
					text = "local"
				} else {
					text = "upvalue"
				}

				fmt.Printf(
					" %s | %s %s | %s | |%s | %s | %s\n",
					util.PadLeft(strconv.Itoa(ip + 5 * (i + 1)), 6, " "),
			
					util.PadRight(strconv.Itoa(int(d.chunk.Metadata[ip].Position.Line + 1)), 4, " "),
					util.PadRight(strconv.Itoa(int(d.chunk.Metadata[ip].Position.Col + 1)), 4, " "),

					util.PadRight(strconv.Itoa(int(d.chunk.Metadata[ip].Length)), 6, " "),
			
					strings.Repeat(" ", MAX_INSTRUCTION_LENGTH - 1),

					util.PadRight(strconv.Itoa(upvalueIndex), 6, " "),
					text,
				)
			}
		}

		// inst [int]
		case instructions.POPN_LOCAL,
			instructions.GET_LOCAL, instructions.SET_LOCAL,
			instructions.GET_UPVALUE, instructions.SET_UPVALUE,
			instructions.GET_GLOBAL, instructions.SET_GLOBAL,
			instructions.CALL, instructions.APPEND_METHODS,
			instructions.EXIT: {
			count, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			fmt.Printf(
				"%s |\n",
				util.PadRight(strconv.Itoa(count), 6, " "),
			)
		}

		// inst [int] [int]
		case instructions.CALL_PROPERTY: {
			index, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			arity, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			fmt.Printf(
				"%s | %s\n",
				util.PadRight(strconv.Itoa(index), 6, " "),
				util.PadRight(strconv.Itoa(arity), 6, " "),
			)
		}

		// inst amount result (add)
		case instructions.JUMP, instructions.JUMP_TRUE, instructions.JUMP_FALSE, instructions.JUMP_HAS_NO_NEXT: {
			count, _ := util.BytesToInt(d.chunk.Code[d.ip : d.ip+4])
			d.ip += 4

			fmt.Printf(
				"%s | %d\n",
				util.PadRight(strconv.Itoa(count), 6, " "),
				d.ip + count,
			)
		}

		// inst amount result (subtract)
		case instructions.LOOP: {
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
