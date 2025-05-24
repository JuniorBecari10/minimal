package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) parseBlock() (ast.BlockStatement, diagnostic.Diagnostic) {
	_, diag := p.expectToken(token.TokenLeftBrace); if diag != nil {
		return ast.BlockStatement{}, diag
	}

	stmts := []ast.Statement{}

	for !p.current.IsEnd() && !p.check(token.TokenRightBrace) {
		decl, diag := p.declaration(true); if diag != nil {
			return ast.BlockStatement{}, diag
		}

		stmts = append(stmts, decl)
	}

	_, diag = p.expectToken(token.TokenRightBrace); if diag != nil {
		return ast.BlockStatement{}, diag
	}

	return ast.BlockStatement{
		Stmts: stmts,
	}, nil
}

func (p *Parser) parseMethods() ([]ast.FnStatement, diagnostic.Diagnostic) {
	_, diag := p.expectToken(token.TokenLeftBrace); if diag != nil {
		return []ast.FnStatement{}, diag
	}

	methods := []ast.FnStatement{}

	for !p.current.IsEnd() && !p.check(token.TokenRightBrace) {
		method, diag := p.fnDecl(); if diag != nil {
			return []ast.FnStatement{}, diag
		}

		methods = append(methods, method.Data.(ast.FnStatement))
	}

	_, diag = p.expectToken(token.TokenRightBrace); if diag != nil {
		return []ast.FnStatement{}, diag
	}

	return methods, nil
}
