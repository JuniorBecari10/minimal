package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) statement() (ast.Statement, diagnostic.Diagnostic) {
	switch p.current.Kind {
		case token.TokenWhileKw: return p.whileStmt()
		case token.TokenForKw: return p.forStmtCheck()
		case token.TokenLoopKw: return p.loopStmt()
		case token.TokenBreakKw: return p.breakStmt()
		case token.TokenContinueKw: return p.continueStmt()
		case token.TokenReturnKw: return p.returnStmt()
		case token.TokenOutKw: return p.outStmt()
		
		default: return p.exprStmt()
	}
}

// ---

func (p *Parser) whileStmt() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance()

	condition, diag := p.parseExpression(); if diag != nil {
		return ast.Statement{}, nil
	}

	block, diag := p.parseBlock(); if diag != nil {
		return ast.Statement{}, nil
	}

	return newStmt(keyword, ast.WhileStatement{
		Condition: condition,
		Block: block,
	}), nil
}

func (p *Parser) forStmtCheck() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance()

	if p.check(token.TokenVarKw) || p.check(token.TokenLetKw) {
		return p.forVarStmt(keyword, p.current.Kind == token.TokenLetKw)
	} else {
		return p.forStmt(keyword)
	}
}

func (p *Parser) forStmt(keyword token.Token) (ast.Statement, diagnostic.Diagnostic) {

}

func (p *Parser) forVarStmt(keyword token.Token, isLet bool) (ast.Statement, diagnostic.Diagnostic) {
	decl, diag := p.varDecl(isLet); if diag != nil {
		return ast.Statement{}, nil
	}

	condition, diag := p.parseExpression(); if diag != nil {
		return ast.Statement{}, nil
	}

	var increment *ast.Expression = nil
	
	if p.match(token.TokenSemicolon) {
		*increment, diag = p.parseExpression(); if diag != nil {
			return ast.Statement{}, nil
		}
	}

	block, diag := p.parseBlock(); if diag != nil {
		return ast.Statement{}, nil
	}

	return newStmt(keyword, ast.ForVarStatement{
		Declaration: decl,
		Condition: condition,
		Increment: increment,
		Block: block,
	}), nil
}

func (p *Parser) loopStmt() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance()
	block, diag := p.parseBlock(); if diag != nil {
		return ast.Statement{}, nil
	}

	return newStmt(keyword, ast.LoopStatement{
		Block: block,
	}), nil
}

func (p *Parser) breakStmt() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance()

	diag := p.expectSemicolon(); if diag != nil {
		return ast.Statement{}, nil
	}

	return newStmt(keyword, ast.ContinueStatement{}), nil
}

func (p *Parser) continueStmt() (ast.Statement, diagnostic.Diagnostic) {
	// To the parser they mean the same thing.
	return p.breakStmt()
}

func (p *Parser) returnStmt() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance()
	var expr *ast.Expression = nil
	
	if p.match(token.TokenSemicolon) {
		var diag diagnostic.Diagnostic
		*expr, diag = p.parseExpression(); if diag != nil {
			return ast.Statement{}, nil
		}
	}
	
	return newStmt(keyword, ast.ReturnStatement{
		Expression: expr,
	}), nil
}

func (p *Parser) outStmt() (ast.Statement, diagnostic.Diagnostic) {
	// To the parser they mean the same thing.
	return p.returnStmt()
}

func (p *Parser) exprStmt() (ast.Statement, diagnostic.Diagnostic) {
	pos := p.current.Pos

	expr, diag := p.parseExpression(); if diag != nil {
		return ast.Statement{}, nil
	}

	diag = p.expectSemicolon(); if diag != nil {
		return ast.Statement{}, nil
	}

	return ast.Statement{
		Base: ast.AstBase{
			Pos:    pos,
			Length: expr.Base.Length,
		},

		Data: ast.ExprStatement{
			Expr: expr,
		},
	}, nil
}
