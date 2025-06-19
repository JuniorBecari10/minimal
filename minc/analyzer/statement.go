package analyzer

import (
	"fmt"
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

	a.newScope()
	var inferredType *types.TypeData = nil
	
	for _, s := range block {
		stmt, diag := a.analyzeStatement(s, &inferredType, mode); if diag != nil {
			diag.PrintDiagnostic()

			if diag.DiagnosticType() == diagnostic.TYPE_ERROR {
				res = RES_ERROR
			}
		}

		generatedTast = append(generatedTast, stmt)
	}

	var blockType types.TypeData = types.TypeVoid{}

	if inferredType != nil {
		blockType = *inferredType
	}

	a.endScope()

	return tast.BlockExpression{
		Stmts: generatedTast,
		BlockType: types.DummyType(blockType), // dummy because it's inferred.
	}, res
}

func (a *Analyzer) analyzeStatement(s ast.Statement, inferredType **types.TypeData, mode BlockAnalyzeMode) (tast.Statement, diagnostic.Diagnostic) {
	newStmt := func(data tast.StmtData) tast.Statement {
		return tast.Statement{
			Base: tast.AstBase(s.Base),
			Data: data,
		}
	}

	switch stmt := s.Data.(type) {
		case ast.FnDeclaration: {
			decl, diag := a.fnDecl(stmt); if diag != nil {
				return tast.Statement{}, diag
			}

			return newStmt(decl), nil
		}

		case ast.RecordDeclaration: {
			// dummy error; not yet supported.
			return tast.Statement{}, a.makeTypeAnnotationsNeeded(stmt.Name)
		}

		case ast.ReturnStatement: {
			stmt, diag := a.returnStmt(s, stmt, inferredType, mode); if diag != nil {
				return tast.Statement{}, diag
			}

			return newStmt(stmt), nil
		}

		case ast.OutStatement: {
			stmt, diag := a.outStmt(s, stmt, inferredType); if diag != nil {
				return tast.Statement{}, diag
			}

			return newStmt(stmt), nil
		}

		case ast.VarDeclaration: {
			decl, diag := a.varDecl(stmt); if diag != nil {
				return tast.Statement{}, diag
			}

			return newStmt(decl), nil
		}

		case ast.WhileStatement: {
			stmt, diag := a.whileStmt(stmt); if diag != nil {
				return tast.Statement{}, diag
			}

			return newStmt(stmt), nil
		}

		case ast.ForStatement: {
			stmt, diag := a.forStmt(stmt); if diag != nil {
				return tast.Statement{}, diag
			}

			return newStmt(stmt), nil
		}

		// when mode is loop, it sets its type to void, otherwise, never.
		case ast.BreakStatement: {
			stmt, diag := a.breakStmt(s, inferredType, mode); if diag != nil {
				return tast.Statement{}, diag
			}

			return newStmt(stmt), nil
		}

		case ast.ContinueStatement: {
			stmt, diag := a.continueStmt(s, inferredType, mode); if diag != nil {
				return tast.Statement{}, diag
			}

			return newStmt(stmt), nil
		}

		case ast.ExprStatement: {
			stmt, diag := a.exprStmt(stmt); if diag != nil {
				return tast.Statement{}, diag
			}

			return newStmt(stmt), nil
		}
	}

	panic(fmt.Sprintf("Unhandled statement type: %#v", s))
}

// ---

// assuming this isn't at top-level, and this doesn't tell the type of the current block.
// 'return' only tells the type of the block if this is the function block; otherwise it is never.
func (a *Analyzer) returnStmt(
	s ast.Statement,
	stmt ast.ReturnStatement,
	inferredType **types.TypeData,
	mode BlockAnalyzeMode,
) (tast.ReturnStatement, diagnostic.Diagnostic) {
	if stmt.Expression == nil {
		// no expression = void.

		// set the inferred type to 'void', if it's a function's body.
		if *inferredType == nil && mode == MODE_FUNCTION {
			var infer types.TypeData = types.TypeVoid{}
			*inferredType = &infer
		} else {
			// else, we set it to never, since returning in an inner block makes it not return anything.
			var infer types.TypeData = types.TypeNever{}
			*inferredType = &infer
		}

		return tast.ReturnStatement{
			Expression: tast.Expression{
				Base: tast.AstBase(s.Base),
				Data: tast.VoidExpression{},
			},
		}, nil
	}

	expr, diag := a.analyzeExpression(*stmt.Expression, false); if diag != nil {
		return tast.ReturnStatement{}, diag
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

	return tast.ReturnStatement{
		Expression: expr,
	}, nil
}

// this always tells the type of the block.
func (a *Analyzer) outStmt(
	s ast.Statement,
	stmt ast.OutStatement,
	inferredType **types.TypeData,
) (tast.OutStatement, diagnostic.Diagnostic) {
	if stmt.Expression == nil {
		// no expression = void

		// set the inferred type to 'void'.
		if *inferredType == nil {
			var infer types.TypeData = types.TypeVoid{}
			*inferredType = &infer
		}

		// ok
		return tast.OutStatement{
			Expression: tast.Expression{
				Base: tast.AstBase(s.Base),
				Data: tast.VoidExpression{},
			},
		}, nil
	}

	expr, diag := a.analyzeExpression(*stmt.Expression, false); if diag != nil {
		return tast.OutStatement{}, diag
	}

	// set the inferred type to be the type of the expression
	if *inferredType == nil {
		infer := expr.Data.Type()
		*inferredType = &infer
	}

	return tast.OutStatement{
		Expression: expr,
	}, nil
}

func (a *Analyzer) forStmt(stmt ast.ForStatement) (tast.ForStatement, diagnostic.Diagnostic) {
	iterable, diag := a.analyzeExpression(stmt.Iterable, false); if diag != nil {
		return tast.ForStatement{}, diag
	}

	// check if type is not iterable
	if !typeIsIterable(iterable.Data.Type()) {
		return tast.ForStatement{}, a.makeExpectedIterableType(types.Type{
			Token: iterable.Base.Token,
			Data: iterable.Data.Type(),
		})
	}

	varType := getIteratorType(iterable)

	block, res := a.analyzeBlock(stmt.Block.Stmts, MODE_LOOP); if res == RES_ERROR {
		return tast.ForStatement{}, diagnostic.HandledDiagnostic{}
	}

	return tast.ForStatement{
		Variable: stmt.Variable,
		VariableType: varType,
		Iterable: iterable,
		Block: block,
	}, nil
}

func (a *Analyzer) whileStmt(stmt ast.WhileStatement) (tast.WhileStatement, diagnostic.Diagnostic) {
	condition, diag := a.analyzeExpression(stmt.Condition, false); if diag != nil {
		return tast.WhileStatement{}, diag
	}

	// check if it's a boolean or it can coerce to it (future-proof, since currently there is no type that can coerce to a bool).
	if _, ok := tryCoercing(condition.Data.Type(), types.TypeBool{}); !ok {
		return tast.WhileStatement{}, a.makeExpectedType(types.TypeBool{}, condition.Data.Type(), condition.Base.Token)
	}

	block, res := a.analyzeBlock(stmt.Block.Stmts, MODE_LOOP); if res == RES_ERROR {
		return tast.WhileStatement{}, diagnostic.HandledDiagnostic{}
	}

	return tast.WhileStatement{
		Condition: condition,
		Block: block,
	}, nil
}

func (a *Analyzer) breakStmt(
	s ast.Statement,
	inferredType **types.TypeData,
	mode BlockAnalyzeMode,
) (tast.BreakStatement, diagnostic.Diagnostic) {
	if !a.isInsideLoop {
		return tast.BreakStatement{}, a.makeBreakContinueOutsideLoop(s.Base.Token)
	}

	if mode == MODE_LOOP {
		var infer types.TypeData = types.TypeVoid{}
		*inferredType = &infer
	} else {
		var infer types.TypeData = types.TypeNever{}
		*inferredType = &infer
	}

	return tast.BreakStatement{}, nil
}

func (a *Analyzer) continueStmt(
	s ast.Statement,
	inferredType **types.TypeData,
	mode BlockAnalyzeMode,
) (tast.ContinueStatement, diagnostic.Diagnostic) {
	_, diag := a.breakStmt(s, inferredType, mode)
	return tast.ContinueStatement{}, diag
}

func (a *Analyzer) exprStmt(stmt ast.ExprStatement) (tast.ExprStatement, diagnostic.Diagnostic) {
	expr, diag := a.analyzeExpression(stmt.Expr, false); if diag != nil {
		return tast.ExprStatement{}, diag
	}

	return tast.ExprStatement{
		Expr: expr,
	}, nil
}
