package analyzer

import (
	"minc/parser"
	"minc/ast"
	"minc/tast"
)

type AnalyzerResult = parser.ParserResult

const (
	RES_OK AnalyzerResult = iota
	RES_ERROR
)

type Analyzer struct {

}

func New(ast []ast.Statement) *Analyzer {
	
}

// Analyzes the code, does a lot of checkings, and if necessary or it's better, modify it, and returns a typed AST.
func (a *Analyzer) Analyze() ([]tast.Statement, AnalyzerResult) {
	
}
