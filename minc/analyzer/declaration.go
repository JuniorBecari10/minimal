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

	// if the return type is unknoen, it will coerce.
	if !typeCanCoerceTo(body.Type(), returnType.Data) {
		return a.makeExpectedType(returnType.Data, body.Type(), returnType.Token)
	}

	// merge the types to one that supports both.
	returnType.Data = mergeTypes(returnType.Data, body.Type())
	
	*generatedTast = append(*generatedTast, newStmt(tast.FnDeclaration{
		Name: decl.Name,
		Parameters: decl.Parameters,
		Return: returnType,
		Body: body,
	}))
	return nil
}

func (a *Analyzer) varDecl(
	decl ast.VarDeclaration,
	generatedTast *tast.Tast,
	newStmt func(tast.StmtData) tast.Statement,
) diagnostic.Diagnostic {
	expr, diag := a.analyzeExpression(decl.Init, false); if diag != nil {
		return diag
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
        if !typeCanCoerceTo(expr.Data.Type(), decl.Type.Data) {
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
