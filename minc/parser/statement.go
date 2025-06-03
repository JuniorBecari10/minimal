package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) statement(requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
	switch p.current.Kind {
		case token.TokenWhileKw: return p.whileStmt()
		case token.TokenForKw: return p.forStmtCheck()
		case token.TokenLoopKw: return p.loopStmt()
		case token.TokenBreakKw: return p.breakStmt(requireSemicolon)
		case token.TokenContinueKw: return p.continueStmt(requireSemicolon)
		case token.TokenReturnKw: return p.returnStmt(requireSemicolon)
		case token.TokenOutKw: return p.outStmt(requireSemicolon)
		
		default: return p.exprStmt()
	}
}

// ---

func (p *Parser) whileStmt() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance()

	condition, diag := p.parseExpression(); if diag != nil {
		return ast.Statement{}, diag
	}

	block, diag := p.parseBlock(); if diag != nil {
		return ast.Statement{}, diag
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
	name, varType, diag := p.parseVariableBinding(); if diag != nil {
		return ast.Statement{}, diag
	}

	_, diag = p.expectToken(token.TokenInKw); if diag != nil {
		return ast.Statement{}, diag
	}

	iterable, diag := p.parseExpression(); if diag != nil {
		return ast.Statement{}, diag
	}

	block, diag := p.parseBlock(); if diag != nil {
		return ast.Statement{}, diag
	}

	return newStmt(keyword, ast.ForStatement{
		Variable: name,
		Type: varType,
		Iterable: iterable,
		Block: block,
	}), nil
}

func (p *Parser) forVarStmt(keyword token.Token, isLet bool) (ast.Statement, diagnostic.Diagnostic) {
	decl, diag := p.varDecl(isLet, true); if diag != nil {
		return ast.Statement{}, diag
	}

	condition, diag := p.parseExpression(); if diag != nil {
		return ast.Statement{}, diag
	}

	var increment *ast.Expression = nil
	
	if p.match(token.TokenSemicolon) {
		incrementDecl, diag := p.parseExpression(); if diag != nil {
			return ast.Statement{}, diag
		}

		increment = &incrementDecl
	}

	block, diag := p.parseBlock(); if diag != nil {
		return ast.Statement{}, diag
	}

	return newStmt(keyword, ast.ForVarStatement{
		Declaration: decl,
		Condition: condition,
		Increment: increment,
		Block: block,
		Immutable: isLet,
	}), nil
}

func (p *Parser) loopStmt() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance()
	block, diag := p.parseBlock(); if diag != nil {
		return ast.Statement{}, diag
	}

	return newStmt(keyword, ast.LoopStatement{
		Block: block,
	}), nil
}

func (p *Parser) breakStmt(requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance()

	if requireSemicolon {
		diag := p.expectSemicolon(); if diag != nil {
			return ast.Statement{}, diag
		}
	}

	return newStmt(keyword, ast.ContinueStatement{}), nil
}

func (p *Parser) continueStmt(requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
	// To the parser they mean the same thing.
	return p.breakStmt(requireSemicolon)
}

func (p *Parser) returnStmt(requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance()
	var expr *ast.Expression = nil
	
	if !p.match(token.TokenSemicolon) {
		var diag diagnostic.Diagnostic
		exprVal, diag := p.parseExpression(); if diag != nil {
			return ast.Statement{}, diag
		}

		expr = &exprVal
	}

	if requireSemicolon {
		diag := p.expectSemicolon(); if diag != nil {
			return ast.Statement{}, diag
		}
	}
	
	return newStmt(keyword, ast.ReturnStatement{
		Expression: expr,
	}), nil
}

func (p *Parser) outStmt(requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
	// To the parser they mean the same thing.
	return p.returnStmt(requireSemicolon)
}

func (p *Parser) exprStmt() (ast.Statement, diagnostic.Diagnostic) {
	pos := p.current.Pos

	expr, diag := p.parseExpression(); if diag != nil {
		return ast.Statement{}, diag
	}

	diag = p.expectSemicolon(); if diag != nil {
		return ast.Statement{}, diag
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
