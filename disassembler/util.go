package disassembler

import "vm-go/compiler"

func (d *Disassembler) nextByte() byte {
	ip := d.ip
	d.ip += 1

	return d.chunk.Code[ip]
}

func (d *Disassembler) isAtEnd() bool {
	return d.ip >= len(d.chunk.Code)
}

// ---

func getInstructionName(inst byte) string {
	switch inst {
	case compiler.OP_PUSH_CONST:
		return "OP_PUSH_CONST"

	case compiler.OP_ADD:
		return "OP_ADD"
	case compiler.OP_SUB:
		return "OP_SUB"
	case compiler.OP_MUL:
		return "OP_MUL"
	case compiler.OP_DIV:
		return "OP_DIV"
	
	case compiler.OP_MODULO:
		return "OP_MODULO"

	case compiler.OP_DEF_VAR:
		return "OP_DEF_VAR"
	case compiler.OP_GET_VAR:
		return "OP_GET_VAR"
	case compiler.OP_SET_VAR:
		return "OP_SET_VAR"

	case compiler.OP_POP:
		return "OP_POP"
	case compiler.OP_POP_VAR:
		return "OP_POP_VAR"
	case compiler.OP_POPN_VAR:
		return "OP_POPN_VAR"

	case compiler.OP_JUMP:
		return "OP_JUMP"
	case compiler.OP_JUMP_FALSE:
		return "OP_JUMP_FALSE"
	case compiler.OP_LOOP:
		return "OP_LOOP"

	case compiler.OP_EQUAL:
		return "OP_EQUAL"
	case compiler.OP_NOT_EQUAL:
		return "OP_NOT_EQUAL"

	case compiler.OP_GREATER:
		return "OP_GREATER"
	case compiler.OP_GREATER_EQUAL:
		return "OP_GREATER_EQUAL"

	case compiler.OP_LESS:
		return "OP_LESS"
	case compiler.OP_LESS_EQUAL:
		return "OP_LESS_EQUAL"

	case compiler.OP_AND:
		return "OP_AND"
	case compiler.OP_OR:
		return "OP_OR"
	case compiler.OP_XOR:
		return "OP_XOR"
	
	case compiler.OP_NOT:
		return "OP_NOT"
	
	case compiler.OP_NEGATE:
		return "OP_NEGATE"

	case compiler.OP_RETURN:
		return "OP_RETURN"

	case compiler.OP_TRUE:
		return "OP_TRUE"
	case compiler.OP_FALSE:
		return "OP_FALSE"

	case compiler.OP_NIL:
		return "OP_NIL"
	case compiler.OP_VOID:
		return "OP_VOID"

	case compiler.OP_CALL:
		return "OP_CALL"
	case compiler.OP_PRINT:
		return "OP_PRINT"

	default:
		return "Unknown"
	}
}
