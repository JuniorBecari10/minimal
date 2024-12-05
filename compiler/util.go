package compiler

import (
	"bytes"
	"fmt"
	"vm-go/ast"
	"vm-go/chunk"
	"vm-go/token"
	"vm-go/util"
	"vm-go/value"
)

func (c *Compiler) resolveVariable(token token.Token) int {
	for i := len(c.locals) - 1; i >= 0; i-- {
		if c.locals[i].name.Lexeme == token.Lexeme {
			if !c.locals[i].initialized {
				// TODO check if the scopeDepth is not 0 and allow its use,
				// since functions are allowed in top level and the use of variables inside them is allowed,
				// because it is guaranteed to have them defined, since we'll have a main function
				c.error(token.Pos, len(token.Lexeme), fmt.Sprintf("'%s' is not defined yet", token.Lexeme))
				return -1
			}

			return i
		}
	}

	c.error(token.Pos, len(token.Lexeme), fmt.Sprintf("'%s' doesn't exist", token.Lexeme))
	return -1
}

func (c *Compiler) compileFnBody(pos token.Position) (chunk.Chunk, bool) {
	c.statements(c.ast)

	c.writeBytePos(OP_VOID, pos)
	c.writeBytePos(OP_RETURN, pos)
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

	// didn't find it
	if index == -1 {
		c.locals = append(c.locals, Local{
			name:        token,
			depth:       c.scopeDepth,
			initialized: true,
		})
		// found it
	} else {
		// if it's a global and we've reached its declaration. also make sure that it isn't a redeclaration
		// we do that by checking if the initialized field is true
		if c.locals[index].depth == 0 && c.scopeDepth == 0 && !c.locals[index].initialized {
			c.locals[index].initialized = true
		} else {
			// Found a variable with the same name, it can be in the same scope or not
			existing := c.locals[index]

			if existing.depth == 0 && c.scopeDepth == 0 && !existing.initialized {
				// If it's a global variable and uninitialized, allow redeclaration
				c.locals[index].initialized = true
			} else if existing.depth == c.scopeDepth {
				// Redeclaration in the same scope is not allowed
				c.error(pos, len(existing.name.Lexeme), fmt.Sprintf("'%s' has already been declared in this scope", token.Lexeme))
				return
			} else {
				// the variable is in an enclosing scope, we'll shadow it
				c.locals = append(c.locals, Local{
					name:        token,
					depth:       c.scopeDepth,
					initialized: true,
				})
			}
		}
	}
}

func (c *Compiler) block(stmts []ast.Statement, pos token.Position) {
	c.beginScope()
	c.statements(stmts)
	c.writeBytes(c.endScope(pos))
}

func (c *Compiler) beginScope() {
	c.scopeDepth += 1
}

// This returns the instructions to pop the variables from the variable stack
func (c *Compiler) endScope(pos token.Position) []byte {
	res := bytes.Buffer{}
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
			res.WriteByte(OP_POPN_VAR)
			res.WriteString(string(util.IntToBytes(count)))
		} else if count == 1 {
			res.WriteByte(OP_POP_VAR)
		}
	}

	return res.Bytes()
}

// ---

func (c *Compiler) addConstant(v value.Value) int {
	for i, constant := range c.chunk.Constants {
		if constant == v {
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
