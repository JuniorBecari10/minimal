package ast

import (
	"vm-go/token"
)

type Statement interface {
	stmt()
}

type Expression interface {
	expr()
}

type AstBase struct {
	Pos token.Position
}

// ---

type VarStatement struct {
	AstBase
	Name token.Token
	Init Expression
}

type BlockStatement struct {
	AstBase
	Stmts []Statement
}

type IfStatement struct {
	AstBase
	Condition Expression
	Then BlockStatement
	Else *BlockStatement // optional
}

type WhileStatement struct {
	AstBase
	Condition Expression
	Block BlockStatement
}

type PrintStatement struct {
	AstBase
	Expr Expression
}

type ExprStatement struct {
	AstBase
	Expr Expression
}

// ---

func (x VarStatement) stmt() {}
func (x BlockStatement) stmt() {}
func (x IfStatement) stmt() {}
func (x WhileStatement) stmt() {}
func (x PrintStatement) stmt() {}
func (x ExprStatement) stmt() {}

// ---

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

type UnaryExpression struct {
	AstBase

	Operand Expression
	Operator token.Token
}

type BinaryExpression struct {
	AstBase

	Left Expression
	Right Expression
	Operator token.Token
}

type GroupExpression struct {
	AstBase
	Expr Expression
}

type IdentifierExpression struct {
	AstBase
	Ident token.Token
}

// ---

func (x NumberExpression) expr() {}
func (x StringExpression) expr() {}
func (x BoolExpression) expr() {}
func (x NilExpression) expr() {}
func (x UnaryExpression) expr() {}
func (x BinaryExpression) expr() {}
func (x GroupExpression) expr() {}
func (x IdentifierExpression) expr() {}
