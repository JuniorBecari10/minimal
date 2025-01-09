package ast

import (
	"vm-go/token"
)

type Expression struct {
	Base AstBase
	Data ExprData
}

type ExprData interface {
	expr()
}

type NumberExpression struct {
	Literal float64
}

type StringExpression struct {
	Literal string
}

type BoolExpression struct {
	Literal bool
}

type RangeExpression struct {
	Start Expression
	End Expression
	Step *Expression // optional
}

type NilExpression struct {}
type VoidExpression struct {}

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

type SelfExpression struct {
	Token token.Token
}

type IdentifierAssignmentExpression struct {
	Name token.Token
	Expr Expression
}

type FnExpression struct {
	Parameters []Parameter
	Body BlockStatement
}

type IfExpression struct {
	Condition Expression
	Then Expression
	Else Expression
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

func (x NumberExpression) expr()     {}
func (x StringExpression) expr()     {}
func (x BoolExpression) expr()       {}
func (x NilExpression) expr()        {}
func (x VoidExpression) expr()        {}
func (x RangeExpression) expr() {}
func (x UnaryExpression) expr()      {}
func (x LogicalExpression) expr()     {}
func (x BinaryExpression) expr()     {}
func (x GroupExpression) expr()      {}
func (x CallExpression) expr()      {}
func (x IdentifierExpression) expr() {}
func (x SelfExpression) expr() {}
func (x IdentifierAssignmentExpression) expr() {}
func (x FnExpression) expr() {}
func (x IfExpression) expr() {}
func (x GetPropertyExpression) expr() {}
func (x SetPropertyExpression) expr() {}
