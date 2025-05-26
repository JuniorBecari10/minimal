package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/types"
	"minlib/token"
)

func (p *Parser) parseBlock() (ast.BlockExpression, diagnostic.Diagnostic) {
	const START = token.TokenColon

	if p.check(START) {
		// statement-level one-statement blocks do not require semicolons.
		return p.parseOneStmtBlock(START, false)
	} else {
		return p.parseBraceBlock()
	}
}

func (p *Parser) parseFnBlock() (ast.BlockExpression, diagnostic.Diagnostic) {
	const START = token.TokenArrow

	if p.check(START) {
		// function one-statement blocks require semicolons.
		return p.parseOneStmtBlock(START, true)
	} else {
		return p.parseBraceBlock()
	}
}

func (p *Parser) parseOneStmtBlock(start token.TokenKind, requireSemicolon bool) (ast.BlockExpression, diagnostic.Diagnostic) {
	_, diag := p.expectToken(start); if diag != nil {
		return ast.BlockExpression{}, diag
	}

	stmt, diag := p.declaration(true, requireSemicolon); if diag != nil {
		return ast.BlockExpression{}, diag
	}

	return ast.BlockExpression{
		Stmts: []ast.Statement{stmt},
	}, nil
}

func (p *Parser) parseBraceBlock() (ast.BlockExpression, diagnostic.Diagnostic) {
	_, diag := p.expectToken(token.TokenLeftBrace); if diag != nil {
		return ast.BlockExpression{}, diag
	}

	stmts := []ast.Statement{}

	for !p.check(token.TokenRightBrace) {
		decl, diag := p.declaration(true, true); if diag != nil {
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

func (p *Parser) parseFunctionDefinition() ([]ast.Parameter, *types.Type, ast.BlockExpression, diagnostic.Diagnostic) {
	errorRet := func(diag diagnostic.Diagnostic) ([]ast.Parameter, *types.Type, ast.BlockExpression, diagnostic.Diagnostic) {
		return []ast.Parameter{}, nil, ast.BlockExpression{}, diag
	}

	params, diag := p.parseParameters(); if diag != nil {
		return errorRet(diag)
	}

	var returnType *types.Type = nil
	if p.check(token.TokenColon) {
		returnTypeDecl, diag := p.parseTypeAnnotation(); if diag != nil {
			return errorRet(diag)
		}

		returnType = &returnTypeDecl
	}

	body, diag := p.parseFnBlock(); if diag != nil {
		return errorRet(diag)
	}

	return params, returnType, body, nil
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

func (p *Parser) parseVariableBinding() (token.Token, *types.Type, diagnostic.Diagnostic) {
	name, diag := p.expectToken(token.TokenIdentifier); if diag != nil {
		return token.Token{}, nil, diag
	}

	var varType *types.Type = nil

	if p.check(token.TokenColon) {
		var diag diagnostic.Diagnostic
		varTypeDecl, diag := p.parseTypeAnnotation(); if diag != nil {
			return token.Token{}, nil, diag
		}

		varType = &varTypeDecl
	}

	return name, varType, nil
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

func (p *Parser) parseArguments() ([]ast.Expression, diagnostic.Diagnostic) {
	_, diag := p.expectToken(token.TokenLeftParen); if diag != nil {
		return nil, diag
	}
	
	args := []ast.Expression{}

	for !p.match(token.TokenRightParen) {
		expr, diag := p.parseExpression(); if diag != nil {
			return nil, diag
		}

		args = append(args, expr)

		if !p.check(token.TokenRightParen) {
			_, diag := p.expectToken(token.TokenComma); if diag != nil {
				return nil, diag
			}
		}
	}

	return args, nil
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

func (p *Parser) makeAssignment(left, right ast.Expression, operator token.Token) (ast.Expression, diagnostic.Diagnostic) {
	switch lValue := left.Data.(type) {
		case ast.IdentifierExpression: {
			return newExpr(operator, ast.IdentifierAssignmentExpression{
				Name: lValue.Token,
				Expr: right,
			}), nil
		}

		case ast.GetPropertyExpression: {
			return newExpr(lValue.Property, ast.SetPropertyExpression{
				Left: lValue.Left,
				Property: lValue.Property,
				Value: right,
			}), nil
		}

		default:
			return ast.Expression{}, p.makeInvalidAssignmentTargetDiagnostic()
	}
}
