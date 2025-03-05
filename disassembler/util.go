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
			return "PUSH_CONST"

		case compiler.OP_PUSH_CLOSURE:
			return "PUSH_CLOSURE"

		case compiler.OP_APPEND_METHODS:
			return "APPEND_METHODS"

		case compiler.OP_ADD:
			return "ADD"
		case compiler.OP_SUB:
			return "SUB"
		case compiler.OP_MUL:
			return "MUL"
		case compiler.OP_DIV:
			return "DIV"
		
		case compiler.OP_MODULO:
			return "MODULO"

		case compiler.OP_DEF_LOCAL:
			return "DEF_LOCAL"
		case compiler.OP_GET_LOCAL:
			return "GET_LOCAL"
		case compiler.OP_SET_LOCAL:
			return "SET_LOCAL"

		case compiler.OP_GET_UPVALUE:
			return "GET_UPVALUE"
		case compiler.OP_SET_UPVALUE:
			return "SET_UPVALUE"

		case compiler.OP_DEF_GLOBAL:
			return "DEF_GLOBAL"
		case compiler.OP_GET_GLOBAL:
			return "GET_GLOBAL"
		case compiler.OP_SET_GLOBAL:
			return "SET_GLOBAL"

		case compiler.OP_GET_PROPERTY:
			return "GET_PROPERTY"
		case compiler.OP_SET_PROPERTY:
			return "SET_PROPERTY"

		case compiler.OP_POP:
			return "POP"
		case compiler.OP_POP_LOCAL:
			return "POP_LOCAL"
		case compiler.OP_POPN_LOCAL:
			return "POPN_LOCAL"

		case compiler.OP_CLOSE_UPVALUE:
			return "CLOSE_UPVALUE"

		case compiler.OP_JUMP:
			return "JUMP"
		case compiler.OP_JUMP_TRUE:
			return "JUMP_TRUE"
		case compiler.OP_JUMP_FALSE:
			return "JUMP_FALSE"
		case compiler.OP_JUMP_HAS_NO_NEXT:
			return "JUMP_HAS_NO_NEXT"
		case compiler.OP_LOOP:
			return "LOOP"

		case compiler.OP_EQUAL:
			return "EQUAL"
		case compiler.OP_NOT_EQUAL:
			return "NOT_EQUAL"

		case compiler.OP_GREATER:
			return "GREATER"
		case compiler.OP_GREATER_EQUAL:
			return "GREATER_EQUAL"

		case compiler.OP_LESS:
			return "LESS"
		case compiler.OP_LESS_EQUAL:
			return "LESS_EQUAL"

		case compiler.OP_AND:
			return "AND"
		case compiler.OP_OR:
			return "OR"
		
		case compiler.OP_NOT:
			return "NOT"
		
		case compiler.OP_NEGATE:
			return "NEGATE"

		case compiler.OP_CALL:
			return "CALL"
		case compiler.OP_CALL_PROPERTY:
			return "CALL_PROPERTY"
		case compiler.OP_RETURN:
			return "RETURN"

		case compiler.OP_PUSH_TRUE:
			return "PUSH_TRUE"
		case compiler.OP_PUSH_FALSE:
			return "PUSH_FALSE"

		case compiler.OP_PUSH_NIL:
			return "PUSH_NIL"
		case compiler.OP_PUSH_VOID:
			return "PUSH_VOID"

		case compiler.OP_MAKE_RANGE:
			return "MAKE_RANGE"
		case compiler.OP_MAKE_INCL_RANGE:
			return "MAKE_INCL_RANGE"
		case compiler.OP_MAKE_ITERATOR:
			return "MAKE_ITERATOR"

        case compiler.OP_GET_NEXT:
            return "GET_NEXT"
        case compiler.OP_ADVANCE:
            return "ADVANCE"

		case compiler.OP_ASSERT_BOOL:
			return "ASSERT_BOOL"

		default:
			return "Unknown"
	}
}
