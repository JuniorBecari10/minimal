package tast

import (
	"minc/ast"
	"minc/types"
)

type Expression = ast.Expression

type ExprData interface {
	Type() types.Type
}

type IntExpression ast.IntExpression
type FloatExpression ast.FloatExpression
type StringExpression ast.StringExpression
type CharExpression ast.CharExpression
type BoolExpression ast.BoolExpression

type RangeExpression ast.RangeExpression

type AsExpression ast.AsExpression
type NilExpression ast.NilExpression
type VoidExpression ast.VoidExpression
type UnaryExpression ast.UnaryExpression
type LogicalExpression ast.LogicalExpression
type BinaryExpression ast.BinaryExpression
type CallExpression ast.CallExpression
type GroupExpression ast.GroupExpression
type IdentifierExpression ast.IdentifierExpression
type SelfExpression ast.SelfExpression
type IdentifierAssignmentExpression ast.IdentifierAssignmentExpression
type FnExpression ast.FnExpression
type BlockExpression ast.BlockExpression
type IfExpression ast.IfExpression
type GetPropertyExpression ast.GetPropertyExpression
type SetPropertyExpression ast.SetPropertyExpression

// ---

func (x IntExpression) Type() types.Type { return types.TypeInt{} }
func (x FloatExpression) Type() types.Type { return types.TypeFloat{} }
func (x StringExpression) Type() types.Type { return types.TypeStr{} }
func (x CharExpression) Type() types.Type { return types.TypeChar{} }
func (x BoolExpression) Type() types.Type { return types.TypeBool{} }
func (x RangeExpression) Type() types.Type { return x.Start }
func (x IntExpression) Type() types.Type { return types.TypeInt{} }
func (x IntExpression) Type() types.Type { return types.TypeInt{} }
