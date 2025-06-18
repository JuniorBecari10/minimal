package analyzer

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
	"minc/types"
)

// the return type can be unknown, but inferrable from its block.
// the parameters need to be explicitly typed, since we can't infer their types from context,
// because this is a statement, and not an expression, like a lambda.
// can return a warning.
func (a *Analyzer) fnDecl(
	decl ast.FnDeclaration,
	generatedTast *tast.Tast,
	newStmt func(tast.StmtData) tast.Statement,
) diagnostic.Diagnostic {
	returnType := types.DummyType(types.TypeVoid{})

	if decl.Return != nil {
		returnType = *decl.Return
	}

	paramTypes := []types.Type{}

	for _, param := range decl.Parameters {
		if param.Type == nil {
			return a.makeExpectedTypeAnnotation(param.Name)
		}

		paramTypes = append(paramTypes, *param.Type)
	}

	// return type may be unknown. check the body and see if the type can be coerced to it.
	body, res := a.analyzeBlock(decl.Body.Stmts, MODE_FUNCTION); if res == RES_ERROR {
		return diagnostic.HandledDiagnostic{}
	}

	// check for unreachable code.
	if _, ok := body.Type().(types.TypeNever); ok {
		return a.makeWarnUnreachable(decl.Name)
	}

	// if the return type is unknown, it will coerce.
	mergedType, ok := tryCoercing(body.Type(), returnType.Data); if !ok {
		return a.makeExpectedType(returnType.Data, body.Type(), returnType.Token)
	}

	// merge the types to one that supports both.
	returnType.Data = mergedType
	
	*generatedTast = append(*generatedTast, newStmt(tast.FnDeclaration{
		Name: decl.Name,
		Parameters: decl.Parameters,
		Return: returnType,
		Body: body,
	}))
	return nil
}

// can return a warning.
func (a *Analyzer) varDecl(
	decl ast.VarDeclaration,
	generatedTast *tast.Tast,
	newStmt func(tast.StmtData) tast.Statement,
) diagnostic.Diagnostic {
	expr, diag := a.analyzeExpression(decl.Init, false); if diag != nil {
		return diag
	}

	// check for unreachable code.
	if _, ok := expr.Data.Type().(types.TypeNever); ok {
		return a.makeWarnUnreachable(decl.Name)
	}
	
	if decl.Type == nil {
		if !typeIsConcrete(expr.Data.Type()) {
			return a.makeExpectedTypeAnnotation(decl.Name)
		}

		*generatedTast = append(*generatedTast, newStmt(tast.VarDeclaration{
			Name: decl.Name,
			Init: expr,
			Type: types.DummyType(expr.Data.Type()),
		}))
	} else {
		if _, ok := tryCoercing(expr.Data.Type(), decl.Type.Data); !ok {
			return a.makeExpectedType(decl.Type.Data, expr.Data.Type(), decl.Type.Token)
		}

		*generatedTast = append(*generatedTast, newStmt(tast.VarDeclaration{
			Name: decl.Name,
			Init: expr,
			Type: *decl.Type,
		}))
	}

	return nil
}
