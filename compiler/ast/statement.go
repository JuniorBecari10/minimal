package ast

import "vm-go/token"

type Statement struct {
	Base AstBase
	Data StmtData
}

type StmtData interface {
	stmt()
}

type FnStatement struct {
	Name token.Token
	Parameters []Parameter
	Body BlockStatement
}

type RecordStatement struct {
	Name token.Token
	Fields []Field
	Methods []FnStatement
}

type ReturnStatement struct {
	Expression *Expression // optional
}

type VarStatement struct {
	Name token.Token
	Init Expression
}

type BlockStatement struct {
	Stmts []Statement
}

type IfStatement struct {
	Condition Expression
	Then      BlockStatement
	Else      *BlockStatement // optional
}

type WhileStatement struct {
	Condition Expression
	Block     BlockStatement
}

type ForStatement struct {
	Variable token.Token // identifier
	Iterable Expression
	Block BlockStatement
}

type ForVarStatement struct {
	Declaration Statement
	Condition Expression
	Increment *Expression // optional
	Block BlockStatement
}

type LoopStatement struct {
	Block BlockStatement
}

type BreakStatement struct {
	Token token.Token
}

type ContinueStatement struct {
	Token token.Token
}

type ExprStatement struct {
	Expr Expression
}

// ---

func (x RecordStatement) stmt()   {}
func (x FnStatement) stmt()       {}
func (x ReturnStatement) stmt()   {}
func (x VarStatement) stmt()      {}
func (x BlockStatement) stmt()    {}
func (x IfStatement) stmt()       {}
func (x WhileStatement) stmt()    {}
func (x ForStatement) stmt()      {}
func (x ForVarStatement) stmt()   {}
func (x LoopStatement) stmt()     {}
func (x ExprStatement) stmt()     {}
func (x BreakStatement) stmt()    {}
func (x ContinueStatement) stmt() {}
