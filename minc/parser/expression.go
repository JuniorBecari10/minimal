package parser

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/types"
	"minlib/token"
	"strconv"
)

type PrefixFn = func() (ast.Expression, diagnostic.Diagnostic)
type InfixFn = func(ast.Expression, token.Position) (ast.Expression, diagnostic.Diagnostic)
type Precedence int

const (
	PREC_LOWEST Precedence = iota
	PREC_ASSIGNMENT    // = += -= *= /= %=
	PREC_OR            // or
	PREC_AND           // and
	PREC_EQUAL         // == !=
	PREC_COMPARISON    // < > <= >=
	PREC_TERM          // + -
	PREC_FACTOR        // * /
	PREC_UNARY         // not -
	PREC_AS            // as
	PREC_CALL          // ()
	PREC_GET_PROPERTY  // .
	PREC_RANGE         // ..
	// PREC_PRIMARY    // literals, identifiers (unused)
)

func (p *Parser) parseExpression() (ast.Expression, diagnostic.Diagnostic) {
	return p.expression(PREC_LOWEST)
}

func (p *Parser) expression(prec Precedence) (ast.Expression, diagnostic.Diagnostic) {
	leftPos := p.current.Pos
	prefixFn := p.getPrefixFn(p.current.Kind)

	if prefixFn == nil {
		return ast.Expression{}, p.makeExpectedExpressionDiagnostic()
	}

	left, diag := prefixFn(); if diag != nil {
		return ast.Expression{}, diag
	}

	for getPrecedenceFor(p.current.Kind) > prec {
		infixFn := p.getInfixFn(p.current.Kind)

		if infixFn == nil {
			break
		}

		left, diag = infixFn(left, leftPos); if diag != nil {
			return ast.Expression{}, diag
		}
	}

	return left, nil
}

func (p *Parser) getPrefixFn(kind token.TokenKind) PrefixFn {
	switch kind {
		case token.TokenIntLiteral: return p.parseInt
		case token.TokenFloatLiteral: return p.parseFloat
		case token.TokenStringLiteral: return p.parseStr
		case token.TokenCharLiteral: return p.parseChar
		case token.TokenIdentifier: return p.parseIdentifier
		case token.TokenSelfKw: return p.parseSelf

		case token.TokenTrueKw: return p.parseBool
		case token.TokenFalseKw: return p.parseBool
		case token.TokenNilKw: return p.parseNil
		case token.TokenVoidKw: return p.parseVoid

		case token.TokenLeftParen: return p.parseGroup
		case token.TokenLeftBrace: return p.parseBlockExpr

		case token.TokenIfKw: return p.parseIf
		case token.TokenFnKw: return p.parseFnExpr

		case token.TokenNotKw: return func() (ast.Expression, diagnostic.Diagnostic) { return p.parseUnary(token.TokenNotKw) }
		case token.TokenMinus: return func() (ast.Expression, diagnostic.Diagnostic) { return p.parseUnary(token.TokenMinus) }

		default: return nil
	}
}

func (p *Parser) getInfixFn(kind token.TokenKind) InfixFn {
	switch kind {
		case token.TokenPlus: return p.infixBinary(token.TokenPlus)
		case token.TokenMinus: return p.infixBinary(token.TokenMinus)

		case token.TokenStar: return p.infixBinary(token.TokenStar)
		case token.TokenSlash: return p.infixBinary(token.TokenSlash)
		case token.TokenPercent: return p.infixBinary(token.TokenPercent)

		case token.TokenPlusEqual: return p.infixOpAssign(token.TokenPlus)
		case token.TokenMinusEqual: return p.infixOpAssign(token.TokenMinus)
		case token.TokenStarEqual: return p.infixOpAssign(token.TokenStar)
		case token.TokenSlashEqual: return p.infixOpAssign(token.TokenSlash)
		case token.TokenPercentEqual: return p.infixOpAssign(token.TokenPercent)

		case token.TokenOrKw: return p.infixLogical(token.TokenOrKw)
		case token.TokenAndKw: return p.infixLogical(token.TokenAndKw)

		case token.TokenGreater: return p.infixBinary(token.TokenGreater)
		case token.TokenGreaterEqual: return p.infixBinary(token.TokenGreaterEqual)
		case token.TokenLess: return p.infixBinary(token.TokenLess)
		case token.TokenLessEqual: return p.infixBinary(token.TokenLessEqual)

		case token.TokenDoubleEqual: return p.infixBinary(token.TokenDoubleEqual)
		case token.TokenBangEqual: return p.infixBinary(token.TokenBangEqual)

		case token.TokenEqual: return p.parseAssignment
		case token.TokenLeftParen: return p.parseCall
		case token.TokenDot: return p.parseDot
		case token.TokenDoubleDot: return p.parseRange
		
		case token.TokenAsKw: return p.parseAs

		default: return nil
	}
}

func getPrecedenceFor(kind token.TokenKind) Precedence {
	switch kind {
		case token.TokenPlus, token.TokenMinus:
			return PREC_TERM

		case token.TokenStar, token.TokenSlash, token.TokenPercent:
			return PREC_FACTOR

		case token.TokenPlusEqual, token.TokenMinusEqual,
			token.TokenStarEqual, token.TokenSlashEqual, token.TokenPercentEqual,
			token.TokenEqual:
			return PREC_ASSIGNMENT

		case token.TokenAndKw:
			return PREC_AND

		case token.TokenOrKw:
			return PREC_OR

		case token.TokenGreater, token.TokenGreaterEqual,
			token.TokenLess, token.TokenLessEqual:
			return PREC_COMPARISON

		case token.TokenDoubleEqual, token.TokenBangEqual:
			return PREC_EQUAL

		case token.TokenLeftParen:
			return PREC_CALL

		case token.TokenDot:
			return PREC_GET_PROPERTY

		case token.TokenDoubleDot:
			return PREC_RANGE
		
		case token.TokenAsKw:
			return PREC_AS

		default:
			return PREC_LOWEST
	}
}

// ---

func (p *Parser) infixBinary(op token.TokenKind) InfixFn {
	return func(left ast.Expression, _ token.Position) (ast.Expression, diagnostic.Diagnostic) {
		return p.parseBinary(left, op)
	}
}

func (p *Parser) infixOpAssign(op token.TokenKind) InfixFn {
	return func(left ast.Expression, _ token.Position) (ast.Expression, diagnostic.Diagnostic) {
		return p.parseOperatorAssignment(left, op)
	}
}

func (p *Parser) infixLogical(op token.TokenKind) InfixFn {
	return func(left ast.Expression, _ token.Position) (ast.Expression, diagnostic.Diagnostic) {
		return p.parseLogical(left, op)
	}
}

// ---

func (p *Parser) parseInt() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()
	value, _ := strconv.Atoi(tok.Lexeme)

	return newExpr(tok, ast.IntExpression{
		Literal: int32(value),
	}), nil
}

func (p *Parser) parseFloat() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()
	value, _ := strconv.ParseFloat(tok.Lexeme, 64)

	return newExpr(tok, ast.FloatExpression{
		Literal: value,
	}), nil
}

func (p *Parser) parseStr() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()

	return newExpr(tok, ast.StringExpression{
		Literal: tok.Lexeme,
	}), nil
}

func (p *Parser) parseChar() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()

	return newExpr(tok, ast.CharExpression{
		Literal: uint8(tok.Lexeme[0]), // guaranteed to be one character long
	}), nil
}

func (p *Parser) parseIdentifier() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()

	return newExpr(tok, ast.IdentifierExpression{
		Token: tok,
	}), nil
}

func (p *Parser) parseSelf() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()
	return newExpr(tok, ast.SelfExpression{}), nil
}

func (p *Parser) parseBool() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()

	return newExpr(tok, ast.BoolExpression{
		Literal: tok.Kind == token.TokenTrueKw,
	}), nil
}

func (p *Parser) parseNil() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()
	var typeArgs []types.Type

	if p.check(token.TokenLess) {
		var diag diagnostic.Diagnostic

		typeArgs, diag = p.parseTypeArguments(); if diag != nil {
			return ast.Expression{}, diag
		}
	}

	return newExpr(tok, ast.NilExpression{
		TypeArguments: typeArgs,
		Token: tok,
	}), nil
}

func (p *Parser) parseVoid() (ast.Expression, diagnostic.Diagnostic) {
	tok, _ := p.advance()

	var expr *ast.Expression = nil

	if p.match(token.TokenLeftParen) {
		e, diag := p.parseExpression(); if diag != nil {
			return ast.Expression{}, diag
		}

		_, diag = p.expectToken(token.TokenRightParen); if diag != nil {
			return ast.Expression{}, diag
		}

		expr = &e
	}

	return newExpr(tok, ast.VoidExpression{
		Expr: expr,
	}), nil
}

func (p *Parser) parseGroup() (ast.Expression, diagnostic.Diagnostic) {
	pos := p.current.Pos

	_, diag := p.expectToken(token.TokenLeftParen); if diag != nil {
		return ast.Expression{}, diag
	}

	expr, diag := p.parseExpression(); if diag != nil {
		return ast.Expression{}, diag
	}

	_, diag = p.expectToken(token.TokenRightParen); if diag != nil {
		return ast.Expression{}, diag
	}

	return newExprToken(pos, expr.Base.Token, ast.GroupExpression{ // add 2 to the length?
		Expr: expr,
	}), nil
}

func (p *Parser) parseBlockExpr() (ast.Expression, diagnostic.Diagnostic) {
	tok := p.current

	block, diag := p.parseBlock(); if diag != nil {
		return ast.Expression{}, diag
	}

	return newExpr(tok, block), nil
}

func (p *Parser) parseIf() (ast.Expression, diagnostic.Diagnostic) {
	keyword, _ := p.advance()

	condition, diag := p.parseExpression(); if diag != nil {
		return ast.Expression{}, diag
	}

	thenToken := p.current
	then, diag := p.parseBlock(); if diag != nil {
		return ast.Expression{}, diag
	}

	elseToken := p.current
	var else_ *ast.BlockExpression = nil

	// Spécial case for 'else if', which is another 'if' inside the 'else' clause,
	// that does not need to be embraced in a block.
	if p.match(token.TokenElseKw) {
		if p.check(token.TokenIfKw) {
			elseIfExpr, diag := p.parseIf(); if diag != nil {
				return ast.Expression{}, diag
			}

			// Create a new block with an ExprStatement inside, which contains the 'if' expression.
			else_ = &ast.BlockExpression{
				Stmts: []ast.Statement{
					{
						Base: elseIfExpr.Base,

						Data: ast.ExprStatement{
							Expr: elseIfExpr,
						},
					},
				},
				Token: elseIfExpr.Base.Token,
			}
		} else {
			elseBlock, diag := p.parseBlock(); if diag != nil {
				return ast.Expression{}, diag
			}

			else_ = &elseBlock
		}
	}

	thenExpr := newExpr(thenToken, then)
	var elseExpr *ast.Expression = nil

	if else_ != nil {
		elseExprDecl := newExpr(elseToken, else_)
		elseExpr = &elseExprDecl
	}

	return newExpr(keyword, ast.IfExpression{
		Condition: condition,
		Then: thenExpr,
		Else: elseExpr,
	}), nil
}

func (p *Parser) parseFnExpr() (ast.Expression, diagnostic.Diagnostic) {
    keyword, _ := p.advance()

	params, returnType, body, diag := p.parseFunctionDefinition(); if diag != nil {
		return ast.Expression{}, diag
	}

	return newExpr(keyword, ast.FnExpression{
		Parameters: params,
		Return: returnType,
		Body: body,
	}), nil
}

// ---

func (p *Parser) parseAssignment(left ast.Expression, pos token.Position) (ast.Expression, diagnostic.Diagnostic) {
    operator, _ := p.advance()

	right, diag := p.parseExpression(); if diag != nil {
		return ast.Expression{}, diag
	}

	return p.makeAssignment(left, right, operator)
}

func (p *Parser) parseCall(left ast.Expression, pos token.Position) (ast.Expression, diagnostic.Diagnostic) {
	leftParen := p.current

	args, diag := p.parseArguments(); if diag != nil {
		return ast.Expression{}, diag
	}

	return newExpr(leftParen, ast.CallExpression{
		Callee: left,
		Arguments: args,
	}), nil
}

func (p *Parser) parseDot(left ast.Expression, pos token.Position) (ast.Expression, diagnostic.Diagnostic) {
	_, diag := p.expectToken(token.TokenDot); if diag != nil {
		return ast.Expression{}, diag
	}

	property, diag := p.expectToken(token.TokenIdentifier); if diag != nil {
		return ast.Expression{}, diag
	}

	return newExpr(property, ast.GetPropertyExpression{
		Left: left,
		Property: property,
	}), nil
}

func (p *Parser) parseRange(left ast.Expression, pos token.Position) (ast.Expression, diagnostic.Diagnostic) {
	operator, _ := p.advance()

	// TODO: check if this doesn't have a diagnostic
	inclusive := p.match(token.TokenEqual)

	right, diag := p.expression(PREC_RANGE); if diag != nil {
		return ast.Expression{}, diag
	}

	var step *ast.Expression = nil

	if p.match(token.TokenColon) {
		stepExpr, diag := p.expression(PREC_RANGE); if diag != nil {
			return ast.Expression{}, diag
		}

		step = &stepExpr
	}

	return newExpr(operator, ast.RangeExpression{
		Start: left,
		End: right,
		Step: step,
		Inclusive: inclusive,
	}), nil
}

func (p *Parser) parseAs(left ast.Expression, pos token.Position) (ast.Expression, diagnostic.Diagnostic) {
	operator, _ := p.advance()

	type_, diag := p.parseType(); if diag != nil {
		return ast.Expression{}, diag
	}

	return newExpr(operator, ast.AsExpression{
		Operand: left,
		ConvertType: type_,
	}), nil
}

// ---

func (p *Parser) parseUnary(kind token.TokenKind) (ast.Expression, diagnostic.Diagnostic) {
	pos := p.current.Pos

	operator, diag := p.expectToken(kind); if diag != nil {
		return ast.Expression{}, diag
	}

	operand, diag := p.expression(PREC_UNARY); if diag != nil {
		return ast.Expression{}, diag
	}
	
	return newExprToken(pos, operator, ast.UnaryExpression{
		Operand: operand,
		Operator: operator,
	}), nil
}

func (p *Parser) parseBinary(left ast.Expression, op token.TokenKind) (ast.Expression, diagnostic.Diagnostic) {
	prec := getPrecedenceFor(op)
	
	operator, diag := p.expectToken(op); if diag != nil {
		return ast.Expression{}, diag
	}

	right, diag := p.expression(prec); if diag != nil {
		return ast.Expression{}, diag
	}

	return newExpr(operator, ast.BinaryExpression{
		Left: left,
		Right: right,
		Operator: operator,
	}), nil
}

func (p *Parser) parseOperatorAssignment(left ast.Expression, finalOp token.TokenKind) (ast.Expression, diagnostic.Diagnostic) {
	operator, _ := p.advance()

	opAsFinal := operator
	opAsFinal.Kind = finalOp
	opAsFinal.Lexeme = string(finalOp)

	right, diag := p.parseExpression(); if diag != nil {
		return ast.Expression{}, diag
	}

	right.Data = ast.BinaryExpression{
		Left: left,
		Right: right,
		Operator: opAsFinal,
	}

	return p.makeAssignment(left, right, operator)
}

func (p *Parser) parseLogical(left ast.Expression, op token.TokenKind) (ast.Expression, diagnostic.Diagnostic) {
	prec := getPrecedenceFor(op)
	
	operator, diag := p.expectToken(op); if diag != nil {
		return ast.Expression{}, diag
	}

	shortCircuit := !p.match(token.TokenStar)

	right, diag := p.expression(prec); if diag != nil {
		return ast.Expression{}, diag
	}

	return newExpr(operator, ast.LogicalExpression{
		Left: left,
		Right: right,
		Operator: operator,
		ShortCircuit: shortCircuit,
	}), nil
}

// ---

func newExpr(tok token.Token, data ast.ExprData) ast.Expression {
	return ast.Expression{
		Base: ast.AstBase{
			Token: tok,
		},
		Data: data,
	}
}

func newExprToken(pos token.Position, tok token.Token, data ast.ExprData) ast.Expression {
	tok.Pos = pos

	return ast.Expression{
		Base: ast.AstBase{
			Token: tok,
		},
		Data: data,
	}
}
