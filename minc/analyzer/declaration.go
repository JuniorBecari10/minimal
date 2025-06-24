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
func (a *Analyzer) fnDecl(decl ast.FnDeclaration) (tast.FnDeclaration, diagnostic.Diagnostic) {
	returnType := types.DummyType(types.TypeVoid{})

	if decl.Return != nil {
		returnType = *decl.Return
	}

	paramTypes := []types.Type{}

	for _, param := range decl.Parameters {
		if param.Type == nil {
			return tast.FnDeclaration{}, a.makeTypeAnnotationsNeeded(param.Name)
		}

		paramTypes = append(paramTypes, *param.Type)
	}

	// return type may be unknown. check the body and see if the type can be coerced to it.
	body, res := a.analyzeBlockAlone(decl.Body, MODE_FUNCTION); if res == RES_ERROR {
		return tast.FnDeclaration{}, diagnostic.HandledDiagnostic{}
	}

	// if the return type is unknown, it will coerce.
	coerced, ok := coerceExpr(tast.Expression{
		Base: tast.AstBase{},
		Data: body,
	}, returnType.Data); if !ok {
		return tast.FnDeclaration{}, a.makeExpectedType(returnType.Data, body.Type(), returnType.Token)
	}

	returnType.Data = coerced.Data.Type()

	a.addVariable(decl.Name, types.Type{
		Token: decl.Name,
		Data: types.TypeFunction{
			Parameters: paramTypes,
			Return: returnType,
		},
	}, true)

	return tast.FnDeclaration{
		Name: decl.Name,
		Parameters: decl.Parameters,
		Return: returnType,
		Body: body,
	}, nil
}

func (a *Analyzer) varDecl(decl ast.VarDeclaration) (tast.VarDeclaration, diagnostic.Diagnostic) {
	expr, diag := a.analyzeExpression(decl.Init, false, nil); if diag != nil {
		return tast.VarDeclaration{}, diag
	}

	if decl.Type == nil {
		if !typeIsConcrete(expr.Data.Type()) {
			return tast.VarDeclaration{}, a.makeTypeAnnotationsNeeded(decl.Name)
		}

		varType := types.DummyType(expr.Data.Type())
		a.addVariable(decl.Name, varType, decl.Immutable)

		return tast.VarDeclaration{
			Name: decl.Name,
			Init: expr,
			Type: varType,
		}, nil
	} else {
		if ok := canCoerce(expr.Data.Type(), decl.Type.Data); !ok {
			return tast.VarDeclaration{}, a.makeExpectedType(decl.Type.Data, expr.Data.Type(), decl.Type.Token)
		}

		var ok bool
		expr, ok = coerceExpr(expr, decl.Type.Data); if !ok {
			return tast.VarDeclaration{}, a.makeExpectedType(decl.Type.Data, expr.Data.Type(), decl.Type.Token)
		}

		a.addVariable(decl.Name, *decl.Type, decl.Immutable)

		return tast.VarDeclaration{
			Name: decl.Name,
			Init: expr,
			Type: *decl.Type,
		}, nil
	}
}
