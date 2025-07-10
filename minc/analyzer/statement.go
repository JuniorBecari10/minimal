package analyzer

import (
	"fmt"
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
	"minc/types"
)

func (a *Analyzer) analyzeStatement(s ast.Statement, inferredType **types.TypeData, mode BlockAnalyzeMode) (tast.Statement, diagnostic.Diagnostic) {
	newStmt := func(data tast.StmtData) tast.Statement {
		return tast.Statement{
			Base: tast.AstBase(s.Base),
			Data: data,
		}
	}

	switch stmt := s.Data.(type) {
		case ast.FnDeclaration: {
			decl, diag := a.fnDecl(stmt)
			return newStmt(decl), diag
		}

		case ast.RecordDeclaration: {
			// dummy error; not yet supported.
			return tast.Statement{}, a.makeTypeAnnotationsNeeded(stmt.Name)
		}

		case ast.ReturnStatement: {
			stmt, diag := a.returnStmt(s, stmt, inferredType, mode)
			return newStmt(stmt), diag
		}

		case ast.OutStatement: {
			stmt, diag := a.outStmt(s, stmt, inferredType, mode)
			return newStmt(stmt), diag
		}

		case ast.VarDeclaration: {
			decl, diag := a.varDecl(stmt)
			return newStmt(decl), diag
		}

		case ast.WhileStatement: {
			stmt, diag := a.whileStmt(stmt)
			return newStmt(stmt), diag
		}

		case ast.ForStatement: {
			stmt, diag := a.forStmt(stmt)
			return newStmt(stmt), diag
		}

		case ast.LoopStatement: {
			stmt, diag := a.loopStmt(stmt)
			return newStmt(stmt), diag
		}

		case ast.BreakStatement: {
			stmt, diag := a.breakStmt(s, inferredType, mode)
			return newStmt(stmt), diag
		}

		case ast.ContinueStatement: {
			stmt, diag := a.continueStmt(s, inferredType, mode)
			return newStmt(stmt), diag
		}

		case ast.ExprStatement: {
			stmt, diag := a.exprStmt(stmt)
			return newStmt(stmt), diag
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
		if *inferredType == nil {
			if mode == MODE_FUNCTION {
				var infer types.TypeData = types.TypeVoid{}
				*inferredType = &infer
			} else {
				// else, we set it to never, since returning in an inner block makes it not return anything.
				var infer types.TypeData = types.TypeNever{}
				*inferredType = &infer
			}
		}

		return tast.ReturnStatement{
			Expression: tast.Expression{
				Base: tast.AstBase(s.Base),
				Data: tast.VoidExpression{},
			},
		}, nil
	}

	// TODO: use the function's expected return type to help inferring this type.
	// also, add an expected type parameter to analyzeBlock, analyzeStatement and some statements, such as return and out.
	expr, diag := a.analyzeExpression(*stmt.Expression, false, a.expectedReturnType); if diag != nil {
		return tast.ReturnStatement{}, diag
	}
	
	// set the inferred type to be the type of the expression, if it's a function's body.
	if *inferredType == nil {
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
	mode BlockAnalyzeMode,
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

	// TODO: use the block's expected return type to help inferring this type.
	// if it's a function, it will use the return type.

	// make the expected type to the function's return type if this is a function body.
	var expectedType *types.TypeData = nil

	if mode == MODE_FUNCTION {
		expectedType = a.expectedReturnType
	}
	fmt.Println(expectedType)

	expr, diag := a.analyzeExpression(*stmt.Expression, false, expectedType); if diag != nil {
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
	inside := a.isInsideLoop
	a.isInsideLoop = true

	defer func() { a.isInsideLoop = inside }()

	iterable, diag := a.analyzeExpression(stmt.Iterable, false, nil); if diag != nil {
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

	block, res := a.analyzeBlockAlone(stmt.Block, MODE_LOOP); if res == RES_ERROR {
		return tast.ForStatement{}, diagnostic.HandledDiagnostic{}
	}

	if !canCoerce(block.BlockType, types.TypeVoid{}) {
		return tast.ForStatement{}, a.makeExpectedType(types.TypeVoid{}, block.BlockType, block.Token)
	}

	return tast.ForStatement{
		Variable: stmt.Variable,
		VariableType: varType,
		Iterable: iterable,
		Block: block,
	}, nil
}

func (a *Analyzer) loopStmt(stmt ast.LoopStatement) (tast.LoopStatement, diagnostic.Diagnostic) {
	inside := a.isInsideLoop
	a.isInsideLoop = true

	defer func() { a.isInsideLoop = inside }()

	block, res := a.analyzeBlockAlone(stmt.Block, MODE_LOOP); if res == RES_ERROR {
		return tast.LoopStatement{}, diagnostic.HandledDiagnostic{}
	}

	if !canCoerce(block.BlockType, types.TypeVoid{}) {
		return tast.LoopStatement{}, a.makeExpectedType(types.TypeVoid{}, block.BlockType, block.Token)
	}

	return tast.LoopStatement{
		Block: block,
	}, nil
}

func (a *Analyzer) whileStmt(stmt ast.WhileStatement) (tast.WhileStatement, diagnostic.Diagnostic) {
	inside := a.isInsideLoop
	a.isInsideLoop = true

	defer func() { a.isInsideLoop = inside }()

	var typeBool types.TypeData = types.TypeBool{}
	condition, diag := a.analyzeExpression(stmt.Condition, false, &typeBool); if diag != nil {
		return tast.WhileStatement{}, diag
	}

	block, res := a.analyzeBlockAlone(stmt.Block, MODE_LOOP); if res == RES_ERROR {
		return tast.WhileStatement{}, diagnostic.HandledDiagnostic{}
	}

	if !canCoerce(block.BlockType, types.TypeVoid{}) {
		return tast.WhileStatement{}, a.makeExpectedType(types.TypeVoid{}, block.BlockType, block.Token)
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
	expr, diag := a.analyzeExpression(stmt.Expr, false, nil); if diag != nil {
		return tast.ExprStatement{}, diag
	}

	var retDiag diagnostic.Diagnostic = nil

	if typeNeedsUnusedWarning(expr.Data.Type()) {
		retDiag = a.makeWarnUnusedValueExprStmt(expr.Data.Type(), expr.Base.Token)
	}

	return tast.ExprStatement{
		Expr: expr,
	}, retDiag
}
