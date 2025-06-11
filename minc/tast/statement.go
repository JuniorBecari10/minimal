package tast

import (
	"go/token"
	"go/types"
	"minc/ast"
)

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
	Return types.Type
	Body BlockExpression
}

type RecordStatement struct {
	Name token.Token
	Fields []Field
	Methods []FnStatement
}

type ReturnStatement struct {
	Expression Expression
}

type OutStatement struct {
	Expression Expression
}

type VarStatement struct {
	Name token.Token
	Init Expression
	Type types.Type
	Immutable bool
}

type WhileStatement struct {
	Condition Expression
	Block     BlockExpression
}

type ForStatement struct {
	Variable token.Token // identifier
	Type types.Type
	Iterable Expression
	Block BlockExpression
}

type ForVarStatement struct {
	Declaration Statement
	Condition Expression
	Increment *Expression // optional
	Block BlockExpression
	Immutable bool
}

type LoopStatement struct {
	Block BlockExpression
}

type BreakStatement ast.BreakStatement
type ContinueStatement ast.ContinueStatement

type ExprStatement struct {
	Expr Expression
}

// ---

func (x RecordStatement) stmt()   {}
func (x FnStatement) stmt()       {}
func (x ReturnStatement) stmt()   {}
func (x OutStatement) stmt()      {}
func (x VarStatement) stmt()      {}
func (x WhileStatement) stmt()    {}
func (x ForStatement) stmt()      {}
func (x ForVarStatement) stmt()   {}
func (x LoopStatement) stmt()     {}
func (x ExprStatement) stmt()     {}
func (x BreakStatement) stmt()    {}
func (x ContinueStatement) stmt() {}
