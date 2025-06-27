package analyzer

import (
	"fmt"
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
	"minc/types"
	"minlib/token"
)

// expectedType is optional
func (a *Analyzer) analyzeExpression(e ast.Expression, shallow bool, expectedType *types.TypeData) (tast.Expression, diagnostic.Diagnostic) {
	newExpr := func(data tast.ExprData) tast.Expression {
		return tast.Expression{
			Base: tast.AstBase(e.Base),
			Data: data,
		}
	}

	switch expr := e.Data.(type) {
		case ast.IntExpression:
			return newExpr(a.analyzeIntExpr(expr)), nil

		case ast.FloatExpression:
			return newExpr(a.analyzeFloatExpr(expr)), nil

		case ast.StringExpression:
			return newExpr(a.analyzeStringExpr(expr)), nil

		case ast.CharExpression:
			return newExpr(a.analyzeCharExpr(expr)), nil

		case ast.BoolExpression:
			return newExpr(a.analyzeBoolExpr(expr)), nil
		
		case ast.RangeExpression: {
			expr, diag := a.analyzeRangeExpr(expr, shallow)
			return newExpr(expr), diag
		}
		case ast.AsExpression: {
			expr, diag := a.analyzeAsExpr(expr, shallow)
			return newExpr(expr), diag
		}

		case ast.NilExpression:
			return newExpr(tast.NilExpression{}), nil
		
		case ast.VoidExpression: {
			expr, diag := a.analyzeVoidExpr(expr, shallow, expectedType)
			return newExpr(expr), diag
		}

		case ast.UnaryExpression: {
			expr, diag := a.analyzeUnaryExpr(expr, shallow, expectedType)
			return newExpr(expr), diag
		}

		case ast.LogicalExpression: {
			expr, diag := a.analyzeLogicalExpr(expr, shallow)
			return newExpr(expr), diag
		}

		case ast.BinaryExpression: {
			expr, diag := a.analyzeBinaryExpr(expr, shallow)
			return newExpr(expr), diag
		}

		case ast.CallExpression: {
			expr, diag := a.analyzeCallExpr(expr, shallow, expectedType)
			return newExpr(expr), diag
		}

		case ast.GroupExpression: {
			expr, diag := a.analyzeGroupExpr(expr, shallow, expectedType)
			return newExpr(expr), diag
		}

		case ast.IdentifierExpression: {
			expr, diag := a.analyzeIdentifierExpr(expr)
			return newExpr(expr), diag
		}

		case ast.SelfExpression: {
			// dummy error. not yet supported.
			return newExpr(tast.SelfExpression{}), a.makeExpectedSemicolon(token.StartToken())
		}

		case ast.IdentifierAssignmentExpression: {
			expr, diag := a.analyzeIdentifierAssignExpr(expr, shallow, expectedType)
			return newExpr(expr), diag
		}

		case ast.FnExpression: {/*
			expr, diag := a.analyzeFnExpr(expr, shallow, expectedType)
			return newExpr(expr), diag
		*/}

		case ast.BlockExpression: {
			expr, diag := a.analyzeBlockExpr(expr, shallow)
			return newExpr(expr), diag
		}

		// because of 'else' blocks being optional.
		case *ast.BlockExpression: {
			expr, diag := a.analyzeBlockExpr(*expr, shallow)
			return newExpr(expr), diag
		}
		
		case ast.IfExpression: {
			expr, diag := a.analyzeIfExpr(expr, shallow, expectedType)
			return newExpr(expr), diag
		}

		case ast.GetPropertyExpression: {
			// dummy error. not yet supported.
			return newExpr(tast.GetPropertyExpression{}), a.makeExpectedSemicolon(token.StartToken())
		}

		case ast.SetPropertyExpression: {
			// dummy error. not yet supported.
			return newExpr(tast.SetPropertyExpression{}), a.makeExpectedSemicolon(token.StartToken())
		}
	}

	panic(fmt.Sprintf("Internal: Invalid expression: %#v", e))
}

func (a *Analyzer) analyzeIntExpr(expr ast.IntExpression) tast.IntExpression {
	return tast.IntExpression{
		Literal: expr.Literal,
	}
}

func (a *Analyzer) analyzeFloatExpr(expr ast.FloatExpression) tast.FloatExpression {
	return tast.FloatExpression{
		Literal: expr.Literal,
	}
}

func (a *Analyzer) analyzeStringExpr(expr ast.StringExpression) tast.StringExpression {
	return tast.StringExpression{
		Literal: expr.Literal,
	}
}

func (a *Analyzer) analyzeCharExpr(expr ast.CharExpression) tast.CharExpression {
	return tast.CharExpression{
		Literal: expr.Literal,
	}
}

func (a *Analyzer) analyzeBoolExpr(expr ast.BoolExpression) tast.BoolExpression {
	return tast.BoolExpression{
		Literal: expr.Literal,
	}
}

// TODO: finish
func (a *Analyzer) analyzeRangeExpr(expr ast.RangeExpression, shallow bool) (tast.RangeExpression, diagnostic.Diagnostic) {
	return tast.RangeExpression{}, nil
}

// TODO: finish
func (a *Analyzer) analyzeAsExpr(expr ast.AsExpression, shallow bool) (tast.AsExpression, diagnostic.Diagnostic) {
    return tast.AsExpression{}, nil
}

func (a *Analyzer) analyzeVoidExpr(expr ast.VoidExpression, shallow bool, expectedType *types.TypeData) (tast.VoidExpression, diagnostic.Diagnostic) {
	if expr.Expr == nil {
		return tast.VoidExpression{
			Expr: nil,
		}, nil
	} else {
		expr, diag := a.analyzeExpression(*expr.Expr, shallow, expectedType); if diag != nil {
			return tast.VoidExpression{}, diag
		}

		return tast.VoidExpression{
			Expr: &expr,
		}, nil
	}
}

func (a *Analyzer) analyzeUnaryExpr(expr ast.UnaryExpression, shallow bool, expectedType *types.TypeData) (tast.UnaryExpression, diagnostic.Diagnostic) {
	// not and -. assuming the parser checked this.

	operand, diag := a.analyzeExpression(expr.Operand, shallow, expectedType); if diag != nil {
		return tast.UnaryExpression{}, diag
	}

	if expr.Operator.Kind == token.TokenNotKw {
		// must be boolean.
		if _, ok := operand.Data.Type().(types.TypeBool); !ok {
			return tast.UnaryExpression{}, a.makeExpectedType(types.TypeBool{}, operand.Data.Type(), operand.Base.Token)
		}
	} else if expr.Operator.Kind == token.TokenMinus {
		// must be numeric (int / float).
		if !typeIsNumeric(operand.Data.Type()) {
			return tast.UnaryExpression{}, a.makeExpectedType(types.TypeFloat{}, operand.Data.Type(), operand.Base.Token)
		}
	} else {
		panic(fmt.Sprintf("Internal: Invalid unary operator: '%s'", expr.Operator.Kind))
	}

	return tast.UnaryExpression{
		Operand: operand,
		Operator: expr.Operator,
	}, nil
}

// expected type is 'bool'.
func (a *Analyzer) analyzeLogicalExpr(expr ast.LogicalExpression, shallow bool) (tast.LogicalExpression, diagnostic.Diagnostic) {
	// all logical expressions have their operands as booleans.
	left, right, diag := a.analyzeBinary(expr.Left, expr.Right, expr.Operator, shallow); if diag != nil {
		return tast.LogicalExpression{}, diag
	}

	// expect boolean type
	if _, ok := left.Data.Type().(types.TypeBool); !ok {
		return tast.LogicalExpression{}, a.makeExpectedType(types.TypeBool{}, left.Data.Type(), left.Base.Token)
	}

	return tast.LogicalExpression{
		Left: left,
		Right: right,
		Operator: expr.Operator,
		ShortCircuit: expr.ShortCircuit,
	}, nil
}

func (a *Analyzer) analyzeBinaryExpr(expr ast.BinaryExpression, shallow bool) (tast.BinaryExpression, diagnostic.Diagnostic) {
	left, right, diag := a.analyzeBinary(expr.Left, expr.Right, expr.Operator, shallow)

	return tast.BinaryExpression{
		Left: left,
		Right: right,
		Operator: expr.Operator,
	}, diag
}

func (a *Analyzer) analyzeCallExpr(expr ast.CallExpression, shallow bool, expectedType *types.TypeData) (tast.CallExpression, diagnostic.Diagnostic) {
	callee, diag := a.analyzeExpression(expr.Callee, shallow, expectedType); if diag != nil {
		return tast.CallExpression{}, diag
	}

	// TODO: coerce?
	// for now we won't coerce.

	fn, ok := callee.Data.Type().(types.TypeFunction); if !ok {
		return tast.CallExpression{}, a.makeExpectedCallableType(callee.Data.Type(), callee.Base.Token)
	}

	typedArgs := []tast.Expression{}

	if len(expr.Arguments) != len(fn.Parameters) {
        return tast.CallExpression{}, a.makeExpectedArity(len(fn.Parameters), len(expr.Arguments), expr.Callee.Base.Token)
	}

	for i, param := range expr.Arguments {
		arg, diag := a.analyzeExpression(param, shallow, &fn.Parameters[i].Data); if diag != nil {
			return tast.CallExpression{}, diag
		}

		coerced, ok := coerceExpr(arg, fn.Parameters[i].Data); if !ok {
			return tast.CallExpression{}, a.makeExpectedType(fn.Parameters[i].Data, arg.Data.Type(), arg.Base.Token)
		}

		typedArgs = append(typedArgs, coerced)
	}

	return tast.CallExpression{
		Callee: callee,
		Arguments: typedArgs,
	}, nil
}

func (a *Analyzer) analyzeGroupExpr(expr ast.GroupExpression, shallow bool, expectedType *types.TypeData) (tast.GroupExpression, diagnostic.Diagnostic) {
	inside, diag := a.analyzeExpression(expr.Expr, shallow, expectedType)
			
	return tast.GroupExpression{
		Expr: inside,
	}, diag
}

func (a *Analyzer) analyzeIdentifierExpr(expr ast.IdentifierExpression) (tast.IdentifierExpression, diagnostic.Diagnostic) {
	variable, res := a.resolveVariable(expr.Token.Lexeme); if res == RES_ERROR {
		return tast.IdentifierExpression{}, a.makeNameNotDefined(expr.Token)
	}

	variable.MarkUsed()

	return tast.IdentifierExpression{
		Variable: variable,
	}, nil
}

func (a *Analyzer) analyzeIdentifierAssignExpr(
	expr ast.IdentifierAssignmentExpression,
	shallow bool,
	expectedType *types.TypeData,
) (tast.IdentifierAssignmentExpression, diagnostic.Diagnostic) {
	variable, res := a.resolveVariable(expr.Name.Lexeme); if res == RES_ERROR {
		return tast.IdentifierAssignmentExpression{}, a.makeNameNotDefined(expr.Name)
	}

	assignExpr, diag := a.analyzeExpression(expr.Expr, shallow, expectedType); if diag != nil {
        return tast.IdentifierAssignmentExpression{}, diag
	}

	if variable.IsImmutable() {
		return tast.IdentifierAssignmentExpression{}, a.makeCannotModifyImmutable(expr.Name)
	}

	// TODO: toggle 'used' on?
	variable.MarkModified()

    return tast.IdentifierAssignmentExpression{
		Variable: variable,
		Expr: assignExpr,
	}, nil
}
/*
func (a *Analyzer) analyzeFnExpr(expr ast.FnExpression, shallow bool, expectedType *types.TypeData) (tast.FnExpression, diagnostic.Diagnostic) {
    // the return type is inferred if not annotated, not automatic 'void', like function declarations.
	returnType := types.DummyType(types.TypeUnknown{})

	if expr.Return != nil {
		returnType = *expr.Return
	}

	expectedFn, okFn := (*expectedType).(types.TypeFunction)
	paramTypes := []types.Type{}

	for _, param := range expr.Parameters {
		if param.Type == nil {
			// try to infer the parameter type based on the inferred type, if it's a function.
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
}
*/
func (a *Analyzer) analyzeBlockExpr(expr ast.BlockExpression, shallow bool) (tast.BlockExpression, diagnostic.Diagnostic) {
	return a.analyzeBlock(expr, MODE_NORMAL, shallow)
}

func (a *Analyzer) analyzeIfExpr(expr ast.IfExpression, shallow bool, expectedType *types.TypeData) (tast.IfExpression, diagnostic.Diagnostic) {
	var expectedCondition types.TypeData = types.TypeBool{}
	
	condition, diag := a.analyzeExpression(expr.Condition, shallow, &expectedCondition); if diag != nil {
		return tast.IfExpression{}, diag
	}

	conditionCoerced, ok := coerceExpr(condition, expectedCondition); if !ok {
		return tast.IfExpression{}, a.makeExpectedType(expectedCondition, condition.Data.Type(), condition.Base.Token)
	}

	then, diag := a.analyzeExpression(expr.Then, shallow, expectedType); if diag != nil {
		return tast.IfExpression{}, diag
	}

	if expr.Else == nil {
		// no 'else' clause, type is 'void'.

		thenCoerced, ok := coerceExpr(then, types.TypeVoid{}); if !ok {
			return tast.IfExpression{}, a.makeExpectedType(types.TypeVoid{}, then.Data.Type(), then.Base.Token)
		}

		return tast.IfExpression{
			Condition: conditionCoerced,
			Then: thenCoerced,
			Else: nil,

			IfType: types.TypeVoid{},
		}, nil
	} else {
		// type is merged from 'then' and 'else' blocks.
		else_, diag := a.analyzeExpression(*expr.Else, shallow, expectedType); if diag != nil {
			return tast.IfExpression{}, diag
		}

		merged := mergeTypes(then.Data.Type(), else_.Data.Type())

		// try again, coercing to the other side.
		if merged == nil {
            merged = mergeTypes(else_.Data.Type(), then.Data.Type())
		}
		
		// if they still can't coerce, throw an error.
        if merged == nil {
			return tast.IfExpression{}, a.makeIncompatibleTypes(then.Data.Type(), else_.Data.Type(), then.Base.Token)
		}

		thenCoerced, ok := coerceExpr(then, merged); if !ok {
			return tast.IfExpression{}, a.makeIncompatibleTypes(merged, then.Data.Type(), then.Base.Token)
		}

		elseCoerced, ok := coerceExpr(else_, merged); if !ok {
			return tast.IfExpression{}, a.makeIncompatibleTypes(merged, else_.Data.Type(), else_.Base.Token)
		}

		return tast.IfExpression{
			Condition: conditionCoerced,
			Then: thenCoerced,
			Else: &elseCoerced,

			IfType: merged,
		}, nil
	}
}
