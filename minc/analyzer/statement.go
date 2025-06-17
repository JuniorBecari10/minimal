package analyzer

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
	"minc/types"
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
				diag := a.fnDecl(stmt, &generatedTast, newStmt); if diag != nil {
					printDiag(diag)
					continue
				}
			}

			case ast.RecordDeclaration: {
				// dummy error; not yet supported.
				printDiag(a.makeExpectedTypeAnnotation(stmt.Name))
				continue
			}

			case ast.ReturnStatement: a.returnStmt(s, stmt, &generatedTast, newStmt, printDiag, &inferredType, mode)
			case ast.OutStatement: a.outStmt(s, stmt, &generatedTast, newStmt, printDiag, &inferredType)

			case ast.VarDeclaration: {
				diag := a.varDecl(stmt, &generatedTast, newStmt); if diag != nil {
					printDiag(diag)
					continue
				}
			}

			case ast.WhileStatement: {
				diag := a.whileStmt(stmt, &generatedTast, newStmt); if diag != nil {
					printDiag(diag)
					continue
				}
			}

			case ast.ForStatement: {
				diag := a.forStmt(stmt, &generatedTast, newStmt); if diag != nil {
					printDiag(diag)
					continue
				}
			}

			// do this later
			case ast.ForVarStatement: {
				// dummy error; not yet supported.
				printDiag(a.makeExpectedTypeAnnotation(stmt.Condition.Base.Token))
				continue
			}
			
            // TODO: remove code repetition
			// when mode is loop, it set the type to void, otherwise, never.
			case ast.BreakStatement: {
				diag := a.breakStmt(s, &generatedTast, newStmt, &inferredType, mode); if diag != nil {
					printDiag(diag)
					continue
				}
			}

			// when mode is loop, it set the type to void, otherwise, never.
			case ast.ContinueStatement: {
				diag := a.continueStmt(s, &generatedTast, newStmt, &inferredType, mode); if diag != nil {
					printDiag(diag)
					continue
				}
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

// ---

// assuming this isn't at top-level, and this doesn't tell the type of the current block.
// 'return' only tells the type of the block if this is the function block; otherwise it is never.
func (a *Analyzer) returnStmt(
	s ast.Statement,
	stmt ast.ReturnStatement,
	generatedTast *tast.Tast,
	newStmt func(tast.StmtData) tast.Statement,
	printDiag func(diagnostic.Diagnostic),
	inferredType **types.TypeData,
	mode BlockAnalyzeMode,
) {
	if stmt.Expression == nil {
		// no expression = void
		*generatedTast = append(*generatedTast, newStmt(tast.ReturnStatement{
			Expression: tast.Expression{
				Base: tast.AstBase(s.Base),
				Data: tast.VoidExpression{},
			},
		}))

		// set the inferred type to 'void', if it's a function's body.
		if *inferredType == nil && mode == MODE_FUNCTION {
			var infer types.TypeData = types.TypeVoid{}
			*inferredType = &infer
		} else {
			// else, we set it to never, since returning in an inner block makes it not return anything.
			var infer types.TypeData = types.TypeNever{}
			*inferredType = &infer
		}

		return
	}

	expr, diag := a.analyzeExpression(*stmt.Expression, false); if diag != nil {
		printDiag(diag)
		return
	}
	
	// set the inferred type to be the type of the expression, if it's a function's body.
	if inferredType == nil {
		if mode == MODE_FUNCTION {
			infer := expr.Data.Type()
			*inferredType = &infer
		} else {
			// else, we set it to never, since returning in an inner block makes it not return anything.
			var infer types.TypeData = types.TypeNever{}
			*inferredType = &infer
		}
	}

	*generatedTast = append(*generatedTast, newStmt(tast.ReturnStatement{
		Expression: expr,
	}))
}

// this always tells the type of the block.
func (a *Analyzer) outStmt(
	s ast.Statement,
	stmt ast.OutStatement,
	generatedTast *tast.Tast,
	newStmt func(tast.StmtData) tast.Statement,
	printDiag func(diagnostic.Diagnostic),
	inferredType **types.TypeData,
) {
	if stmt.Expression == nil {
		// no expression = void
		*generatedTast = append(*generatedTast, newStmt(tast.OutStatement{
			Expression: tast.Expression{
				Base: tast.AstBase(s.Base),
				Data: tast.VoidExpression{},
			},
		}))

		// set the inferred type to 'void'.
		if *inferredType == nil {
			var infer types.TypeData = types.TypeVoid{}
			*inferredType = &infer
		}

		return
	}

	expr, diag := a.analyzeExpression(*stmt.Expression, false); if diag != nil {
		printDiag(diag)
		return
	}

	// set the inferred type to be the type of the expression
	if *inferredType == nil {
		infer := expr.Data.Type()
		*inferredType = &infer
	}

	*generatedTast = append(*generatedTast, newStmt(tast.OutStatement{
		Expression: expr,
	}))
}

func (a *Analyzer) forStmt(
	stmt ast.ForStatement,
	generatedTast *tast.Tast,
	newStmt func(tast.StmtData) tast.Statement,
) diagnostic.Diagnostic {
	iterable, diag := a.analyzeExpression(stmt.Iterable, false); if diag != nil {
		return diag
	}

	if !typeIsIterable(iterable.Data.Type()) {
		return a.makeExpectedIterableType(types.Type{
			Token: iterable.Base.Token,
			Data: iterable.Data.Type(),
		})
	}

	varType := getIteratorType(iterable)

	block, res := a.analyzeBlock(stmt.Block.Stmts, MODE_LOOP); if res == RES_ERROR {
		return diagnostic.HandledDiagnostic{}
	}

	*generatedTast = append(*generatedTast, newStmt(tast.ForStatement{
		Variable: stmt.Variable,
		VariableType: varType,
		Iterable: iterable,
		Block: block,
	}))

	return nil
}

func (a *Analyzer) whileStmt(
	stmt ast.WhileStatement,
	generatedTast *tast.Tast,
	newStmt func(tast.StmtData) tast.Statement,
) diagnostic.Diagnostic {
	condition, diag := a.analyzeExpression(stmt.Condition, false); if diag != nil {
		return diag
	}

	// check if it's a boolean or it can coerce to it.
	if _, ok := tryCoercing(condition.Data.Type(), types.TypeBool{}); !ok {
		return a.makeExpectedType(types.TypeBool{}, condition.Data.Type(), condition.Base.Token)
	}

	block, res := a.analyzeBlock(stmt.Block.Stmts, MODE_LOOP); if res == RES_ERROR {
		return diagnostic.HandledDiagnostic{}
	}

	*generatedTast = append(*generatedTast, newStmt(tast.WhileStatement{
		Condition: condition,
		Block: block,
	}))

	return nil
}

func (a *Analyzer) breakStmt(
	s ast.Statement,
	generatedTast *tast.Tast,
	newStmt func(tast.StmtData) tast.Statement,
	inferredType **types.TypeData,
	mode BlockAnalyzeMode,
) diagnostic.Diagnostic {
	if !a.isInsideLoop {
		return a.makeBreakContinueOutsideLoop(s.Base.Token)
	}

	if mode == MODE_LOOP {
		var infer types.TypeData = types.TypeVoid{}
		*inferredType = &infer
	} else {
		var infer types.TypeData = types.TypeNever{}
		*inferredType = &infer
	}

	*generatedTast = append(*generatedTast, newStmt(tast.BreakStatement{}))
	return nil
}

func (a *Analyzer) continueStmt(
	s ast.Statement,
	generatedTast *tast.Tast,
	newStmt func(tast.StmtData) tast.Statement,
	inferredType **types.TypeData,
	mode BlockAnalyzeMode,
) diagnostic.Diagnostic {
	if !a.isInsideLoop {
		return a.makeBreakContinueOutsideLoop(s.Base.Token)
	}

	if mode == MODE_LOOP {
		var infer types.TypeData = types.TypeVoid{}
		*inferredType = &infer
	} else {
		var infer types.TypeData = types.TypeNever{}
		*inferredType = &infer
	}

	*generatedTast = append(*generatedTast, newStmt(tast.ContinueStatement{}))
	return nil
}
