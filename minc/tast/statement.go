package tast

import (
	"minc/types"
	"minc/ast"
	"minlib/token"
)

type Statement struct {
	Base AstBase
	Data StmtData
}

type StmtData interface {
	stmt()
}

type FnDeclaration struct {
	Name token.Token
	Parameters []Parameter
	Return types.Type
	Body BlockExpression
}

type RecordDeclaration struct {
	Name token.Token
	Fields []Field
	Methods []FnDeclaration
}

// if in the AST the expression is omitted, here it should be void.
type ReturnStatement struct {
	Expression Expression
}

type OutStatement struct {
	Expression Expression
}

type VarDeclaration struct {
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
	VariableType types.TypeData
	Iterable Expression
	Block BlockExpression
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

func (x RecordDeclaration) stmt() {}
func (x FnDeclaration) stmt()     {}
func (x ReturnStatement) stmt()   {}
func (x OutStatement) stmt()      {}
func (x VarDeclaration) stmt()    {}
func (x WhileStatement) stmt()    {}
func (x ForStatement) stmt()      {}
func (x LoopStatement) stmt()     {}
func (x ExprStatement) stmt()     {}
func (x BreakStatement) stmt()    {}
func (x ContinueStatement) stmt() {}
