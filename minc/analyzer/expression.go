package analyzer

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
)

func (a *Analyzer) analyzeExpression(expr ast.Expression, shallow bool) (tast.Expression, diagnostic.Diagnostic) {
	switch expr.Data.(type) {
		case ast.IntExpression: {}
		case ast.FloatExpression: {}
		case ast.StringExpression: {}
		case ast.CharExpression: {}
		case ast.BoolExpression: {}
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
}
