package tast

import (
	"minc/ast"
	"minc/types"
	"minlib/token"
)

// the declaration is the same, but the interface it refers to is different.
type Expression struct {
	Base AstBase
	Data ExprData
}

type ExprData interface {
	Type() types.TypeData
}

// these are primitive wrappers, so there's no problem
type IntExpression ast.IntExpression
type FloatExpression ast.FloatExpression
type StringExpression ast.StringExpression
type CharExpression ast.CharExpression
type BoolExpression ast.BoolExpression

type RangeExpression struct {
	Start Expression
	End Expression
	Step *Expression // optional

    Inclusive bool
}

type AsExpression struct {
	Operand Expression
	ConvertType types.Type
}

// almost the same as AsExpression (the type does not contain token information), but represents an implicit, restricted conversion.
type CoerceExpression struct {
	Operand Expression
	ConvertType types.TypeData
}

type NilExpression ast.NilExpression

type VoidExpression struct {
	Expr *Expression // optional
}

type UnaryExpression struct {
	Operand  Expression
	Operator token.Token
}

type LogicalExpression struct {
	Left     Expression
	Right    Expression
	
	Operator token.Token
	ShortCircuit bool
}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator token.Token
}

type CallExpression struct {
	Callee Expression
	Arguments []Expression
}

type GroupExpression struct {
	Expr Expression
}

type IdentifierExpression struct {
	Token token.Token
	VariableType types.Type // filled during type-checking and name-resolution phase
}

type SelfExpression struct {
	VariableType types.Type // filled during type-checking and name-resolution phase
}

type IdentifierAssignmentExpression struct {
	Name token.Token
	Expr Expression
}

type FnExpression struct {
	Parameters []Parameter
	Return types.Type
	Body BlockExpression
}

type BlockExpression struct {
	Stmts []Statement
	BlockType types.Type
}

type IfExpression struct {
	Condition Expression
	Then Expression
	Else *Expression // optional
}

type GetPropertyExpression struct {
	Left Expression
	Property token.Token

	PropertyType types.Type
}

type SetPropertyExpression struct {
	Left Expression
	Property token.Token
	Value Expression
}

// ---

func (IntExpression) Type() types.TypeData    { return types.TypeInt{} }
func (FloatExpression) Type() types.TypeData  { return types.TypeFloat{} }
func (StringExpression) Type() types.TypeData { return types.TypeStr{} }
func (CharExpression) Type() types.TypeData   { return types.TypeChar{} }
func (BoolExpression) Type() types.TypeData   { return types.TypeBool{} }

func (x RangeExpression) Type() types.TypeData {
	return types.TypeRange{
		Inside: types.DummyType(x.Start.Data.Type()),
	}
}

func (x AsExpression) Type() types.TypeData {
	// TODO: add helper function that calculates the actual type
	return types.TypeInt{}
}

func (x CoerceExpression) Type() types.TypeData {
	// TODO: add helper function that calculates the actual type
	return types.TypeInt{}
}

func (NilExpression) Type() types.TypeData   { return types.TypeUntypedNil{} }
func (VoidExpression) Type() types.TypeData  { return types.TypeVoid{} }

func (x UnaryExpression) Type() types.TypeData {
	return x.Operand.Data.Type()
}

func (x LogicalExpression) Type() types.TypeData {
	return x.Left.Data.Type()
}

func (x BinaryExpression) Type() types.TypeData {
	return x.Left.Data.Type()
}

func (x CallExpression) Type() types.TypeData {
	if fn, ok := x.Callee.Data.Type().(types.TypeFunction); ok {
		return fn.Return.Data
	}
	return types.TypeUnknown{}
}

func (x GroupExpression) Type() types.TypeData {
	return x.Expr.Data.Type()
}

func (x IdentifierExpression) Type() types.TypeData {
	return x.VariableType.Data
}

func (x SelfExpression) Type() types.TypeData {
	return x.VariableType.Data
}

func (IdentifierAssignmentExpression) Type() types.TypeData {
	return types.TypeVoid{}
}

func (x FnExpression) Type() types.TypeData {
	parameterTypes := make([]types.Type, 0, len(x.Parameters))
	
	for _, param := range x.Parameters {
		type_ := types.DummyType(types.TypeUnknown{})

		if param.Type != nil {
            type_ = *param.Type
		}

		parameterTypes = append(parameterTypes, type_)
	}

	return types.TypeFunction{
		Parameters: parameterTypes,
		Return:     x.Return,
	}
}

func (x BlockExpression) Type() types.TypeData {
	return x.BlockType.Data
}

func (x IfExpression) Type() types.TypeData {
	if x.Else == nil {
		return types.TypeVoid{}
	}
	return x.Then.Data.Type()
}

func (x GetPropertyExpression) Type() types.TypeData {
	return x.PropertyType.Data
}

func (SetPropertyExpression) Type() types.TypeData {
	return types.TypeVoid{}
}
