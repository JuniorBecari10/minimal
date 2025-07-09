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
	inside := a.isInsideLoop
	a.isInsideLoop = false // functions cancel this

	defer func() { a.isInsideLoop = inside }()

	returnType := types.DummyType(types.TypeVoid{})
	retAnnotated := false

	if decl.Return != nil {
		returnType = *decl.Return
		retAnnotated = true
	}

	oldExpectedReturn := a.expectedReturnType
	a.expectedReturnType = &returnType.Data

	defer func() { a.expectedReturnType = oldExpectedReturn }()

	a.newScope()
	defer a.endScope()

	paramTypes := []types.Type{}

	// function declarations must have annotated parameters with concrete types.
	// the compiler won't try to infer them.
	for _, param := range decl.Parameters {
		if param.Type == nil {
			return tast.FnDeclaration{}, a.makeTypeAnnotationsNeeded(param.Name)
		} else if !typeIsConcrete(param.Type.Data) {
			return tast.FnDeclaration{}, a.makeExpectedConcreteType(*param.Type)
		}

		paramTypes = append(paramTypes, *param.Type)
		a.addVariable(param.Name, *param.Type, true)
	}

	// return type may be unknown. check the body and see if the type can be coerced to it.
	body, res := a.analyzeBlockAlone(decl.Body, MODE_FUNCTION); if res == RES_ERROR {
		return tast.FnDeclaration{}, diagnostic.HandledDiagnostic{}
	}

	bodyCoerced, ok := coerceExpr(body.IntoExpr(), returnType.Data); if !ok {
		if retAnnotated {
			return tast.FnDeclaration{}, a.makeExpectedReturnType(returnType.Data, body.BlockType, returnType.Token, retAnnotated)
		} else {
			return tast.FnDeclaration{}, a.makeExpectedReturnType(returnType.Data, body.BlockType, decl.Name, retAnnotated)
		}
	}

	// change return data to the block's if it's unknown (that means, open for inference)
	if _, ok := returnType.Data.(types.TypeUnknown); ok {
		returnType.Data = body.BlockType
	}
	
	if !typeIsConcrete(returnType.Data) {
		return tast.FnDeclaration{}, a.makeExpectedConcreteType(returnType)
	}

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
		// if the type is not annotated, require the type of the expression to be concrete without coercions.
		
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
		// if it is, require the type to be concrete and try to coerce the expression to it.

		if !typeIsConcrete(decl.Type.Data) {
			return tast.VarDeclaration{}, a.makeExpectedConcreteType(*decl.Type)
		}

		coercedExpr, ok := coerceExpr(expr, decl.Type.Data); if !ok {
			return tast.VarDeclaration{}, a.makeExpectedType(decl.Type.Data, coercedExpr.Data.Type(), decl.Type.Token)
		}

		a.addVariable(decl.Name, *decl.Type, decl.Immutable)

		return tast.VarDeclaration{
			Name: decl.Name,
			Init: coercedExpr,
			Type: *decl.Type,
		}, nil
	}
}
