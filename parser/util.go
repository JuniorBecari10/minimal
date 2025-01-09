package parser

import (
	"fmt"
	"vm-go/ast"
	"vm-go/token"
	"vm-go/util"
)

func (p *Parser) parseBlock() ast.BlockStatement {
	p.expect(token.TokenLeftBrace)
	stmts := []ast.Statement{}

	for !p.isAtEnd(0) && !p.check(token.TokenRightBrace) {
		stmts = append(stmts, p.declaration(true))
	}

	p.expect(token.TokenRightBrace)
	return ast.BlockStatement{
		Stmts: stmts,
	}
}

func (p *Parser) parseMethods() []ast.FnStatement {
	p.expect(token.TokenLeftBrace)
	methods := []ast.FnStatement{}

	for !p.isAtEnd(0) && !p.check(token.TokenRightBrace) {
		methods = append(methods, p.fnStatement().Data.(ast.FnStatement))
	}

	p.expect(token.TokenRightBrace)
	return methods
}

func (p *Parser) parseParameters() []ast.Parameter {
	p.expect(token.TokenLeftParen)
	params := []ast.Parameter{}

	for !p.match(token.TokenRightParen) && !p.isAtEnd(0) && !p.panicMode {
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

func (p *Parser) parseFields() []ast.Field {
	params := p.parseParameters()
	fields := []ast.Field{}

	for _, param := range params {
		fields = append(fields, ast.Field(param))
	}

	return fields
}

func (p *Parser) expect(kind token.TokenKind) bool {
	return !p.expectToken(kind).IsAbsent()
}

func (p *Parser) expectToken(kind token.TokenKind) token.Token {
	if !p.check(kind) {
		if p.isAtEnd(0) {
			p.error(fmt.Sprintf("Expected '%s', but reached end.", kind))
		} else {
			p.error(fmt.Sprintf("Expected '%s', but got '%s' instead.", kind, p.peek(0).Kind))
		}
		return token.AbsentToken()
	}

	return p.advance()
}

func (p *Parser) expectTokenNoAdvance(kind token.TokenKind) token.Token {
	if !p.check(kind) {
		if p.isAtEnd(0) {
			p.error(fmt.Sprintf("Expected '%s', reached end.", kind))
		} else {
			p.error(fmt.Sprintf("Expected '%s', got '%s' instead.", kind, p.peek(0).Kind))
		}
		return token.AbsentToken()
	}

	return p.peek(0)
}

func (p *Parser) requireSemicolon() {
	if !p.check(token.TokenSemicolon) {
		if p.isAtEnd(0) {
			p.rawError("Expected ';' after statement, but reached end.", 1, token.Position{
				Line: p.peek(-1).Pos.Line,
				Col: p.peek(-1).Pos.Col + len(p.peek(-1).Lexeme),
			})
		} else {
			p.rawError(fmt.Sprintf("Expected ';' after statement, but got '%s' instead.", p.peek(0).Kind), 1, token.Position{
				Line: p.peek(-1).Pos.Line,
				Col: p.peek(-1).Pos.Col + len(p.peek(-1).Lexeme),
			})
		}
		return
	}

	p.advance()
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
        kind := p.peek(0).Kind

        // Return if a synchronization point is found
        switch kind {
			case token.TokenVarKw, token.TokenLeftBrace, token.TokenRightBrace,
				token.TokenIfKw, token.TokenElseKw, token.TokenWhileKw, token.TokenBreakKw, token.TokenContinueKw,
				token.TokenForKw, token.TokenFnKw, token.TokenReturnKw, token.TokenRecordKw,
				token.TokenSemicolon:
				return
        }

        // Advance to the next token if no synchronization point is found
        p.advance()
    }
}


func (p *Parser) error(message string) {
	last := p.peek(0)

	if last.IsAbsent() {
		last = p.tokens[len(p.tokens)-1]
	}

	pos := last.Pos

	p.rawError(message, len(last.Lexeme), pos)
}

func (p *Parser) rawError(message string, len int, pos token.Position) {
	if p.panicMode {
		return
	}
	
	util.Error(pos, len, message, p.fileData)

	p.hadError = true
	p.panicMode = true
}
