package analyzer

import (
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

type BlockAnalyzeMode int

const (
	MODE_NORMAL BlockAnalyzeMode = iota
	MODE_FUNCTION
	MODE_LOOP
)

// TODO: add flags canBeMutable (for functions and refine the error message)
type Local struct {
	name token.Token
	localType types.Type

	immutable bool
	depth uint32
	modified bool
	used bool
}

type Global struct {
	name token.Token
	globalType types.Type

	immutable bool
	initialized bool // to check if it has already been declared or just hoisted
	modified bool
	used bool
}

type Analyzer struct {
	ast ast.Ast

	locals []Local
	globals []Global

	scopeDepth uint32
	isInsideLoop bool

	fileData *file.FileData
}

func New(ast ast.Ast, fileData *file.FileData) *Analyzer {
	return &Analyzer{
		ast: ast,
		
		locals: []Local{},
		globals: []Global{},

		scopeDepth: 0,
		isInsideLoop: false,

		fileData: fileData,
	}
}

// Analyzes the code, does a lot of checkings, and if necessary or it's better, modifies it,
// and returns a typed AST, ready to be compiled.
func (a *Analyzer) Analyze() (tast.Tast, AnalyzerResult) {
	a.addNatives()

	res := a.hoistTopLevel(); if res != RES_OK {
		return nil, res
	}

	block, res := a.analyzeBlock(a.ast, MODE_NORMAL)
	
	endRes := a.endTopLevel(); if endRes == RES_ERROR {
		res = endRes
	}

	return block.Stmts, res
}
