package compiler

import (
	"fmt"
	"reflect"
	"vm-go/ast"
	"vm-go/chunk"
	"vm-go/token"
	"vm-go/util"
	"vm-go/value"
)

/*
	If/Else
	Control Flow:

		[ condition ]

	+--	OP_JUMP_FALSE
	|	OP_POP
	|
	|	[ then branch ]
	|
	|	OP_JUMP --------+
	+-> OP_POP			|
						|
		[ else branch ] |
						|
	continues... <------+
*/

// 'then' and 'else_' are functions because this function accepts both statements and expressions, and functions that
// compile these things don't return and have side effects, so the caller creates a new function and inserts the
// code inside it to generate the desired branch. 'else_' is a pointer because it's nullable. If it's not present,
// please pass in 'nil' for it.
func (c *Compiler) compileIf(condition ast.Expression, then func(), else_ *func(), pos token.Position) {
	c.expression(condition)
	c.writeBytePos(OP_JUMP_FALSE, pos)

	jumpFalseOffsetIndex := len(c.chunk.Code)
	c.writeBytes(util.IntToBytes(0)) // dummy
	c.writeBytePos(OP_POP, pos)

	then()

	if else_ != nil {
		c.writeBytePos(OP_JUMP, pos)
		jumpOffsetIndex := len(c.chunk.Code)
		c.writeBytes(util.IntToBytes(0)) // dummy

		// insert the real offset into the instruction, right before OP_POP
		c.backpatch(jumpFalseOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpFalseOffsetIndex - 4)) // index

		c.writeBytePos(OP_POP, pos)
		(*else_)()

		c.backpatch(jumpOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpOffsetIndex - 4)) // index
	} else {
		// insert the real offset into the instruction, if there's no else
		c.backpatch(jumpFalseOffsetIndex, util.IntToBytes(len(c.chunk.Code) - jumpFalseOffsetIndex - 4)) // index
	}
}

func (c *Compiler) compileFunction(parameters []ast.Parameter, body ast.BlockStatement, name *string, pos token.Position) {
	fnCompiler := newFnCompiler(body.Stmts, c)

	for _, param := range parameters {
		fnCompiler.addVariable(param.Name, param.Name.Pos)
	}

	fnChunk, hadError := fnCompiler.compileFnBody(pos)

	if hadError {
		c.hadError = true
		return
	}

	function := value.ValueFunction{
		Arity: len(parameters),
		Chunk: fnChunk,
		Name: name,
	}

	index := c.addConstant(function)
	c.writeBytePos(OP_PUSH_CLOSURE, pos)
	c.writeBytes(util.IntToBytes(index))

	c.writeBytes(util.IntToBytes(len(fnCompiler.upvalues)))

	// emit upvalue data
	// structure: 0/1 | index
	for i := range len(fnCompiler.upvalues) {
		if fnCompiler.upvalues[i].isLocal {
			c.writeBytePos(1, pos)
		} else {
			c.writeBytePos(0, pos)
		}

		c.writeBytes(util.IntToBytes(fnCompiler.upvalues[i].index))
	}
}

func (c *Compiler) compileFnBody(pos token.Position) (chunk.Chunk, bool) {
	c.statements(c.ast)

	if len(c.chunk.Code) > 0 && c.chunk.Code[len(c.chunk.Code) - 1] != OP_RETURN {
		c.endScope(pos)

		c.writeBytePos(OP_VOID, pos)
		c.writeBytePos(OP_RETURN, pos)
	}

	return c.chunk, c.hadError
}

func (c *Compiler) writeByte(b byte) {
	c.chunk.Code = append(c.chunk.Code, b)
}

func (c *Compiler) writeBytePos(b byte, pos token.Position) {
	c.chunk.Positions = append(c.chunk.Positions, pos)
	c.chunk.Code = append(c.chunk.Code, b)
}

func (c *Compiler) writeBytes(bytes []byte) {
	c.chunk.Code = append(c.chunk.Code, bytes...)

	// append dummy positions in the positions array
	for range len(bytes) {
		c.chunk.Positions = append(c.chunk.Positions, token.Position{})
	}
}

func (c *Compiler) backpatch(index int, bytes []byte) {
	// Ensure the position is valid
	if index < 0 || index + len(bytes) > len(c.chunk.Code) {
		// TODO: separate this into a function
		fmt.Printf("internal: invalid position: %d\n", index)
		c.hadError = true
		return
	}

	// Overwrite the bytes at the specified position
	for i, b := range bytes {
		c.chunk.Code[index + i] = b
	}
}

func (c *Compiler) addDeclarationInstruction(pos token.Position) {
	if c.scopeDepth == 0 {
		c.writeBytePos(OP_DEF_GLOBAL, pos)
	} else {
		c.writeBytePos(OP_DEF_LOCAL, pos)
	}
}

func (c *Compiler) resolveVariable(token token.Token, set bool) (int, Opcode) {
	// search it in locals
	index, opcode := c.resolveLocal(token, set)

	// found it in locals.
	if index != -1 {
		return index, opcode
	}

	// didn't find inside locals, let's search it in enclosing compilers for upvalues
	index, opcode = c.resolveUpvalue(token, set)

	// found it in enclosing compilers.
	if index != -1 {
		return index, opcode
	}

	// didn't find in enclosing compilers, let's search it in globals
	index, opcode = c.resolveGlobal(token, set)

	// found it in globals.
	if index != -1 {
		return index, opcode
	}

	// the variable doesn't exist.
	c.error(token.Pos, len(token.Lexeme), fmt.Sprintf("'%s' doesn't exist in this or in a parent scope.", token.Lexeme))
	return -1, OP_GET_LOCAL
}

func (c *Compiler) resolveLocal(token token.Token, set bool) (int, Opcode) {
	for i := len(c.locals) - 1; i >= 0; i-- {
		if c.locals[i].name.Lexeme == token.Lexeme {
			var opcode Opcode

			if set {
				opcode = OP_SET_LOCAL
			} else {
				opcode = OP_GET_LOCAL
			}

			return i, opcode
		}
	}

	return -1, OP_GET_LOCAL
}

func (c *Compiler) resolveUpvalue(token token.Token, set bool) (int, Opcode) {
	if c.enclosing == nil {
		return -1, OP_GET_UPVALUE
	}

	var opcode Opcode

	if set {
		opcode = OP_SET_UPVALUE
	} else {
		opcode = OP_GET_UPVALUE
	}

	// search it in enclosing's locals.
	// the opcode is not necessary
	index, _ := c.enclosing.resolveLocal(token, set)

	// found it in enclosing's locals.
	if index != -1 {
		// mark the variable as captured, so the compiler can emit an instruction
		// to close the upvalues that reference it, when it goes out of scope.
		c.enclosing.locals[index].isCaptured = true

		upIndex := c.addUpvalue(index, true)
		return upIndex, opcode
	}

	// didn't find inside the locals of the enclosing function. let's search it in its enclosing recursively.
	index, _ = c.enclosing.resolveUpvalue(token, set)

	// found it. let's capture that upvalue and add it to the list of upvalues of this function.
	if index != -1 {
		upIndex := c.addUpvalue(index, false)
		return upIndex, opcode
	}

	// didn't find it in any enclosing function.
	return -1, OP_GET_UPVALUE
}

func (c *Compiler) resolveGlobal(token token.Token, set bool) (int, Opcode) {
	for i := len(c.globals) - 1; i >= 0; i-- {
		if c.globals[i].name.Lexeme == token.Lexeme {
			// the scope depth is also verified, because if the compiler is in an inner scope, the global is
			// guaranteed to be initialized, because the program always starts at main(), and it is called after
			// all globals are initialized.
			if !c.globals[i].initialized && c.scopeDepth == 0 {
				c.error(token.Pos, len(token.Lexeme), fmt.Sprintf("'%s' is used before being initialized.", token.Lexeme))
			}

			var opcode Opcode

			if set {
				opcode = OP_SET_GLOBAL
			} else {
				opcode = OP_GET_GLOBAL
			}

			return i, opcode
		}
	}

	return -1, OP_GET_GLOBAL
}

func (c *Compiler) addUpvalue(index int, isLocal bool) int {
	// check if an upvalue to the same variable already exists
	// if so, return it
    for i, upvalue := range c.upvalues {
        if upvalue.index == index && upvalue.isLocal == isLocal {
            return i
        }
    }

	// if not, add a new one and return its index (which is len(c.upvalues) - 1),
	// because it's the last element.
    c.upvalues = append(c.upvalues, Upvalue{
        index:   index,
        isLocal: isLocal,
    })

    return len(c.upvalues) - 1
}

func (c *Compiler) addVariable(token token.Token, pos token.Position) {
	// Find the variable to check if it already exists or not, in this scope
	index := -1
	for i := len(c.locals) - 1; i >= 0; i-- {
		if c.locals[i].name.Lexeme == token.Lexeme {
			index = i
			break
		}

		if c.locals[i].depth < c.scopeDepth {
			break
		}
	}

	// It didn't find it, we can declare it safely, if it's a local.
	if index == -1 {
		if c.scopeDepth == 0 {
			// If it's a global, it won't be in the locals list, and we don't need to declare it again,
			// because when we hoisted them, we have already declared them.

			// We just need to check for redeclaration, and mark it as initialized if it's not;
			for i := len(c.globals) - 1; i >= 0; i-- {
				if c.globals[i].name.Lexeme == token.Lexeme {
					if c.globals[i].initialized {
						// It's a redeclaration, so we throw an error.
						// this message is for the global scope.
						c.error(pos, len(c.globals[i].name.Lexeme), fmt.Sprintf("'%s' has already been declared in this scope.", token.Lexeme))
						return
					} else {
						// It's not; mark it as initialized.
						c.globals[i].initialized = true
					}
				}
			}
		} else {
			// It's a local, so we declare it.
			c.locals = append(c.locals, Local{
				name:        token,
				depth:       c.scopeDepth,
				isCaptured:  false,
			})
		}
	} else {
		// We found the variable, it can be in the same scope or not
		existing := c.locals[index]

		// If it's in the same scope, throw an error, because it's a redeclaration.
		// this message is for local scopes.
		if existing.depth == c.scopeDepth {
			c.error(pos, len(existing.name.Lexeme), fmt.Sprintf("'%s' has already been declared in this scope.", token.Lexeme))
			return
		} else {
			// The variable is in an enclosing scope, we'll shadow it by declaring it in this scope
			c.locals = append(c.locals, Local{
				name:        token,
				depth:       c.scopeDepth,
				isCaptured:  false,
			})
		}
	}
}

func (c *Compiler) block(stmts []ast.Statement, pos token.Position) {
	c.beginScope()
	c.statements(stmts)
	c.endScope(pos)
}

func (c *Compiler) beginScope() {
	c.scopeDepth += 1
}

func (c *Compiler) endScope(pos token.Position) {
	currentScopeDepth := c.scopeDepth
	c.scopeDepth -= 1

	if len(c.locals) > 0 {
		count := 0
		realCount := 0

		for i := len(c.locals) - 1; i >= 0; i-- {
			local := &c.locals[i]

			// If the local variable is not in the current scope, stop popping.
			if local.depth != currentScopeDepth {
				break
			}

			if local.isCaptured {
				// Pop the other variables, which were counted before this captured one.
				c.emitPop(count, pos)

				// Pop the captured variable and reset the count for further counting.
				c.writeBytePos(OP_CLOSE_UPVALUE, pos)
				count = 0
				realCount++
			} else {
				// Count non-captured variables for potential batch popping.
				count++
				realCount++
			}
		}

		// Remove the variables from locals and pop the remaining variables, if any.
		c.locals = c.locals[:len(c.locals)-realCount]
		c.emitPop(count, pos)
	}
}

// ---

func (c *Compiler) emitPop(count int, pos token.Position) {
	// If count <= 0 this function does nothing.

	if count > 1 {
		c.writeBytePos(OP_POPN_LOCAL, pos)
		c.writeBytes(util.IntToBytes(count))
	} else if count == 1 {
		c.writeBytePos(OP_POP_LOCAL, pos)
	}
}

func (c *Compiler) addConstant(v value.Value) int {
	for i, constant := range c.chunk.Constants {
		if reflect.DeepEqual(constant, v) {
			return i
		}
	}

	c.chunk.Constants = append(c.chunk.Constants, v)
	return len(c.chunk.Constants) - 1
}

func (c *Compiler) error(pos token.Position, length int, message string) {
	if c.hadError {
		return
	}

	util.Error(pos, length, message, c.fileData)

	c.hadError = true
	c.panicMode = true
}

func (c *Compiler) errorNoBody(message string) {
	if c.hadError {
		return
	}

	fmt.Printf("[-] Error at %s: %s\n", c.fileData.Name, message)

	c.hadError = true
	c.panicMode = true
}
