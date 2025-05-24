package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) declaration(allowStats bool) (ast.Statement, diagnostic.Diagnostic) {
	if p.panicMode {
		p.synchronize()
	}

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
	keyword, diag := p.advance(); if diag != nil {
		return ast.Statement{}, diag
	}

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
		diag := p.requireSemicolon(); if diag != nil {
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

}

func (p *Parser) varDecl() (ast.Statement, diagnostic.Diagnostic) {

}
