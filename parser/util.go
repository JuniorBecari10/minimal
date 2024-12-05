package parser

import (
	"fmt"
	"vm-go/ast"
	"vm-go/token"
	"vm-go/util"
)

func (p *Parser) parseBlock() ast.BlockStatement {
	pos := p.expectToken(token.TokenLeftBrace).Pos
	stmts := []ast.Statement{}

	for !p.isAtEnd(0) && !p.check(token.TokenRightBrace) && !p.hadError {
		stmts = append(stmts, p.statement())
	}

	p.expect(token.TokenRightBrace)
	return ast.BlockStatement{
		AstBase: ast.AstBase{
			Pos: pos,
		},
		Stmts: stmts,
	}
}

func (p *Parser) parseParameters() []ast.Parameter {
	p.expect(token.TokenLeftParen)
	params := []ast.Parameter{}

	for !p.match(token.TokenRightParen) {
		name := p.expectToken(token.TokenIdentifier)
		params = append(params, ast.Parameter{
			Name: name,
		})

		if !p.check(token.TokenRightParen) {
			p.expect(token.TokenComma)
		}
	}

	return params
}

func (p *Parser) expect(kind token.TokenKind) bool {
	return !p.expectToken(kind).IsAbsent()
}

func (p *Parser) expectToken(kind token.TokenKind) token.Token {
	if !p.check(kind) {
		if p.isAtEnd(0) {
			p.error(fmt.Sprintf("Expected '%s', reached end", kind))
		} else {
			p.error(fmt.Sprintf("Expected '%s', got '%s'", kind, p.peek(0).Kind))
		}
		return token.AbsentToken()
	}

	return p.advance()
}

func (p *Parser) expectTokenNoAdvance(kind token.TokenKind) token.Token {
	if !p.check(kind) {
		if p.isAtEnd(0) {
			p.error(fmt.Sprintf("Expected '%s', reached end", kind))
		} else {
			p.error(fmt.Sprintf("Expected '%s', got '%s'", kind, p.peek(0).Kind))
		}
		return token.AbsentToken()
	}

	return p.peek(0)
}

func (p *Parser) requireSemicolon() {
	p.expect(token.TokenSemicolon)
}

func (p *Parser) check(kind token.TokenKind) bool {
	return p.peek(0).Kind == kind
}

// it advances
func (p *Parser) match(kind token.TokenKind) bool {
	if p.check(kind) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) advance() token.Token {
	peek := p.peek(0)
	p.current += 1

	return peek
}

func (p *Parser) peek(offset int) token.Token {
	if p.isAtEnd(offset) {
		return token.AbsentToken()
	}

	return p.tokens[p.current + offset]
}

func (p *Parser) isAtEnd(offset int) bool {
	return p.current + offset >= len(p.tokens)
}

func (p *Parser) synchronize() {
	p.panicMode = false

	for !p.isAtEnd(0) {
		switch p.peek(0).Kind {
		case token.TokenVarKw, token.TokenLeftBrace, token.TokenIfKw:
			return
		}

		if p.peek(0).Kind == token.TokenSemicolon {
			return
		}

		p.advance()
	}
}

func (p *Parser) error(message string) {
	if p.panicMode {
		return
	}

	// TODO: if reached end, get the position of the last token
	last := p.peek(-1)
	pos := last.Pos

	pos.Col += len(last.Lexeme)
	util.Error(pos, len(last.Lexeme), message, p.fileData)

	p.hadError = true
	p.panicMode = true
}
