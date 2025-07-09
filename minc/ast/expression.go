package ast

import (
	"minc/types"
	"minlib/token"
)

type Expression struct {
	Base AstBase
	Data ExprData
}

type ExprData interface {
	expr()
}

type IntExpression struct {
	Literal int32
}

type FloatExpression struct {
	Literal float64
}

type StringExpression struct {
	Literal string
}

type CharExpression struct {
	Literal uint8
}

type BoolExpression struct {
	Literal bool
}

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

type NilExpression struct {
	Token token.Token
	TypeArguments []types.Type
}

type SomeExpression struct {
	Inside Expression
	TypeArguments []types.Type
}

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
}

type SelfExpression struct {}

type IdentifierAssignmentExpression struct {
	Name token.Token
	Expr Expression
}

type FnExpression struct {
	Parameters []Parameter
	Return *types.Type // optional
	Body BlockExpression
}

type BlockExpression struct {
	Stmts []Statement
	Token token.Token
	ShowSemicolonWarning bool
}

type IfExpression struct {
	Condition Expression
	Then Expression
	Else *Expression // optional
}

type GetPropertyExpression struct {
	Left Expression
	Property token.Token
}

type SetPropertyExpression struct {
	Left Expression
	Property token.Token
	Value Expression
}

// ---

func (b BlockExpression) IntoExpr() Expression {
	return Expression{
		Base: AstBase{
			Token: b.Token,
		},
		Data: b,
	}
}

// ---

func (x IntExpression) expr()                  {}
func (x FloatExpression) expr()                {}
func (x CharExpression) expr()                 {}
func (x StringExpression) expr()               {}
func (x BoolExpression) expr()                 {}
func (x NilExpression) expr()                  {}
func (x SomeExpression) expr()                 {}
func (x VoidExpression) expr()                 {}
func (x RangeExpression) expr()                {}
func (x AsExpression) expr()                   {}
func (x UnaryExpression) expr()                {}
func (x LogicalExpression) expr()              {}
func (x BinaryExpression) expr()               {}
func (x GroupExpression) expr()                {}
func (x CallExpression) expr()                 {}
func (x IdentifierExpression) expr()           {}
func (x SelfExpression) expr()                 {}
func (x IdentifierAssignmentExpression) expr() {}
func (x FnExpression) expr()                   {}
func (x BlockExpression) expr()                {}
func (x IfExpression) expr()                   {}
func (x GetPropertyExpression) expr()          {}
func (x SetPropertyExpression) expr()          {}
