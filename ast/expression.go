package ast

import (
	"vm-go/token"
)

type Expression interface {
	expr()
}

type NumberExpression struct {
	AstBase
	Literal float64
}

type StringExpression struct {
	AstBase
	Literal string
}

type BoolExpression struct {
	AstBase
	Literal bool
}

type NilExpression struct {
	AstBase
}

type VoidExpression struct {
	AstBase
}

type UnaryExpression struct {
	AstBase

	Operand  Expression
	Operator token.Token
}

type BinaryExpression struct {
	AstBase

	Left     Expression
	Right    Expression
	Operator token.Token
}

type CallExpression struct {
	AstBase
	Callee Expression
	Arguments []Expression
}

type GroupExpression struct {
	AstBase
	Expr Expression
}

type IdentifierExpression struct {
	AstBase
	Ident token.Token
}

type IdentifierAssignmentExpression struct {
	AstBase
	Name token.Token
	Expr Expression
}

// ---

func (x NumberExpression) expr()     {}
func (x StringExpression) expr()     {}
func (x BoolExpression) expr()       {}
func (x NilExpression) expr()        {}
func (x VoidExpression) expr()        {}
func (x UnaryExpression) expr()      {}
func (x BinaryExpression) expr()     {}
func (x GroupExpression) expr()      {}
func (x CallExpression) expr()      {}
func (x IdentifierExpression) expr() {}
func (x IdentifierAssignmentExpression) expr() {}
