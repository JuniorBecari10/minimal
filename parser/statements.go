package parser

import (
	"vm-go/ast"
	"vm-go/token"
)

func (p *Parser) statement() ast.Statement {
	if p.panicMode {
		p.synchronize()
	}

	t := p.peek(0)
	
	switch t.Kind {
		case token.TokenIfKw: return p.ifStatement()
		case token.TokenWhileKw: return p.whileStatement()
		case token.TokenForKw: return p.forStatement()
		case token.TokenVarKw: return p.varStatement()
		case token.TokenPrintKw: return p.printStatement()
		case token.TokenLeftBrace: return p.blockStatement()
		
		default: return p.exprStatement()
	}
}

// ---

func (p *Parser) ifStatement() ast.Statement {
	pos := p.advance().Pos // 'if' keyword
	condition := p.expression(0)

	thenPos := p.expectToken(token.TokenLeftBrace).Pos
	then := p.parseBlock()

	var elseBranch *ast.BlockStatement = nil

	if p.match(token.TokenElseKw) {
		if p.check(token.TokenIfKw) {
			cond, _ := p.ifStatement().(ast.IfStatement) // this won't panic

			elseBranch = &ast.BlockStatement{
				AstBase: ast.AstBase{
					Pos: cond.Pos,
				},
				Stmts: []ast.Statement{
					cond,
				},
			}
		} else {
			elsePos := p.expectToken(token.TokenLeftBrace).Pos
			elseBlock := p.parseBlock()

			elseBranch = &ast.BlockStatement{
				AstBase: ast.AstBase{
					Pos: elsePos,
				},
				Stmts: elseBlock,
			}
		}
	}

	return ast.IfStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Condition: condition,
		Then: ast.BlockStatement{
			AstBase: ast.AstBase{
				Pos: thenPos,
			},
			Stmts: then,
		},
		Else: elseBranch,
	}
}

// TODO: use this function to disambiguate between for (each) and for (var) loops
func (p *Parser) forStatement() ast.Statement {
	pos := p.advance().Pos // 'for' keyword
	p.expectTokenNoAdvance(token.TokenVarKw) // for now, it is mandatory

	decl := p.varStatement() // already requires a semicolon
	condition := p.expression(PrecLowest)
	p.requireSemicolon()

	// increment can be omitted
	
}

func (p *Parser) whileStatement() ast.Statement {
	pos := p.advance().Pos // 'var' keyword
	condition := p.expression(PrecLowest)

	blockPos := p.expectToken(token.TokenLeftBrace).Pos
	block := p.parseBlock()

	return ast.WhileStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Condition: condition,
		Block: ast.BlockStatement{
			AstBase: ast.AstBase{
				Pos: blockPos,
			},
			Stmts: block,
		},
	}
}

func (p *Parser) varStatement() ast.VarStatement {
	pos := p.advance().Pos // 'var' keyword
	name := p.expectToken(token.TokenIdentifier)
	p.expect(token.TokenEqual)

	expr := p.expression(0)
	p.requireSemicolon()

	return ast.VarStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Name: name,
		Init: expr,
	}
}

func (p *Parser) printStatement() ast.Statement {
	pos := p.advance().Pos // 'print' keyword
	expr := p.expression(0)

	p.requireSemicolon()

	return ast.PrintStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Expr: expr,
	}
}

func (p *Parser) blockStatement() ast.BlockStatement {
	pos := p.advance().Pos // '{'
	stmts := p.parseBlock()

	return ast.BlockStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Stmts: stmts,
	}
}

func (p *Parser) exprStatement() ast.Statement {
	pos := p.peek(0).Pos
	expr := p.expression(PrecLowest)

	p.requireSemicolon()

	return ast.ExprStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},
		Expr: expr,
	}
}
