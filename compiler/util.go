package compiler

import (
	"bytes"
	"vm-go/token"
	"vm-go/util"
	"vm-go/value"
)

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

	util.Error(pos, message)

	c.hadError = true
	c.panicMode = true
}
