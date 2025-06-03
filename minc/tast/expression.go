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
	Type() types.Type
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

func (x IntExpression) Type() types.Type { return types.TypeInt{} }
func (x FloatExpression) Type() types.Type { return types.TypeFloat{} }
func (x StringExpression) Type() types.Type { return types.TypeStr{} }
func (x CharExpression) Type() types.Type { return types.TypeChar{} }
func (x BoolExpression) Type() types.Type { return types.TypeBool{} }
func (x RangeExpression) Type() types.Type { return types.TypeRange{ Inside: x.Start.Data.Type() } }
func (x AsExpression) Type() types.Type { return types.TypeInt{} } // TODO: add helper function that calculates the actual type
func (x NilExpression) Type() types.Type { return types.TypeUntypedNil{} }
func (x VoidExpression) Type() types.Type { return types.TypeVoid{} }
func (x UnaryExpression) Type() types.Type { return x.Operand.Data.Type() } // the operator doesn't change the type.
func (x LogicalExpression) Type() types.Type { return x.Left.Data.Type() } // same. assumes both sides have the same type
func (x BinaryExpression) Type() types.Type { return x.Left.Data.Type() } // same

func (x CallExpression) Type() types.Type {
    if fn, ok := x.Callee.Data.Type().(types.TypeFunction); ok {
        return fn.Return
    }

	// should not happen
    return types.TypeUnknown{}
}

func (x GroupExpression) Type() types.Type { return x.Expr.Data.Type() }
func (x IdentifierExpression) Type() types.Type { return x.VariableType }
func (x SelfExpression) Type() types.Type { return x.VariableType }
func (x IdentifierAssignmentExpression) Type() types.Type { return types.TypeVoid{} } // assignments return void

func (x FnExpression) Type() types.Type {
	parameterTypes := make([]types.Type, 0, len(x.Parameters))

	for _, param := range x.Parameters {
		parameterTypes = append(parameterTypes, param.Type)
	}

	return types.TypeFunction{
		Parameters: parameterTypes,
		Return: x.Return,
	}
}

func (x BlockExpression) Type() types.Type { return x.BlockType }

func (x IfExpression) Type() types.Type {
	if x.Else == nil {
		return types.TypeVoid{}
	} else {
		return x.Then.Data.Type()
	}
}

func (x GetPropertyExpression) Type() types.Type { return x.PropertyType }
func (x SetPropertyExpression) Type() types.Type { return types.TypeVoid{} }
