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

func (c *Compiler) compileFnBody(pos token.Position) (chunk.Chunk, bool) {
	c.statements(c.ast)

	if len(c.chunk.Code) > 0 && c.chunk.Code[len(c.chunk.Code) - 1] != OP_RETURN {
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

func (c *Compiler) resolveVariable(token token.Token) (int, Opcode) {
	for i := len(c.locals) - 1; i >= 0; i-- {
		if c.locals[i].name.Lexeme == token.Lexeme {
			return i, OP_GET_LOCAL
		}
	}
	// didn't find inside the locals, let's search it in globals
	for i := len(c.globals) - 1; i >= 0; i-- {
		if c.globals[i].name.Lexeme == token.Lexeme {
			// the scope depth is also verified, because if the compiler is in an inner scope, the global is
			// guaranteed to be initialized, because the program always starts at main(), and it is called after
			// all globals are initialized.
			if !c.globals[i].initialized && c.scopeDepth == 0 {
				c.error(token.Pos, len(token.Lexeme), fmt.Sprintf("'%s' is used before being initialized", token.Lexeme))
			}

			return i, OP_GET_GLOBAL
		}
	}

	// the variable doesn't exist
	c.error(token.Pos, len(token.Lexeme), fmt.Sprintf("'%s' doesn't exist", token.Lexeme))
	return -1, OP_GET_LOCAL
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
						c.error(pos, len(c.globals[i].name.Lexeme), fmt.Sprintf("'%s' has already been declared in this scope", token.Lexeme))
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
			})
		}
	} else {
		// We found the variable, it can be in the same scope or not
		existing := c.locals[index]

		// If it's in the same scope, throw an error, because it's a redeclaration
		if existing.depth == c.scopeDepth {
			c.error(pos, len(existing.name.Lexeme), fmt.Sprintf("'%s' has already been declared in this scope", token.Lexeme))
			return
		} else {
			// The variable is in an enclosing scope, we'll shadow it by declaring it in this scope
			c.locals = append(c.locals, Local{
				name:        token,
				depth:       c.scopeDepth,
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
		for i := len(c.locals) - 1; i >= 0; i-- {
			if c.locals[i].depth != currentScopeDepth {
				break
			}

			count += 1
		}

		c.locals = c.locals[:len(c.locals)-count]
		c.chunk.Positions = append(c.chunk.Positions, pos)

		if count > 1 {
			c.writeByte(OP_POPN_VAR)
			c.writeBytes(util.IntToBytes(count))
		} else if count == 1 {
			c.writeByte(OP_POP_VAR)
		}
	}
}

// ---

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
