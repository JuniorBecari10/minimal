package parser

import (
	"fmt"
	"minc/ast"
	"minlib/token"
	"strconv"
)

func (p *Parser) expression(precedence int) ast.Expression {
	peek := p.peek(0)

	pos := peek.Pos
	prefixFn, ok := p.prefixMap[peek.Kind]

	if !ok {
		p.error(fmt.Sprintf("Expected expression after %s, but found %s instead.",
			p.peek(-1).FormatError(),
			p.peek(0).FormatError()))

		return ast.Expression{}
	}

	left := prefixFn()

	for p.precedenceMap[p.peek(0).Kind] > precedence {
		infixFn, ok := p.infixMap[p.peek(0).Kind]

		if !ok {
			break
		}

		left = infixFn(left, pos)
	}

	return left
}

// ---

func (p *Parser) parseInt() ast.Expression {
	tok := p.advance()
	value, _ := strconv.Atoi(tok.Lexeme)

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme),
		},
		Data: ast.IntExpression{
			Literal: int32(value),
		},
	}
}

func (p *Parser) parseFloat() ast.Expression {
	tok := p.advance()
	value, _ := strconv.ParseFloat(tok.Lexeme, 64)

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme),
		},
		Data: ast.FloatExpression{
			Literal: value,
		},
	}
}

func (p *Parser) parseString() ast.Expression {
	tok := p.advance()

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme) + 2, // the quotes
		},
		Data: ast.StringExpression{
			Literal: tok.Lexeme,
		},
	}
}

func (p *Parser) parseChar() ast.Expression {
	tok := p.advance()

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme) + 2,
		},
		Data: ast.CharExpression{
			Literal: uint8(tok.Lexeme[0]), // guaranteed to be one character long
		},
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := p.expectToken(token.TokenIdentifier)

	return ast.Expression{
		Base: ast.AstBase{
			Pos: ident.Pos,
			Length: len(ident.Lexeme),
		},
		Data: ast.IdentifierExpression{
			Token: ident,
		},
	}
}

func (p *Parser) parseSelf() ast.Expression {
	self := p.expectToken(token.TokenSelfKw)

	return ast.Expression{
		Base: ast.AstBase{
			Pos: self.Pos,
			Length: len(self.Lexeme),
		},
		Data: ast.SelfExpression{
			Token: self,
		},
	}
}

func (p *Parser) parseBool() ast.Expression {
	tok := p.advance()

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme),
		},
		Data: ast.BoolExpression{
			Literal: tok.Kind == token.TokenTrueKw,
		},
	}
}

func (p *Parser) parseNil() ast.Expression {
	tok := p.advance()

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme),
		},
		Data: ast.NilExpression{},
	}
}

func (p *Parser) parseVoid() ast.Expression {
	tok := p.advance()

	var expr *ast.Expression = nil

	if p.match(token.TokenLeftParen) {
		e := p.parseExpression()

		p.expect(token.TokenRightParen)
		expr = &e
	}

	return ast.Expression{
		Base: ast.AstBase{
			Pos: tok.Pos,
			Length: len(tok.Lexeme),
		},
		Data: ast.VoidExpression{
			Expr: expr,
		},
	}
}

// disambuguation between parenthesized expression and lambda expression.
func (p *Parser) lParen() ast.Expression {
	// lambda: ')' | ( ident ',' | ')' )
	if p.peek(1).Kind == token.TokenRightParen ||
		(p.peek(1).Kind == token.TokenIdentifier && (p.peek(2).Kind == token.TokenRightParen || p.peek(2).Kind == token.TokenComma)) {
		return p.parseLambda()
	} else {
		return p.parseGroup()
	}
}

func (p *Parser) parseLambda() ast.Expression {
	pos := p.peek(0).Pos
	parameters := p.parseParameters()

	p.expect(token.TokenArrow)

	var body ast.BlockStatement

	if p.peek(0).Kind == token.TokenLeftBrace {
		body = p.parseBlock()
	} else {
		peek := p.peek(0)
		expr := p.parseExpression()

		body = ast.BlockStatement{
			Stmts: []ast.Statement{
				{
					Base: ast.AstBase{
						Pos: peek.Pos,
						Length: len(peek.Lexeme),
					},
					Data: ast.ReturnStatement{
						Expression: &expr,
					},
				},
			},
		}
	}

	return ast.Expression{
		Base: ast.AstBase{
			Pos: pos,
			Length: 1, // TODO: change
		},
		Data: ast.FnExpression{
			Parameters: parameters,
			Body: body,
		},
	}
}

func (p *Parser) parseIfExpr() ast.Expression {
	if_ := p.advance()
	cond := p.parseExpression()

	p.expect(token.TokenColon)
	then := p.parseExpression()

	p.expect(token.TokenElseKw)
	var else_ ast.Expression

	// Check 'else if' chain and not require the colon if so.
	if p.check(token.TokenIfKw) {
		else_ = p.parseIfExpr()
	} else {
		p.expect(token.TokenColon)
		else_ = p.parseExpression()
	}

	return ast.Expression{
		Base: ast.AstBase{
			Pos: if_.Pos,
			Length: len(if_.Lexeme),
		},
		Data: ast.IfExpression{
			Condition: cond,
			Then: then,
			Else: else_,
		},
	}
}

func (p *Parser) parseGroup() ast.Expression {
	pos := p.peek(0).Pos
	p.expect(token.TokenLeftParen)

	expr := p.parseExpression()
	p.expect(token.TokenRightParen)

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    pos,
			Length: expr.Base.Length + 2, // The parentheses
		},
		Data: ast.GroupExpression{
			Expr: expr,
		},
	}
}

func (p *Parser) parseUnary(op token.TokenKind) ast.Expression {
	pos := p.peek(0).Pos

	operator := p.expectToken(op)
	operand := p.expression(PrecUnary)

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    pos,
			Length: len(operator.Lexeme) + operand.Base.Length,
		},
		Data: ast.UnaryExpression{
			Operand:  operand,
			Operator: operator,
		},
	}
}

// --- Infix ---

func (p *Parser) parseBinary(left ast.Expression, op token.TokenKind) ast.Expression {
	precedence := p.precedenceMap[op]

	operator := p.expectToken(op)
	right := p.expression(precedence)

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    operator.Pos,
			Length: len(operator.Lexeme),
		},
		Data: ast.BinaryExpression{
			Left:     left,
			Right:    right,
			Operator: operator,
		},
	}
}

func (p *Parser) parseLogical(left ast.Expression, op token.TokenKind) ast.Expression {
	precedence := p.precedenceMap[op]

	operator := p.expectToken(op)
	shortCircuit := !p.match(token.TokenStar)

	right := p.expression(precedence)

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    operator.Pos,
			Length: len(operator.Lexeme),
		},
		Data: ast.LogicalExpression{
			Left:         left,
			Right:        right,
			Operator:     operator,
			ShortCircuit: shortCircuit,
		},
	}
}

func (p *Parser) parseCall(left ast.Expression, pos token.Position) ast.Expression {
	p.expectToken(token.TokenLeftParen)
	arguments := []ast.Expression{}

	for !p.match(token.TokenRightParen) {
		arguments = append(arguments, p.parseExpression())

		if !p.check(token.TokenRightParen) {
			p.expect(token.TokenComma)
		}
	}

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    pos,
			Length: left.Base.Length, // the callee
		},
		Data: ast.CallExpression{
			Callee:    left,
			Arguments: arguments,
		},
	}
}

func (p *Parser) parseAssignment(left ast.Expression, pos token.Position) ast.Expression {
	operator := p.expectToken(token.TokenEqual)
	right := p.parseExpression() // accept one level higher because assignment is right-associative

	return p.makeAssignment(left, right, operator)
}

// this operator will act as a macro for the default assignment, which generates
// the same AST nodes if you repeat the l-value and put it into a binary expression.
// examples:
// v += 10; -> v = v + 10;
// w -= 20; -> w = w - 20;
// x *= 30; -> x = x * 30;
// y /= 40; -> y = y / 40;
// z %= 50; -> z = z % 50;
func (p *Parser) parseOperatorAssignment(left ast.Expression, finalOp token.TokenKind) ast.Expression {
	operator := p.advance()

	// this will convert the assignment operator into the correspondent binary operator.
	// examples:
	// '+=' -> '+'
	// '-=' -> '-'
	// ...
	operatorAsFinal := operator
	operatorAsFinal.Kind = finalOp
	operatorAsFinal.Lexeme = string(finalOp)

	// this will convert the right-hand expression into a binary expression that contains the left-hand side
	// of the assignment into the left-hand of the binary expression.
	// examples:
	// x += 10 -> x = x + 10
	// y -= 20 -> y = y - 20
	right := p.parseExpression()
	right.Data = ast.BinaryExpression{
		Left: left,
		Right: right,
		Operator: operatorAsFinal,
	}

	return p.makeAssignment(left, right, operator)
}

func (p *Parser) parseDot(left ast.Expression, pos token.Position) ast.Expression {
	p.expect(token.TokenDot)
	property := p.expectToken(token.TokenIdentifier)

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    property.Pos,
			Length: len(property.Lexeme),
		},
		Data: ast.GetPropertyExpression{
			Left:    left,
			Property: property,
		},
	}
}

func (p *Parser) parseRange(left ast.Expression, pos token.Position) ast.Expression {
	operator := p.expectToken(token.TokenDoubleDot)

    inclusive := p.match(token.TokenEqual)
	right := p.expression(PrecRange)

	var step *ast.Expression = nil

	if p.match(token.TokenColon) {
		expr := p.expression(PrecRange)
		step = &expr
	}

	return ast.Expression{
		Base: ast.AstBase{
			Pos:    operator.Pos,
			Length: len(operator.Lexeme),
		},
		Data: ast.RangeExpression{
			Start: left,
			End: right,
			Step: step,

            Inclusive: inclusive,
		},
	}
}
