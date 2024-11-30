package ast

import "vm-go/token"

type Statement interface {
	stmt()
}

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
	Then      BlockStatement
	Else      *BlockStatement // optional
}

type WhileStatement struct {
	AstBase
	Condition Expression
	Block     BlockStatement
}

type ForVarStatement struct {
	AstBase
	Declaration VarStatement
	Condition Expression
	Increment *Expression // optional
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

func (x VarStatement) stmt()   {}
func (x BlockStatement) stmt() {}
func (x IfStatement) stmt()    {}
func (x WhileStatement) stmt() {}
func (x ForVarStatement) stmt() {}
func (x PrintStatement) stmt() {}
func (x ExprStatement) stmt()  {}
