package disassembler

import "minlib/instructions"

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
		case instructions.PUSH_CONST:
			return "PUSH_CONST"

		case instructions.PUSH_CLOSURE:
			return "PUSH_CLOSURE"

		case instructions.APPEND_METHODS:
			return "APPEND_METHODS"

		case instructions.ADD:
			return "ADD"
		case instructions.SUB:
			return "SUB"
		case instructions.MUL:
			return "MUL"
		case instructions.DIV:
			return "DIV"
		
		case instructions.MOD:
			return "MOD"

		case instructions.DEF_LOCAL:
			return "DEF_LOCAL"
		case instructions.GET_LOCAL:
			return "GET_LOCAL"
		case instructions.SET_LOCAL:
			return "SET_LOCAL"

		case instructions.GET_UPVALUE:
			return "GET_UPVALUE"
		case instructions.SET_UPVALUE:
			return "SET_UPVALUE"

		case instructions.DEF_GLOBAL:
			return "DEF_GLOBAL"
		case instructions.GET_GLOBAL:
			return "GET_GLOBAL"
		case instructions.SET_GLOBAL:
			return "SET_GLOBAL"

		case instructions.GET_PROPERTY:
			return "GET_PROPERTY"
		case instructions.SET_PROPERTY:
			return "SET_PROPERTY"

		case instructions.POP:
			return "POP"
		case instructions.POP_LOCAL:
			return "POP_LOCAL"
		case instructions.POPN_LOCAL:
			return "POPN_LOCAL"

		case instructions.CLOSE_UPVALUE:
			return "CLOSE_UPVALUE"

		case instructions.JUMP:
			return "JUMP"
		case instructions.JUMP_TRUE:
			return "JUMP_TRUE"
		case instructions.JUMP_FALSE:
			return "JUMP_FALSE"
		case instructions.JUMP_HAS_NO_NEXT:
			return "JUMP_HAS_NO_NEXT"
		case instructions.LOOP:
			return "LOOP"

		case instructions.EQUAL:
			return "EQUAL"
		case instructions.NOT_EQUAL:
			return "NOT_EQUAL"

		case instructions.GREATER:
			return "GREATER"
		case instructions.GREATER_EQUAL:
			return "GREATER_EQUAL"

		case instructions.LESS:
			return "LESS"
		case instructions.LESS_EQUAL:
			return "LESS_EQUAL"

		case instructions.AND:
			return "AND"
		case instructions.OR:
			return "OR"
		
		case instructions.NOT:
			return "NOT"
		
		case instructions.NEGATE:
			return "NEGATE"

		case instructions.CALL:
			return "CALL"
		case instructions.CALL_PROPERTY:
			return "CALL_PROPERTY"
		case instructions.RETURN:
			return "RETURN"

		case instructions.PUSH_TRUE:
			return "PUSH_TRUE"
		case instructions.PUSH_FALSE:
			return "PUSH_FALSE"

		case instructions.PUSH_NIL:
			return "PUSH_NIL"
		case instructions.PUSH_VOID:
			return "PUSH_VOID"

		case instructions.MAKE_RANGE:
			return "MAKE_RANGE"
		case instructions.MAKE_INCL_RANGE:
			return "MAKE_INCL_RANGE"
		case instructions.MAKE_ITERATOR:
			return "MAKE_ITERATOR"

        case instructions.GET_NEXT:
            return "GET_NEXT"
        case instructions.ADVANCE:
            return "ADVANCE"

		case instructions.ASSERT_BOOL:
			return "ASSERT_BOOL"

		default:
			return "Unknown"
	}
}
