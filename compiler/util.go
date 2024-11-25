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

func (c *Compiler) beginScope() {
	c.scopeDepth += 1
}

func (c *Compiler) endScope() []byte {
	res := bytes.Buffer{}
	c.scopeDepth -= 1

	if len(c.variables) > 0 {
		lastDepth := c.variables[len(c.variables)-1].depth

		count := 0
		for i := len(c.variables) - 1; i >= 0; i-- {
			if c.variables[i].depth != lastDepth {
				break
			}

			count += 1
		}

		c.variables = c.variables[:len(c.variables)-count]

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
	for i, constant := range c.constants {
		if constant == v {
			return i
		}
	}

	c.constants = append(c.constants, v)
	return len(c.constants) - 1
}

func (c *Compiler) error(pos token.Position, message string) {
	if c.hadError {
		return
	}

	util.Error(pos, message, c.fileData)

	c.hadError = true
	c.panicMode = true
}
