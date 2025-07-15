package analyzer

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
	"minc/types"
	"minlib/token"
)

// this returns a block, with its type inferred by its statements.
// the first statement that can set the type of the block directly inside it will do it.
// this function is meant to analyze blocks stored alone in the AST nodes, like in 'while', 'for' and 'if'.
func (a *Analyzer) analyzeBlockAlone(block ast.BlockExpression, mode BlockAnalyzeMode) (tast.BlockExpression, AnalyzerResult) {
	res := RES_OK

	blockRes, diag := a.analyzeBlock(block, mode, false); if diag != nil {
		a.diagnostics = append(a.diagnostics, diag)
		res = RES_ERROR
	}

	return blockRes, res
}

// synchronization point. this function prints the diagnostics and doesn't bubble them up.
// this analyzes a block expression and tries to infer its type.
func (a *Analyzer) analyzeBlock(block ast.BlockExpression, mode BlockAnalyzeMode, shallow bool) (tast.BlockExpression, diagnostic.Diagnostic) {
	if shallow {
		// don't enter the block; return an empty one with unknown type.
		return tast.BlockExpression{
			Stmts: []tast.Statement{},
			BlockType: types.TypeUnknown{},
			Token: block.Token,
		}, nil
	// Special case where the block has one statement and it is an expression statement.
	} else if len(block.Stmts) == 1 {
		// Confirmation of length needs to be before checking the type of the statement.
		if stmt, ok := block.Stmts[0].Data.(ast.ExprStatement); ok {
			var diag diagnostic.Diagnostic

            // if it has a semicolon here, report it as a warning, since it is redundant.
			if stmt.Semicolon != nil && block.ShowSemicolonWarning {
				diag = a.makeWarnRedundantSemicolon(*stmt.Semicolon)
			}
			
			// make the expected type to the function's return type if this is a function body.
			var expectedType *types.TypeData = nil

			if mode == MODE_FUNCTION {
				expectedType = a.expectedReturnType
			}

			// yes, if diag == nil.
			exprInside, exprDiag := a.analyzeExpression(stmt.Expr, shallow, expectedType); if diag == nil {
				if mode == MODE_FUNCTION {
					// TODO: if it's an expected type diagnostic, change it to be an expected return one.
					diag = exprDiag
				} else {
					diag = exprDiag
				}
			}

			if diag != nil {
				return tast.BlockExpression{}, diag
			}

			return tast.BlockExpression{
				Stmts: []tast.Statement{
					{
						Base: exprInside.Base,
						Data: tast.OutStatement{
						    Expression: exprInside,
						},
					},
				},
				BlockType: exprInside.Data.Type(),
				Token: block.Token,
			}, diag
		}
	}

	// both else blocks
	blockRes, res := a.analyzeStatements(block.Stmts, block.Token, mode); if res == RES_ERROR {
		return tast.BlockExpression{}, diagnostic.HandledDiagnostic{}
	}

	// check if there is more than one statement
	if len(block.Stmts) > 1 {
		if stmt, ok := block.Stmts[len(block.Stmts) - 1].Data.(ast.ExprStatement); ok {
			// if it doesn't have a semicolon, report an error.
			if stmt.Semicolon == nil {
				return blockRes, a.makeExpectedSemicolon(stmt.Expr.Base.Token)
			}
		}
	}

	return blockRes, nil
}

func (a *Analyzer) analyzeStatements(stmts ast.Ast, tok token.Token, mode BlockAnalyzeMode) (tast.BlockExpression, AnalyzerResult) {
	a.newScope()
    defer a.endScope()

	var inferredType *types.TypeData = nil
	
	// res is carried forward
	generatedTast, res := a.statements(stmts, &inferredType, mode)
	var blockType types.TypeData = types.TypeVoid{}

	if inferredType != nil {
		// block types are not always concrete, so we don't need to check.
		blockType = *inferredType
	}

	return tast.BlockExpression{
		Stmts: generatedTast,
		BlockType: blockType,
		Token: tok,
	}, res
}

func (a *Analyzer) analyzeTopLevelStatements() (tast.Tast, AnalyzerResult) {
	return a.statements(a.ast, nil, MODE_NORMAL)
}

func (a *Analyzer) statements(ast ast.Ast, inferredType **types.TypeData, mode BlockAnalyzeMode) (tast.Tast, AnalyzerResult) {
	generatedTast := make(tast.Tast, 0, len(ast))
	res := RES_OK

	for _, s := range ast {
		stmt, diag := a.analyzeStatement(s, inferredType, mode); if diag != nil {
			a.diagnostics = append(a.diagnostics, diag)

			if diag.DiagnosticType() == diagnostic.TYPE_ERROR {
				res = RES_ERROR
			}
		}

		generatedTast = append(generatedTast, stmt)
	}

	return generatedTast, res
}
