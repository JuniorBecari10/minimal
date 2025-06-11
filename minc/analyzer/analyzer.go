package analyzer

import (
	"fmt"
	"minc/ast"
	"minc/parser"
	"minc/tast"
	"minc/types"
	"minlib/file"
	"minlib/token"
)

type AnalyzerResult = parser.ParserResult

const (
	RES_OK AnalyzerResult = iota
	RES_ERROR
)

type Local struct {
	name token.Token
	localType types.Type

	depth uint32
	isCaptured bool
	immutable bool
}

type Global struct {
	name token.Token
	globalType types.Type

	immutable bool
	initialized bool // to check if it has already been declared or just hoisted
}

type Upvalue struct {
	index int
	isLocal bool
}

type Analyzer struct {
	ast ast.Ast

	locals []Local
	globals []Global
	upvalues []Upvalue

	scopeDepth uint32
	isInsideLoop bool

	fileData *file.FileData
}

func New(ast ast.Ast, fileData *file.FileData) *Analyzer {
	return &Analyzer{
		ast: ast,
		
		locals: []Local{},
		globals: []Global{},
		upvalues: []Upvalue{},

		scopeDepth: 0,
		isInsideLoop: false,

		fileData: fileData,
	}
}

// Analyzes the code, does a lot of checkings, and if necessary or it's better, modify it, and returns a typed AST.
func (a *Analyzer) Analyze() (tast.Tast, AnalyzerResult) {
	a.addNatives()

	res := a.hoistTopLevel(); if res != RES_OK {
		return nil, res
	}

	return a.analyzeBlock(a.ast)
}

// pre-declares all top-level declarations in order to properly do name resolution.
// it doesn't do all the required setup, it just adds the name to the list with its type, using shallow inference.
func (a *Analyzer) hoistTopLevel() AnalyzerResult {
	res := RES_OK

	for _, d := range a.ast {
		switch decl := d.Data.(type) {
			// In 'fn' statements we check only the declaration. All types must be concrete.
			case ast.FnDeclaration: {
				a.globals = append(a.globals, Global{
					name: decl.Name,
					globalType: types.TypeFunction{},

					immutable: true,
					initialized: false,
				})
			}
			
			// In 'var'/'let' declarations we check the type and if omitted, we try to infer it shallowly.
			case ast.VarDeclaration: {
				if decl.Type == nil {
					// type isn't annotated. infer it shallowly.

				} else {
					// type is annotated; add the variable with its type.
					// actual type checking is done later.
					a.globals = append(a.globals, Global {
						name: decl.Name,
						globalType: *decl.Type,

						immutable: decl.Immutable,
						initialized: true,
					})
				}
			}

			// In records, all types must be explicitly annotated and concrete.
			case ast.RecordDeclaration: {

			}

			// Should not reach here.
			default:
				panic(fmt.Sprintf("Unknown declaration %v of type %T", decl, decl))
		}
	}

	return res
}

// Adds native functions and variables to the global scope.
func (a *Analyzer) addNatives() {
	// fn print()
	a.globals = append(a.globals, newNative("print", types.TypeFunction{
		Parameters: []types.Type{},
		Return: types.TypeVoid{},
	}))
	
	// fn println()
	a.globals = append(a.globals, newNative("println", types.TypeFunction{
		Parameters: []types.Type{},
		Return: types.TypeVoid{},
	}))
	
	// fn input(prompt: str): str
	a.globals = append(a.globals, newNative("print", types.TypeFunction{
		Parameters: []types.Type{ types.TypeStr{} },
		Return: types.TypeStr{},
	}))
	
	// fn time(): int
	a.globals = append(a.globals, newNative("time", types.TypeFunction{
		Parameters: []types.Type{},
		Return: types.TypeVoid{},
	}))
}
