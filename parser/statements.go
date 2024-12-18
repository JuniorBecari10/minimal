package parser

import (
	"vm-go/ast"
	"vm-go/token"
)

func (p *Parser) declaration(allowStatements bool) ast.Statement {
	if p.panicMode {
		p.synchronize()
	}

	t := p.peek(0)
	
	switch t.Kind {
		case token.TokenRecordKw: return p.recordStatement()
		case token.TokenFnKw: return p.fnStatement()
		case token.TokenVarKw: return p.varStatement()
		
		default: {
			if allowStatements {
				return p.statement()
			} else {
				p.error("Statements are not allowed at top-level.")
				p.advance()
				return nil
			}
		}
	}
}

func (p *Parser) statement() ast.Statement {
	t := p.peek(0)
	
	switch t.Kind {
		case token.TokenIfKw: return p.ifStatement()
		case token.TokenWhileKw: return p.whileStatement()
		case token.TokenForKw: return p.forStatement()
		case token.TokenBreakKw: return p.breakStatement()
		case token.TokenContinueKw: return p.continueStatement()
		case token.TokenReturnKw: return p.returnStatement()
		case token.TokenLeftBrace: return p.blockStatement()
		
		default: return p.exprStatement()
	}
}

// ---

func (p *Parser) recordStatement() ast.Statement {
	pos := p.advance().Pos // 'record' keyword
	name := p.expectToken(token.TokenIdentifier)
	fields := p.parseFields()

	p.requireSemicolon()

	return ast.RecordStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},
		Name: name,
		Fields: fields,
	}
}

func (p *Parser) fnStatement() ast.Statement {
	pos := p.advance().Pos // 'fn' keyword

	name := p.expectToken(token.TokenIdentifier)
	parameters := p.parseParameters()
	body := p.parseBlock()

	return ast.FnStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},
		Name: name,
		Parameters: parameters,
		Body: body,
	}
}

func (p *Parser) returnStatement() ast.Statement {
	pos := p.advance().Pos // 'return' keyword
	var expr *ast.Expression = nil

	if !p.check(token.TokenSemicolon) {
		expr_ := p.expression(PrecLowest)
		expr = &expr_
	}

	p.requireSemicolon()
	return ast.ReturnStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},
		Expression: expr,
	}
}

func (p *Parser) ifStatement() ast.Statement {
	pos := p.advance().Pos // 'if' keyword

	condition := p.expression(0)
	then := p.parseBlock()

	var else_ *ast.BlockStatement = nil

	if p.match(token.TokenElseKw) {
		if p.check(token.TokenIfKw) {
			cond, _ := p.ifStatement().(ast.IfStatement) // this won't panic

			else_ = &ast.BlockStatement{
				AstBase: ast.AstBase{
					Pos: cond.Pos,
				},
				Stmts: []ast.Statement{
					cond,
				},
			}
		} else {
			elseBlock := p.parseBlock()
			else_ = &elseBlock
		}
	}

	return ast.IfStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Condition: condition,
		Then: then,
		Else: else_,
	}
}

// TODO: use this function to disambiguate between for (each) and for (var) loops
func (p *Parser) forStatement() ast.Statement {
	pos := p.advance().Pos // 'for' keyword
	p.expectTokenNoAdvance(token.TokenVarKw) // for now, it is mandatory

	declaration := p.varStatement() // already requires a semicolon
	condition := p.expression(PrecLowest)
	var increment *ast.Expression = nil
	
	if p.match(token.TokenSemicolon) {
		expr := p.expression(PrecLowest)
		increment = &expr // it can escape the scope - the GC collects it later
	}

	block := p.parseBlock()

	return ast.ForVarStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Declaration: declaration,
		Condition: condition,
		Increment: increment,
		
		Block: block,
	}
}

func (p *Parser) whileStatement() ast.Statement {
	pos := p.advance().Pos // 'while' keyword
	condition := p.expression(PrecLowest)
	block := p.parseBlock()

	return ast.WhileStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},

		Condition: condition,
		Block: block,
	}
}

func (p *Parser) breakStatement() ast.Statement {
	token := p.advance() // 'break' keyword
	p.requireSemicolon()

	return ast.BreakStatement{
		AstBase: ast.AstBase{
			Pos: token.Pos,
		},
		Token: token,
	}
}

func (p *Parser) continueStatement() ast.Statement {
	token := p.advance() // 'continue' keyword
	p.requireSemicolon()

	return ast.ContinueStatement{
		AstBase: ast.AstBase{
			Pos: token.Pos,
		},
		Token: token,
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

func (p *Parser) blockStatement() ast.BlockStatement {
	return p.parseBlock()
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
