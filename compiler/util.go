package compiler

import (
	"bytes"
	"fmt"
	"vm-go/token"
	"vm-go/util"
	"vm-go/value"
)

func (c *Compiler) resolveVariable(token token.Token) int {
	for i := len(c.variables) - 1; i >= 0; i-- {
		if c.variables[i].name.Lexeme == token.Lexeme {
			if !c.variables[i].initialized {
				// TODO check if the scopeDepth is not 0 and allow its use,
				// since functions are allowed in top level and the use of variables inside them is allowed,
				// because it is guaranteed to have them defined, since we'll have a main function
				c.error(token.Pos, fmt.Sprintf("'%s' is not defined yet", token.Lexeme))
				return -1
			}

			return i
		}
	}

	c.error(token.Pos, fmt.Sprintf("'%s' doesn't exist", token.Lexeme))
	return -1
}

func (c *Compiler) writeByte(b uint8) {
	c.chunk.Code = append(c.chunk.Code, b)
}

func (c *Compiler) writeBytePos(b uint8, pos token.Position) {
	c.chunk.Positions = append(c.chunk.Positions, pos)
	c.chunk.Code = append(c.chunk.Code, b)
}

func (c *Compiler) writeBytes(bytes []uint8) {
	for _, b := range bytes {
		c.chunk.Code = append(c.chunk.Code, b)
	}
}

func (c *Compiler) backpatch(index int, bytes []uint8) {
	// Ensure the position is valid
	if index < 0 || index+len(bytes) > len(c.chunk.Code) {
		// TODO: separate this into a function
		fmt.Printf("internal: invalid position: %d\n", index)
		c.hadError = true
		return
	}

	// Overwrite the bytes at the specified position
	for i, b := range bytes {
		c.chunk.Code[index+i] = b
	}
}

func (c *Compiler) beginScope() {
	c.scopeDepth += 1
}

// this returns the instructions to pop the variables from the variable stack
func (c *Compiler) endScope(pos token.Position) []byte {
	res := bytes.Buffer{}
	c.scopeDepth -= 1

	if len(c.variables) > 0 {
		lastDepth := c.variables[len(c.variables) - 1].depth

		count := 0
		for i := len(c.variables) - 1; i >= 0; i-- {
			if c.variables[i].depth != lastDepth {
				break
			}

			count += 1
		}

		c.variables = c.variables[:len(c.variables)-count]
		c.chunk.Positions = append(c.chunk.Positions, pos)

		if count > 1 {
			res.WriteByte(OP_POPN_VAR)
			res.WriteString(string(util.IntToBytes(count)))
		} else {
			res.WriteByte(OP_POP_VAR)
		}
	}

	return res.Bytes()
}

// ---

func (c *Compiler) AddConstant(v value.Value) int {
	for i, constant := range c.chunk.Constants {
		if constant == v {
			return i
		}
	}

	c.chunk.Constants = append(c.chunk.Constants, v)
	return len(c.chunk.Constants) - 1
}

func (c *Compiler) error(pos token.Position, message string) {
	if c.hadError {
		return
	}

	util.Error(pos, message, c.fileData)

	c.hadError = true
	c.panicMode = true
}
