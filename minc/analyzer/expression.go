package analyzer

import (
	"fmt"
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
	"minc/types"
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
			expr, diag := a.analyzeRangeExpr(expr, shallow); if diag != nil {
                return tast.Expression{}, diag
			}

			return newExpr(expr), nil
		}
		case ast.AsExpression: {
			expr, diag := a.analyzeAsExpr(expr, shallow); if diag != nil {
				return tast.Expression{}, nil
			}

			return newExpr(expr), nil
		}

		case ast.NilExpression:
			return newExpr(tast.NilExpression{}), nil
		
		case ast.VoidExpression: {
			expr, diag := a.analyzeVoidExpr(expr, shallow, expectedType); if diag != nil {
				return tast.Expression{}, nil
			}

			return newExpr(expr), nil
		}

		case ast.UnaryExpression: {}
		case ast.LogicalExpression: {}
		case ast.BinaryExpression: {}
		case ast.CallExpression: {}
		case ast.GroupExpression: {}
		case ast.IdentifierExpression: {}
		case ast.SelfExpression: {}
		case ast.IdentifierAssignmentExpression: {}
		case ast.FnExpression: {}

		case ast.BlockExpression: {
			expr, diag := a.analyzeBlockExpr(expr, shallow); if diag != nil {
				return tast.Expression{}, nil
			}

			return newExpr(expr), nil
		}
		
		case ast.IfExpression: {}
		case ast.GetPropertyExpression: {}
		case ast.SetPropertyExpression: {}
	}

	panic(fmt.Sprintf("Internal: invalid expression: %#v", e))
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

func (a *Analyzer) analyzeRangeExpr(expr ast.RangeExpression, shallow bool) (tast.RangeExpression, diagnostic.Diagnostic) {
	return tast.RangeExpression{}, nil
}

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

func (a *Analyzer) analyzeBlockExpr(expr ast.BlockExpression, shallow bool) (tast.BlockExpression, diagnostic.Diagnostic) {
	if shallow {
		// don't enter the block; return an empty one with unknown type.
		return tast.BlockExpression{
			Stmts: []tast.Statement{},
			BlockType: types.DummyType(types.TypeUnknown{}),
		}, nil
	} else {
        block, res := a.analyzeBlock(expr.Stmts, MODE_NORMAL); if res == RES_ERROR {
			return tast.BlockExpression{}, diagnostic.HandledDiagnostic{}
		}

		return block, nil
	}
}
