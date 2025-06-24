package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) statement(requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
	switch p.current.Kind {
		case token.TokenWhileKw: return p.whileStmt()
		case token.TokenForKw: return p.forStmt()
		case token.TokenLoopKw: return p.loopStmt()
		case token.TokenBreakKw: return p.breakStmt(requireSemicolon)
		case token.TokenContinueKw: return p.continueStmt(requireSemicolon)
		case token.TokenReturnKw: return p.returnStmt(requireSemicolon)
		case token.TokenOutKw: return p.outStmt(requireSemicolon)
		
		default: return p.exprStmt(requireSemicolon)
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

func (p *Parser) forStmt() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance()

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

	// Check for a semicolon *immediately* after return
	if p.match(token.TokenSemicolon) {
		return newStmt(keyword, ast.ReturnStatement{
			Expression: nil,
		}), nil
	}

	exprVal, diag := p.parseExpression(); if diag != nil {
		return ast.Statement{}, diag
	}
	expr = &exprVal

	// Now we require semicolon *after* the return value
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

func (p *Parser) exprStmt(requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
	// save current.Pos if needed.

	// TODO: maybe pass a flag here saying that this is supposed to be a statement, and therefore refine the error message,
	// to say 'expected statement', instead of 'expected expression'.
	expr, diag := p.parseExpression(); if diag != nil {
		return ast.Statement{}, diag
	}

	var semicolon *token.Token = nil

	// If this is the last statement in a brace block (current token is '}'), don't require the semicolon.
	// Don't skip it, since it will be required to close the block later.
	if requireSemicolon && !p.check(token.TokenRightBrace) && !p.exprCanSkipSemicolon(expr.Data) {
		semicolon = &p.current
		
		diag = p.expectSemicolon(); if diag != nil {
			return ast.Statement{}, diag
		}
	}

	return ast.Statement{
		Base: expr.Base,

		Data: ast.ExprStatement{
			Expr: expr,
			Semicolon: semicolon,
		},
	}, nil
}
