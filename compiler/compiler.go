package compiler

import (
	"vm-go/ast"
	"vm-go/chunk"
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
	OP_MODULO

	OP_DEF_LOCAL
	OP_GET_LOCAL
	OP_SET_LOCAL

	OP_DEF_GLOBAL
	OP_GET_GLOBAL
	OP_SET_GLOBAL

	OP_POP
	OP_POP_VAR
	OP_POPN_VAR

	OP_JUMP
	OP_JUMP_TRUE
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

	OP_NOT
	OP_NEGATE

	OP_RETURN

	OP_TRUE
	OP_FALSE
	
	OP_NIL
	OP_VOID

	OP_CALL
	OP_PRINT
)

type Local struct {
	name token.Token
	depth int
}

type Global struct {
	name token.Token
	initialized bool // to check redeclaration
}

type Compiler struct {
	ast []ast.Statement

	locals []Local
	globals []Global

	chunk chunk.Chunk
	scopeDepth int

	hadError bool
	panicMode bool

	fileData *util.FileData
}

func NewCompiler(ast []ast.Statement, fileData *util.FileData) *Compiler {
	return &Compiler{
		ast: ast,
		locals: []Local{},
		globals: []Global{},

		chunk: chunk.Chunk{},
		scopeDepth: 0,

		hadError: false,
		panicMode: false,

		fileData: fileData,
	}
}

func newFnCompiler(ast []ast.Statement, fileData *util.FileData, globals []Global, scopeDepth int) *Compiler {
	compiler := NewCompiler(ast, fileData)
	compiler.globals = globals
	compiler.scopeDepth = scopeDepth + 1

	return compiler
}

func (c *Compiler) Compile() (chunk.Chunk, bool) {
	c.addNativeFunctions()
	c.hoistTopLevel()

	c.statements(c.ast)
	c.callMain()

	return c.chunk, c.hadError
}

// ---

func (c *Compiler) statements(stmts []ast.Statement){
	for _, stmt := range stmts {
		if c.panicMode {
			c.panicMode = false
		}

		c.statement(stmt)
	}
}

// TODO: hoist the inner scopes too, to improve the error messages?
func (c *Compiler) hoistTopLevel() {
	for _, decl := range c.ast {
		switch s := decl.(type) {
			case ast.VarStatement: {
				c.globals = append(c.globals, Global{
					name: s.Name,
					initialized: false,
				})
			}
			case ast.FnStatement: {
				c.globals = append(c.globals, Global{
					name: s.Name,
					initialized: false,
				})
			}
		}
	}
}

func (c *Compiler) addNativeFunctions() {
	// they will be set to initialized to prevent shadowing in the global scope

	// fn print() -> void
	c.globals = append(c.globals, Global{
		name: token.Token{ Lexeme: "print" },
		initialized: true,
	})

	// fn println() -> void
	c.globals = append(c.globals, Global{
		name: token.Token{ Lexeme: "println" },
		initialized: true,
	})

	// fn time() -> number
	c.globals = append(c.globals, Global{
		name: token.Token{ Lexeme: "time" },
		initialized: true,
	})
}

func (c *Compiler) callMain() {
	for i, global := range c.globals {
		if global.name.Lexeme == "main" {
			// check if it's a function
			// in the meanwhile, this error will be caught at runtime

			c.writeBytePos(OP_GET_GLOBAL, global.name.Pos)
			c.writeBytes(util.IntToBytes(i))

			c.writeBytePos(OP_CALL, global.name.Pos)
			c.writeBytes(util.IntToBytes(0))

			return
		}
	}

	c.error(token.Position{}, 1, "A main function wasn't found")
}
