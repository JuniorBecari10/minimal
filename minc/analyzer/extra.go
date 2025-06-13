package analyzer

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
	"minc/types"
)

type BlockAnalyzeMode int

const (
	MODE_NORMAL BlockAnalyzeMode = iota
	MODE_FUNCTION
	MODE_LOOP
)

// synchronization point. this function prints the diagnostics and doesn't bubble them up
// this returns a block, with its type inferred by its statements.
// the first statement that can set the type of the block directly inside it will do it.
func (a *Analyzer) analyzeBlock(block ast.Ast, mode BlockAnalyzeMode) (tast.BlockExpression, AnalyzerResult) {
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
			// 'return' only tells the type of the block if this is the function block; otherwise it is never.
			case ast.ReturnStatement: {
				if stmt.Expression == nil {
					// no expression = void
					generatedTast = append(generatedTast, newStmt(tast.ReturnStatement{
						Expression: tast.Expression{
							Base: tast.AstBase(s.Base),
							Data: tast.VoidExpression{},
						},
					}))

					// set the inferred type to 'void', if it's a function's body.
					if inferredType == nil && mode == MODE_FUNCTION {
						var infer types.TypeData = types.TypeVoid{}
						inferredType = &infer
					} else {
						// else, we set it to never, since returning in an inner block makes it not return anything.
						var infer types.TypeData = types.TypeNever{}
						inferredType = &infer
					}

					continue
				}

                expr, diag := a.analyzeExpression(*stmt.Expression, false); if diag != nil {
					printDiag(diag)
					continue
				}
				
				// set the inferred type to be the type of the expression, if it's a function's body.
				if inferredType == nil {
					if mode == MODE_FUNCTION {
						infer := expr.Data.Type()
						inferredType = &infer
					} else {
						// else, we set it to never, since returning in an inner block makes it not return anything.
						var infer types.TypeData = types.TypeNever{}
						inferredType = &infer
					}
				}

				generatedTast = append(generatedTast, newStmt(tast.ReturnStatement{
					Expression: expr,
				}))
			}

			// this always tells the type of the block.
			case ast.OutStatement: {
				if stmt.Expression == nil {
					// no expression = void
					generatedTast = append(generatedTast, newStmt(tast.OutStatement{
						Expression: tast.Expression{
							Base: tast.AstBase(s.Base),
							Data: tast.VoidExpression{},
						},
					}))

					// set the inferred type to 'void'.
					if inferredType == nil {
						var infer types.TypeData = types.TypeVoid{}
						inferredType = &infer
					}

					continue
				}

                expr, diag := a.analyzeExpression(*stmt.Expression, false); if diag != nil {
					printDiag(diag)
					continue
				}

				// set the inferred type to be the type of the expression
				if inferredType == nil {
					infer := expr.Data.Type()
					inferredType = &infer
				}

				generatedTast = append(generatedTast, newStmt(tast.OutStatement{
					Expression: expr,
				}))
			}

			case ast.VarDeclaration: {

			}

			case ast.WhileStatement: {
				condition, diag := a.analyzeExpression(stmt.Condition, false); if diag != nil {
					printDiag(diag)
					continue
				}

				// check if it's a boolean or it can coerce to it.
				if !typeCanCoerceTo(condition.Data.Type(), types.TypeBool{}) {
					printDiag(a.makeExpectedType(types.TypeBool{}, condition.Data.Type(), condition.Base.Token))
					continue
				}
			}

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
