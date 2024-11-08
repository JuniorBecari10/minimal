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
	OP_PUSH_CONST = iota
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

type Variable struct {
	name token.Token
	depth int
	initialized bool
}

type Compiler struct {
	program []ast.Statement
	locals []Variable
	scopeDepth int

	code bytes.Buffer
	constants []util.Value

	hadError bool
	panicMode bool
}

func NewCompiler(program []ast.Statement) *Compiler {
	return &Compiler{
		program: program,
		locals: []Variable{},
		scopeDepth: 0,

		code: bytes.Buffer{},
		constants: []util.Value{},

		hadError: false,
		panicMode: false,
	}
}

func (c *Compiler) Compile() (string, []util.Value, bool) {
	c.hoistTopLevel()
	for _, stmt := range c.program {
		if c.panicMode {
			c.panicMode = false
		}

		c.statement(stmt)
	}
	
	return c.code.String(), c.constants, c.hadError
}

// ---

func (c *Compiler) hoistTopLevel() {
	for _, decl := range c.program {
		switch s := decl.(type) {
			case ast.VarStatement:
				c.locals = append(c.locals, Variable{
					name: s.Name,
					depth: c.scopeDepth,
					initialized: false,
				})
		}
	}
}

func (c *Compiler) statement(stmt ast.Statement) {
	switch s := stmt.(type) {
		case ast.VarStatement: {
			index := -1
			for i := len(c.locals) - 1; i >= 0; i-- {
				if c.locals[i].name.Lexeme == s.Name.Lexeme {
					index = i
					break
				}
			}

			// didn't find
			if index == -1 {
				c.locals = append(c.locals, Variable{
					name: s.Name,
					depth: c.scopeDepth,
					initialized: true,
				})
			} else {
				// if it's a global and we're reached its declaratiomn. also make sure that it isn't a redeclaration
				if c.locals[index].depth == 0 && c.scopeDepth == 0 && !c.locals[index].initialized {
					c.locals[index].initialized = true
				} else {
					c.error(s.Pos, fmt.Sprintf("'%s' has already been declared in this scope", s.Name.Lexeme))
					return
				}
			}

			c.expression(s.Init)

			if c.hadError {
				return
			}

			c.emitByte(OP_DEF_VAR) // pop from stack and push to variable stack
		}

		case ast.BlockStatement: {
			c.beginScope()

			for _, stmt := range s.Stmts {
				if c.panicMode {
					c.panicMode = false
				}

				c.statement(stmt)
			}
			
			c.endScope()
		}

		case ast.PrintStatement: {
			c.expression(s.Expr)

			if c.hadError {
				return
			}

			c.emitByte(OP_PRINT)
		}

		case ast.ExprStatement:
			c.expression(s.Expr)
	}
}

// ---

func (c *Compiler) expression(expr ast.Expression) {
	if c.hadError {
		return
	}

	switch e := expr.(type) {
		case ast.NumberExpression: {
			index := c.AddConstant(util.Value(e.Literal))

			c.emitByte(OP_PUSH_CONST)
			c.emitBytes(util.IntToBytes(index)) // index has 4 bytes
		}

		case ast.IdentifierExpression: {
			for i := len(c.locals) - 1; i >= 0; i-- {
				if c.locals[i].name.Lexeme == e.Ident.Lexeme {
					if !c.locals[i].initialized {
						// TODO check if the scopeDepth is not 0 and allow its use,
						// since functions are allowed in top level and the use of variables inside them is allowed,
						// because it is guaranteed to have them defined
						c.error(e.Pos, fmt.Sprintf("'%s' is not defined yet", e.Ident.Lexeme))
						return
					}

					c.emitByte(OP_GET_VAR)
					c.emitBytes(util.IntToBytes(i)) // index has 4 bytes

					return
				}
			}

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

		c.locals = c.locals[:len(c.locals) - count]

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
	c.panicMode = true
}
