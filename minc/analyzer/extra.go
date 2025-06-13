package analyzer

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
	"minc/types"
)

// synchronization point. this function prints the diagnostics and doesn't bubble them up
// this returns a block, with its type inferred by its statements.
// 'return' only tells the type of the block if this is the function block; otherwise it is never.
func (a *Analyzer) analyzeBlock(block ast.Ast, isFunction bool) (tast.BlockExpression, AnalyzerResult) {
	generatedTast := make(tast.Tast, 0, len(block))
	res := RES_OK

	var inferredType *types.TypeData = nil
	
	printDiag := func(diag diagnostic.Diagnostic) {
		diag.PrintDiagnostic()
		res = RES_ERROR
	}

	for _, s := range block {
		newStmt := func(data tast.StmtData) tast.Statement {
			return tast.Statement{
				Base: tast.AstBase(s.Base),
				Data: data,
			}
		}
		
		switch stmt := s.Data.(type) {
			case ast.FnDeclaration: {
				
			}

			case ast.RecordDeclaration: {
				// dummy error; not yet supported.
				printDiag(a.makeExpectedTypeAnnotation(stmt.Name))
				continue
			}

			// assuming this isn't at top-level, and this doesn't tell the type of the current block.
			case ast.ReturnStatement: {
				if stmt.Expression == nil {
					// no expression = void
					generatedTast = append(generatedTast, newStmt(tast.ReturnStatement{
						Expression: tast.Expression{
							Base: tast.AstBase(s.Base),
							Data: tast.VoidExpression{},
						},
					}))

					continue
				}

                expr, diag := a.analyzeExpression(*stmt.Expression, false); if diag != nil {
					printDiag(diag)
					continue
				}

				generatedTast = append(generatedTast, newStmt(tast.ReturnStatement{
					Expression: expr,
				}))
			}

			// this tells the type of the block.
			case ast.OutStatement: {
				if stmt.Expression == nil {
					// no expression = void
					generatedTast = append(generatedTast, newStmt(tast.OutStatement{
						Expression: tast.Expression{
							Base: tast.AstBase(s.Base),
							Data: tast.VoidExpression{},
						},
					}))

					continue
				}

                expr, diag := a.analyzeExpression(*stmt.Expression, false); if diag != nil {
					printDiag(diag)
					continue
				}

				generatedTast = append(generatedTast, newStmt(tast.OutStatement{
					Expression: expr,
				}))
			}

			case ast.VarDeclaration: {}

			case ast.WhileStatement: {}

			case ast.ForStatement: {}

			case ast.ForVarStatement: {}
			
            // TODO: remove code repetition
			case ast.BreakStatement: {
				if !a.isInsideLoop {
					printDiag(a.makeBreakContinueOutsideLoop(s.Base.Token))
				}

				generatedTast = append(generatedTast, newStmt(tast.BreakStatement{}))
			}

			case ast.ContinueStatement: {
				if !a.isInsideLoop {
					printDiag(a.makeBreakContinueOutsideLoop(s.Base.Token))
				}

				generatedTast = append(generatedTast, newStmt(tast.ContinueStatement{}))
			}

			case ast.ExprStatement: {
				expr, diag := a.analyzeExpression(stmt.Expr, false); if diag != nil {
					printDiag(diag)
					continue
				}

				generatedTast = append(generatedTast, newStmt(tast.ExprStatement{
					Expr: expr,
				}))
			}
		}
	}

	var blockType types.TypeData = types.TypeVoid{}

	if inferredType != nil {
		blockType = *inferredType
	}

	return tast.BlockExpression{
		Stmts: generatedTast,
		BlockType: types.DummyType(blockType), // dummy because it's inferred.
	}, res
}

func (a *Analyzer) analyzeExpression(expr ast.Expression, shallow bool) (tast.Expression, diagnostic.Diagnostic) {

}
