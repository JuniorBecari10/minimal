package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) parseBlock() (ast.BlockExpression, diagnostic.Diagnostic) {
	_, diag := p.expectToken(token.TokenLeftBrace); if diag != nil {
		return ast.BlockExpression{}, diag
	}

	stmts := []ast.Statement{}

	for !p.check(token.TokenRightBrace) {
		decl, diag := p.declaration(true); if diag != nil {
			return ast.BlockExpression{}, diag
		}

		stmts = append(stmts, decl)
	}

	_, diag = p.expectToken(token.TokenRightBrace); if diag != nil {
		return ast.BlockExpression{}, diag
	}

	return ast.BlockExpression{
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

		methods = append(methods, method.Data.(ast.FnStatement)) // 'fnDecl' always returns a FnStatement.
	}

	_, diag = p.expectToken(token.TokenRightBrace); if diag != nil {
		return []ast.FnStatement{}, diag
	}

	return methods, nil
}

func (p *Parser) parseParameters() ([]ast.Parameter, diagnostic.Diagnostic) {
	_, diag := p.expectToken(token.TokenLeftParen); if diag != nil {
		return []ast.Parameter{}, diag
	}

	params := []ast.Parameter{}

	for !p.match(token.TokenRightParen) {
		name, diag := p.expectToken(token.TokenIdentifier); if diag != nil {
			return []ast.Parameter{}, diag
		}

		paramType, diag := p.parseTypeAnnotation(); if diag != nil {
			return []ast.Parameter{}, diag
		}

		params = append(params, ast.Parameter{
			Name: name,
			Type: paramType,
		})

		if !p.check(token.TokenRightParen) {
			_, diag := p.expectToken(token.TokenComma); if diag != nil {
				return []ast.Parameter{}, diag
			}
		}
	}

	return params, nil
}

func (p *Parser) parseFields() ([]ast.Field, diagnostic.Diagnostic) {
	params, diag := p.parseParameters(); if diag != nil {
		return []ast.Field{}, diag
	}

	fields := []ast.Field{}

	for _, param := range params {
		fields = append(fields, ast.Field{
			Name: param.Name,
			Type: param.Type,
		})
	}
	
	return fields, nil
}
