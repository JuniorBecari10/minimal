package ast

import "vm-go/token"

type Statement interface {
	stmt()
}

type FnStatement struct {
	AstBase
	Name token.Token
	Parameters []Parameter
	Body BlockStatement
}

type ReturnStatement struct {
	AstBase
	Expression *Expression // optional
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

func (x FnStatement) stmt() {}
func (x ReturnStatement) stmt() {}
func (x VarStatement) stmt()   {}
func (x BlockStatement) stmt() {}
func (x IfStatement) stmt()    {}
func (x WhileStatement) stmt() {}
func (x ForVarStatement) stmt() {}
func (x PrintStatement) stmt() {}
func (x ExprStatement) stmt()  {}
