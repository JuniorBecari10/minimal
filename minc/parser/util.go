package parser

import (
	"fmt"
	"minc/ast"
	"minlib/token"
	"minlib/util"
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

func (p *Parser) parseExpression() ast.Expression {
	return p.expression(PrecLowest)
}

func (p *Parser) expect(kind token.TokenKind) bool {
	return !p.expectToken(kind).IsEnd()
}

func (p *Parser) expectToken(kind token.TokenKind) token.Token {
	if !p.check(kind) {
		if p.isAtEnd(0) {
			p.error(fmt.Sprintf("Expected %s after %s, but reached end.", kind, p.peek(-1).FormatError()))
		} else {
			p.error(fmt.Sprintf("Expected %s after %s, but found %s instead.", kind, p.peek(-1).FormatError(), p.peek(0).FormatError()))
		}
		return token.EndToken()
	}

	return p.advance()
}

func (p *Parser) requireSemicolon() {
	if !p.check(token.TokenSemicolon) {
		if p.isAtEnd(0) {
			p.printErrorHelp("Expected ';' after statement, but reached end.",
				"Insert a semicolon where the arrow is pointing.", 1, token.Position{
					Line: p.peek(-1).Pos.Line,
					Col: p.peek(-1).Pos.Col + uint32(len(p.peek(-1).Lexeme)),
				})
		} else {
			p.printErrorHelp(
				fmt.Sprintf("Expected ';' after statement, but found %s instead.",
					p.peek(0).FormatError()),
					"Insert a semicolon where the arrow is pointing.", 1, token.Position{
						Line: p.peek(-1).Pos.Line,
						Col: p.peek(-1).Pos.Col + uint32(len(p.peek(-1).Lexeme)),
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

func (p *Parser) makeAssignment(left ast.Expression, right ast.Expression, operator token.Token) ast.Expression {
	switch lValue := left.Data.(type) {
	case ast.IdentifierExpression:
		return ast.Expression{
			Base: ast.AstBase{
				Pos:    operator.Pos,
				Length: len(operator.Lexeme),
			},
			Data: ast.IdentifierAssignmentExpression{
				Name: lValue.Token,
				Expr: right,
			},
		}

	case ast.GetPropertyExpression:
		return ast.Expression{
			Base: ast.AstBase{
				Pos:  	lValue.Property.Pos,
				Length: len(lValue.Property.Lexeme),
			},
			Data: ast.SetPropertyExpression{
				Left:    lValue.Left,
				Property: lValue.Property,
				Value:   right,
			},
		}

	default:
		p.error(fmt.Sprintf("Invalid assignment target: '%v'.", left))
		return ast.Expression{}
	}
}

func (p *Parser) peek(offset int) token.Token {
	if p.isAtStart(offset) {
		return token.StartToken()
	} else if p.isAtEnd(offset) {
		return token.EndToken()
	}

	return p.tokens[p.current + offset]
}

func (p *Parser) isAtStart(offset int) bool {
	return p.current + offset < 0
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

	if last.IsEnd() {
		last = p.tokens[len(p.tokens)-1]
	}

	pos := last.Pos
	p.printError(message, len(last.Lexeme), pos)
}

func (p *Parser) reportError(pos token.Position, length int, message string, help *string) {
	if p.panicMode {
		return
	}

	util.PrintError(pos, length, message, help, p.fileData)

	p.hadError = true
	p.panicMode = true
}

func (p *Parser) printError(message string, length int, pos token.Position) {
	p.reportError(pos, length, message, nil)
}

func (p *Parser) printErrorHelp(message, help string, length int, pos token.Position) {
	p.reportError(pos, length, message, &help)
}
