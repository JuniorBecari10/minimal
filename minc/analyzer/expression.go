package analyzer

import (
	"fmt"
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
)

func (a *Analyzer) analyzeExpression(e ast.Expression, shallow bool) (tast.Expression, diagnostic.Diagnostic) {
	newExpr := func(data tast.ExprData) tast.Expression {
		return tast.Expression{
			Base: tast.AstBase(e.Base),
			Data: data,
		}
	}

	switch expr := e.Data.(type) {
		case ast.IntExpression:
			return newExpr(tast.IntExpression{
				Literal: expr.Literal,
			}), nil

		case ast.FloatExpression:
			return newExpr(tast.FloatExpression{
				Literal: expr.Literal,
			}), nil
		
		case ast.StringExpression:
			return newExpr(tast.StringExpression{
				Literal: expr.Literal,
			}), nil
		
		case ast.CharExpression:
			return newExpr(tast.CharExpression{
				Literal: expr.Literal,
			}), nil
		
		
		case ast.BoolExpression:
			return newExpr(tast.BoolExpression{
				Literal: expr.Literal,
			}), nil
		
		case ast.RangeExpression: {}
		case ast.AsExpression: {}
		case ast.NilExpression: {}
		case ast.VoidExpression: {}
		case ast.UnaryExpression: {}
		case ast.LogicalExpression: {}
		case ast.BinaryExpression: {}
		case ast.CallExpression: {}
		case ast.GroupExpression: {}
		case ast.IdentifierExpression: {}
		case ast.SelfExpression: {}
		case ast.IdentifierAssignmentExpression: {}
		case ast.FnExpression: {}
		case ast.BlockExpression: {}
		case ast.IfExpression: {}
		case ast.GetPropertyExpression: {}
		case ast.SetPropertyExpression: {}
	}

	panic(fmt.Sprintf("Internal: invalid expression: %#v", e))
}
