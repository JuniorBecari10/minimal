package analyzer

import (
	"minc/ast"
	"minc/tast"
	"minc/types"
)

func (a *Analyzer) analyzeBlock(block ast.Ast) (tast.Tast, AnalyzerResult) {
	generatedTast := make(tast.Tast, 0, len(block))

	for _, s := range block {
		switch stmt := s.Data.(type) {
			case ast.FnDeclaration: {
				
			}

			case ast.RecordDeclaration: {}

			case ast.ReturnStatement: {
				

                expr, res := a.analyzeExpression(*stmt.Expression)

				if res == RES_ERROR {
					return generatedTast, res
				}

				generatedTast = append(generatedTast, tast.ReturnStatement{
					Expression: expr,
				})
			}

			case ast.OutStatement: {}

			case ast.VarDeclaration: {}

			case ast.WhileStatement: {}

			case ast.ForStatement: {}

			case ast.ForVarStatement: {}
	
			case ast.BreakStatement: {}

			case ast.ContinueStatement: {}

			case ast.ExprStatement: {}
		}
	}

	return generatedTast, RES_OK
}

func (a *Analyzer) analyzeExpression(expr ast.Expression) (tast.Expression, AnalyzerResult) {

}
