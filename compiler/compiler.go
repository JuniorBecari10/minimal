package compiler

import (
	"bytes"
	"fmt"
	"vm-go/ast"
	"vm-go/token"
	"vm-go/util"
)

type Opcode int

const (
	OP_CONSTANT = iota
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV

	OP_DEF_VAR
	OP_GET_VAR
	OP_SET_VAR

	OP_POP_VAR
	OP_POPN_VAR

	OP_PRINT
)

type Local struct {
	name token.Token
	depth int
}

type Compiler struct {
	program []ast.Statement
	locals []Local
	scopeDepth int

	code bytes.Buffer
	constants []util.Value
	hadError bool
}

func NewCompiler(program []ast.Statement) *Compiler {
	return &Compiler{
		program: program,
		locals: []Local{},
		scopeDepth: 0,

		code: bytes.Buffer{},
		constants: []util.Value{},
		hadError: false,
	}
}

func (c *Compiler) Compile() (string, []util.Value) {
	for _, stmt := range c.program {
		c.statement(stmt)
	}
	
	return c.code.String(), c.constants
}

// ---

func (c *Compiler) statement(stmt ast.Statement) {
	switch s := stmt.(type) {
		case ast.VarStatement: {
			c.locals = append(c.locals, Local{
				name: s.Name,
				depth: c.scopeDepth,
			})

			c.expression(s.Init)
			c.emitByte(OP_DEF_VAR) // pop from stack and push to variable stack
		}

		case ast.BlockStatement: {
			c.beginScope()

			for _, stmt := range s.Stmts {
				c.statement(stmt)
			}
			
			c.endScope()
		}

		case ast.PrintStatement: {
			c.expression(s.Expr)
			c.emitByte(OP_PRINT)
		}

		case ast.ExprStatement:
			c.expression(s.Expr)
	}
}

// ---

func (c *Compiler) expression(expr ast.Expression) {
	switch e := expr.(type) {
		case ast.NumberExpression: {
			index := c.AddConstant(util.Value(e.Literal))

			c.emitByte(OP_CONSTANT)
			c.emitBytes(util.IntToBytes(index)) // index has 4 bytes
		}

		case ast.IdentifierExpression: {
			for i := len(c.locals) - 1; i >= 0; i-- {
				if c.locals[i].name.Lexeme == e.Ident.Lexeme {
					c.emitByte(OP_GET_VAR)
					c.emitBytes(util.IntToBytes(i)) // index has 4 bytes

					return
				}
			}

			// TODO: throw error and abort compilation if one happened
			c.error(e.Pos, fmt.Sprintf("'%s' doesn't exist", e.Ident.Lexeme))
		}

		case ast.BinaryExpression: {
			c.expression(e.Left)
			c.expression(e.Right)

			switch e.Operator.Kind {
				case token.TokenPlus: c.emitByte(OP_ADD)
				case token.TokenMinus: c.emitByte(OP_SUB)
				case token.TokenStar: c.emitByte(OP_MUL)
				case token.TokenSlash: c.emitByte(OP_DIV)

				default:
					panic(fmt.Sprintf("Unknown binary operator: '%s'", e.Operator.Kind))
			}
		}

		case ast.GroupExpression: c.expression(e.Expr)
	}
}

// ---

func (c *Compiler) beginScope() {
	c.scopeDepth += 1
}

func (c *Compiler) endScope() {
	c.scopeDepth -= 1

	if len(c.locals) > 0 {
		lastDepth := c.locals[len(c.locals) - 1].depth

		count := 0
		for i := len(c.locals) - 1; i >= 0; i-- {
			if c.locals[i].depth != lastDepth {
				break
			}

			count += 1
		}

		if count > 1 {
			c.emitByte(OP_POPN_VAR)
			c.emitBytes(util.IntToBytes(count)) // count has 4 bytes
		} else {
			c.emitByte(OP_POP_VAR)
		}
	}
}

// ---

func (c *Compiler) emitByte(b byte) {
	c.code.WriteByte(b)
}

func (c *Compiler) emitBytes(b []byte) {
	c.code.WriteString(string(b))
}

func (c *Compiler) AddConstant(v util.Value) int {
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
}
