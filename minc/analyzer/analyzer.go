package analyzer

import (
	"minc/ast"
	"minc/parser"
	"minc/tast"
	"minlib/file"
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

type Analyzer struct {
	ast ast.Ast

	locals []Local
	globals []Global
	natives []Global

	scopeDepth uint32
	isInsideLoop bool

	fileData *file.FileData
}

func New(ast ast.Ast, fileData *file.FileData) *Analyzer {
	return &Analyzer{
		ast: ast,
		
		locals: []Local{},
		globals: []Global{},
		natives: []Global{},

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

	tast, res := a.analyzeTopLevelStatements()
	
	endRes := a.endTopLevel(); if endRes == RES_ERROR {
		res = endRes
	}

	return tast, res
}
