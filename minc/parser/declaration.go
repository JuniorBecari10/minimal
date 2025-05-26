package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) declaration(allowStats, requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
	switch p.current.Kind {
		case token.TokenRecordKw: return p.recordDecl(requireSemicolon)
		case token.TokenFnKw: return p.fnDecl()

		case token.TokenVarKw,
			 token.TokenLetKw:
			 return p.varDecl(p.current.Kind == token.TokenLetKw, requireSemicolon)

		default: {
			if allowStats {
				return p.statement(requireSemicolon)
			} else {
				diag := p.makeStatementsNotAllowedDiagnostic()
				p.advance()
				
				return ast.Statement{}, diag
			}
		}
	}
}

// ---

func (p *Parser) recordDecl(requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
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
	} else if requireSemicolon {
		diag := p.expectSemicolon(); if diag != nil {
			return ast.Statement{}, diag
		}
	}

	return newStmt(keyword, ast.RecordStatement{
		Name:   name,
		Fields: fields,
		Methods: methods,
	}), nil
}

// function declarations do not require a semicolon.
func (p *Parser) fnDecl() (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance() // Guaranteed.

	name, diag := p.expectToken(token.TokenIdentifier); if diag != nil {
		return ast.Statement{}, diag
	}

	params, returnType, body, diag := p.parseFunctionDefinition(); if diag != nil {
		return ast.Statement{}, nil
	}

	return newStmt(keyword, ast.FnStatement{
		Name: name,
		Parameters: params,
		Return: returnType,
		Body: body,
	}), nil
}

func (p *Parser) varDecl(isLet bool, requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance() // Guaranteed.

	name, varType, diag := p.parseVariableBinding(); if diag != nil {
		return ast.Statement{}, nil
	}

	_, diag = p.expectToken(token.TokenEqual); if diag != nil {
		return ast.Statement{}, diag
	}

	expr, diag := p.parseExpression(); if diag != nil {
		return ast.Statement{}, diag
	}

	if requireSemicolon {
		diag = p.expectSemicolon(); if diag != nil {
			return ast.Statement{}, diag
		}
	}

	return newStmt(keyword, ast.VarStatement{
		Name: name,
		Init: expr,
		Type: varType,
		Immutable: isLet,
	}), nil
}

// ---

func newStmt(keyword token.Token, data ast.StmtData) ast.Statement {
	return ast.Statement{
		Base: ast.AstBase{
			Pos:    keyword.Pos,
			Length: len(keyword.Lexeme),
		},

		Data: data,
	}
}
