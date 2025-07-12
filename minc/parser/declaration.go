package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minlib/token"
)

func (p *Parser) parseTopLevelDeclaration() (ast.Statement, ParserResult) {
	return p.parseStatement(false, true)
}

// parses a statement and add the diagnostic to the list inside the parser if an error occurs.
// this is the synchronization point; this does not bubble up the error.
// it returns a ParseResult for better error handling on the caller side.
func (p *Parser) parseStatement(allowStatements, requireSemicolon bool) (ast.Statement, ParserResult) {
	if p.current.IsEnd() {
		// Unexpected EOF when a statement was required
		p.hadError = true
		return ast.Statement{}, RES_ERROR
	}

	decl, diag := p.declaration(allowStatements, requireSemicolon); if diag != nil {
		p.diagnostics = append(p.diagnostics, diag)
		p.hadError = true

		p.synchronize()
		return ast.Statement{}, RES_ERROR
	}

	return decl, RES_OK
}

func (p *Parser) declaration(allowStatements, requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
	switch p.current.Kind {
		case token.TokenRecordKw: return p.recordDecl(requireSemicolon)
		case token.TokenFnKw: return p.fnDecl()

		case token.TokenVarKw,
			 token.TokenLetKw:
			 return p.varDecl(p.current.Kind == token.TokenLetKw, requireSemicolon)

		default: {
			if allowStatements {
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

	methods := []ast.FnDeclaration{}

	if p.check(token.TokenLeftBrace) {
		methods, diag = p.parseMethods(); if diag != nil {
			return ast.Statement{}, diag
		}
	} else if requireSemicolon {
		diag := p.expectSemicolon(); if diag != nil {
			return ast.Statement{}, diag
		}
	}

	return newStmt(keyword, ast.RecordDeclaration{
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

	params, returnType, body, isOneStmtBlock, diag := p.parseFunctionDefinition(); if diag != nil {
		return ast.Statement{}, diag
	}

	if isOneStmtBlock {
		// we need to require semicolons.
		diag := p.expectSemicolon(); if diag != nil {
			return ast.Statement{}, diag
		}
	}

	return newStmt(keyword, ast.FnDeclaration{
		Name: name,
		Parameters: params,
		Return: returnType,
		Body: body,
	}), nil
}

func (p *Parser) varDecl(isLet bool, requireSemicolon bool) (ast.Statement, diagnostic.Diagnostic) {
	keyword, _ := p.advance() // Guaranteed.

	name, varType, diag := p.parseVariableBinding(); if diag != nil {
		return ast.Statement{}, diag
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

	return newStmt(keyword, ast.VarDeclaration{
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
			Token: keyword,
		},

		Data: data,
	}
}
