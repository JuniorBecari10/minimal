package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) statement() (ast.Statement, diagnostic.Diagnostic) {
	switch p.current.Kind {
		case token.TokenIfKw: return p.ifStmt()
		case token.TokenWhileKw: return p.whileStmt()
		case token.TokenForKw: return p.forStatCheck()
		case token.TokenLoopKw: return p.loopStmt()
		case token.TokenBreakKw: return p.breakStmt()
		case token.TokenContinueKw: return p.continueStmt()
		case token.TokenReturnKw: return p.returnStmt()
		case token.TokenLeftBrace: return p.blockStmt()
		
		default: return p.exprStmt()
	}
}
