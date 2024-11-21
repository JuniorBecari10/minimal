package compiler

import (
	"bytes"
	"vm-go/ast"
	"vm-go/token"
	"vm-go/value"
)

type Opcode int

const (
	OP_PUSH_CONST = iota
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV

	OP_DEF_VAR
	OP_GET_VAR
	OP_SET_VAR

	OP_POP
	OP_POP_VAR
	OP_POPN_VAR

	OP_JUMP
	OP_JUMP_FALSE
	OP_LOOP

	OP_EQUAL
	OP_NOT_EQUAL

	OP_GREATER
	OP_GREATER_EQUAL

	OP_LESS
	OP_LESS_EQUAL

	OP_AND
	OP_OR
	OP_XOR

	OP_NOT
	OP_NEGATE

	OP_PRINT
)

type Variable struct {
	name token.Token
	depth int
	initialized bool
}

type Compiler struct {
	program []ast.Statement
	variables []Variable

	scopeDepth int
	constants []value.Value

	hadError bool
	panicMode bool
}

func NewCompiler(program []ast.Statement) *Compiler {
	return &Compiler{
		program: program,
		variables: []Variable{},

		scopeDepth: 0,
		constants: []value.Value{},

		hadError: false,
		panicMode: false,
	}
}

func (c *Compiler) Compile() ([]byte, []value.Value, bool) {
	c.hoistTopLevel()
	
	bytecode := c.statements(c.program)
	return bytecode, c.constants, c.hadError
}

// ---

func (c *Compiler) statements(stmts []ast.Statement) []byte {
	res := bytes.Buffer{}

	for _, stmt := range stmts {
		if c.panicMode {
			c.panicMode = false
		}

		res.WriteString(string(c.statement(stmt)))
	}

	return res.Bytes()
}

func (c *Compiler) hoistTopLevel() {
	for _, decl := range c.program {
		switch s := decl.(type) {
			case ast.VarStatement: {
				c.variables = append(c.variables, Variable{
					name: s.Name,
					depth: c.scopeDepth,
					initialized: false,
				})
			}
		}
	}
}
