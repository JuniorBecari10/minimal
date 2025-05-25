package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/types"
	"minlib/token"
)

func (p *Parser) declaration(allowStats bool) (ast.Statement, diagnostic.Diagnostic) {
	switch p.current.Kind {
		case token.TokenRecordKw: return p.recordDecl()
		case token.TokenFnKw: return p.fnDecl()
		case token.TokenVarKw: return p.varDecl()

		default: {
			if allowStats {
				return p.statement()
			} else {
				diag := p.makeStatementsNotAllowedDiagnostic()
				p.advance()
				
				return ast.Statement{}, diag
			}
		}
	}
}

func (p *Parser) recordDecl() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance() // Guaranteed.

	name, diag := p.expectToken(token.TokenIdentifier); if diag != nil {
		return ast.Statement{}, diag
	}

	fields, diag := p.parseFields(); if diag != nil {
		return ast.Statement{}, diag
	}

	methods := []ast.FnStatement{}

	if p.check(token.TokenLeftBrace) {
		methods, diag = p.parseMethods(); if diag != nil {
			return ast.Statement{}, diag
		}
	} else {
		diag := p.expectSemicolon(); if diag != nil {
			return ast.Statement{}, diag
		}
	}

	return ast.Statement{
		Base: ast.AstBase{
			Pos:    keyword.Pos,
			Length: len(keyword.Lexeme),
		},

		Data: ast.RecordStatement{
			Name:   name,
			Fields: fields,
			Methods: methods,
		},
	}, nil
}

func (p *Parser) fnDecl() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance() // Guaranteed.

	name, diag := p.expectToken(token.TokenIdentifier); if diag != nil {
		return ast.Statement{}, diag
	}

	params, diag := p.parseParameters(); if diag != nil {
		return ast.Statement{}, diag
	}

	var returnType *types.Type = nil
	if p.check(token.TokenColon) {
		*returnType, diag = p.parseTypeAnnotation(); if diag != nil {
			return ast.Statement{}, diag
		}
	}

	body, diag := p.parseBlock(); if diag != nil {
		return ast.Statement{}, diag
	}

	return ast.Statement{
		Base: ast.AstBase{
			Pos:    keyword.Pos,
			Length: len(keyword.Lexeme),
		},

		Data: ast.FnStatement{
			Name: name,
			Parameters: params,
			Body: body,
			ReturnType: returnType,
		},
	}, nil
}

func (p *Parser) varDecl() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance() // Guaranteed.

	name, diag := p.expectToken(token.TokenIdentifier); if diag != nil {
		return ast.Statement{}, diag
	}

	var varType *types.Type = nil
	if p.check(token.TokenColon) {
		var diag diagnostic.Diagnostic
		*varType, diag = p.parseTypeAnnotation(); if diag != nil {
			return ast.Statement{}, diag
		}
	}

	_, diag = p.expectToken(token.TokenEqual); if diag != nil {
		return ast.Statement{}, diag
	}

	expr, diag := p.parseExpression(); if diag != nil {
		return ast.Statement{}, diag
	}

	diag = p.expectSemicolon(); if diag != nil {
		return ast.Statement{}, diag
	}

	return ast.Statement{
		Base: ast.AstBase{
			Pos:    keyword.Pos,
			Length: len(keyword.Lexeme),
		},

		Data: ast.VarStatement{
			Name: name,
			Init: expr,
			Type: varType,
		},
	}, nil
}
