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
				return ast.Statement{}
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
		case token.TokenLoopKw: return p.loopStatement()
		case token.TokenBreakKw: return p.breakStatement()
		case token.TokenContinueKw: return p.continueStatement()
		case token.TokenReturnKw: return p.returnStatement()
		case token.TokenLeftBrace: return p.blockStatement()
		
		default: return p.exprStatement()
	}
}

// ---

func (p *Parser) recordStatement() ast.Statement {
	keyword := p.advance()
	name := p.expectToken(token.TokenIdentifier)
	fields := p.parseFields()

	p.requireSemicolon()

	return ast.Statement{
		Base: ast.AstBase{
			Pos:    keyword.Pos,
			Length: len(keyword.Lexeme),
		},

		Data: ast.RecordStatement{
			Name:   name,
			Fields: fields,
		},
	}
}

func (p *Parser) fnStatement() ast.Statement {
	keyword := p.advance()

	name := p.expectToken(token.TokenIdentifier)
	parameters := p.parseParameters()
	body := p.parseBlock()

	return ast.Statement{
		Base: ast.AstBase{
			Pos: keyword.Pos,
			Length: len(keyword.Lexeme),
		},
		Data: ast.FnStatement{
			Name: name,
			Parameters: parameters,
			Body: body,
		},
	}
}

func (p *Parser) returnStatement() ast.Statement {
	keyword := p.advance()
	var expr *ast.Expression = nil

	if !p.check(token.TokenSemicolon) {
		expr_ := p.expression(PrecLowest)
		expr = &expr_
	}

	p.requireSemicolon()
	return ast.Statement{
		Base: ast.AstBase{
			Pos: keyword.Pos,
			Length: len(keyword.Lexeme),
		},
		Data: ast.ReturnStatement{
			Expression: expr,
		},
	}
}

func (p *Parser) ifStatement() ast.Statement {
	keyword := p.advance()

	condition := p.expression(0)
	then := p.parseBlock()

	var else_ *ast.BlockStatement = nil

	if p.match(token.TokenElseKw) {
		if p.check(token.TokenIfKw) {
			elseIf := p.ifStatement()
			cond, _ := elseIf.Data.(ast.IfStatement) // this won't panic

			else_ = &ast.BlockStatement{
				Stmts: []ast.Statement{
					{
						Base: ast.AstBase{
							Pos: elseIf.Base.Pos,
						},

						Data: cond,
					},
				},
			}
		} else {
			elseBlock := p.parseBlock()
			else_ = &elseBlock
		}
	}

	return ast.Statement{
		Base: ast.AstBase{
			Pos: keyword.Pos,
			Length: len(keyword.Lexeme),
		},

		Data: ast.IfStatement{
			Condition: condition,
			Then: then,
			Else: else_,
		},
	}
}

// TODO: use this function to disambiguate between for (each) and for (var) loops
func (p *Parser) forStatement() ast.Statement {
	keyword := p.advance()
	p.expectTokenNoAdvance(token.TokenVarKw) // for now, it is mandatory
 	
	// already requires a semicolon
	declaration := p.varStatement()
	condition := p.expression(PrecLowest)

	var increment *ast.Expression
	
	if p.match(token.TokenSemicolon) {
		expr := p.expression(PrecLowest)
		increment = &expr // it can escape the scope - the GC collects it later
	}

	block := p.parseBlock()

	return ast.Statement{
		Base: ast.AstBase{
			Pos: keyword.Pos,
			Length: len(keyword.Lexeme),
		},
		Data: ast.ForVarStatement{
			Declaration: declaration, // won't panic
			Condition: condition,
			Increment: increment,
			Block: block,
		},
	}
}

func (p *Parser) whileStatement() ast.Statement {
	keyword := p.advance()
	condition := p.expression(PrecLowest)
	block := p.parseBlock()

	return ast.Statement{
		Base: ast.AstBase{
			Pos: keyword.Pos,
			Length: len(keyword.Lexeme),
		},
		Data: ast.WhileStatement{
			Condition: condition,
			Block: block,
		},
	}
}

func (p *Parser) loopStatement() ast.Statement {
	keyword := p.advance()
	block := p.parseBlock()

	return ast.Statement{
		Base: ast.AstBase{
			Pos: keyword.Pos,
			Length: len(keyword.Lexeme),
		},
		Data: ast.LoopStatement{
			Block: block,
		},
	}
}

func (p *Parser) breakStatement() ast.Statement {
	keyword := p.advance()
	p.requireSemicolon()

	return ast.Statement{
		Base: ast.AstBase{
			Pos: keyword.Pos,
			Length: len(keyword.Lexeme),
		},
		Data: ast.BreakStatement{
			Token: keyword,
		},
	}
}

func (p *Parser) continueStatement() ast.Statement {
	keyword := p.advance()
	p.requireSemicolon()

	return ast.Statement{
		Base: ast.AstBase{
			Pos: keyword.Pos,
			Length: len(keyword.Lexeme),
		},
		Data: ast.ContinueStatement{
			Token: keyword,
		},
	}
}

func (p *Parser) varStatement() ast.Statement {
	keyword := p.advance()
	name := p.expectToken(token.TokenIdentifier)
	p.expect(token.TokenEqual)

	expr := p.expression(0)
	p.requireSemicolon()

	return ast.Statement{
		Base: ast.AstBase{
			Pos: keyword.Pos,
			Length: len(keyword.Lexeme),
		},
		Data: ast.VarStatement{
			Name: name,
			Init: expr,
		},
	}
}

func (p *Parser) blockStatement() ast.Statement {
	pos := p.peek(0).Pos

	return ast.Statement{
		Base: ast.AstBase{
			Pos: pos,
			Length: 1,
		},
		Data: p.parseBlock(),
	}
}

func (p *Parser) exprStatement() ast.Statement {
	pos := p.peek(0).Pos
	expr := p.expression(PrecLowest)

	p.requireSemicolon()

	return ast.Statement{
		Base: ast.AstBase{
			Pos: pos,
			Length: expr.Base.Length,
		},
		Data: ast.ExprStatement{
			Expr: expr,
		},
	}
}
