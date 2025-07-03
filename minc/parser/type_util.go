package parser

import (
	"minc/diagnostic"
	"minc/types"
	"minlib/token"
)

func (p *Parser) parseTypeAnnotation() (types.Type, diagnostic.Diagnostic) {
	_, diag := p.expectToken(token.TokenColon); if diag != nil {
		return types.Type{}, diag
	}

	return p.parseType()
}

func (p *Parser) parseFnType() (types.Type, diagnostic.Diagnostic) {
	// 'fn' keyword is already advanced.
	fnKw := p.previous

	params, diag := p.parseParameterTypes(); if diag != nil {
		return types.Type{}, diag
	}

	// dummy token, because this token doesn't exist and when omitted, the type is void, and I won't add another type for functions
	// with an optional return type.
	// This is a semantic problem, but for now let's keep it like this.
	var returnType types.Type = types.DummyType(types.TypeVoid{})

	if p.check(token.TokenColon) {
		returnType, diag = p.parseTypeAnnotation(); if diag != nil {
			return types.Type{}, diag
		}
	}

	return types.Type{
        Token: fnKw,

		Data: types.TypeFunction{
			Parameters: params,
			Return: returnType,
		},
	}, nil
}

func (p *Parser) parseRangeType() (types.Type, diagnostic.Diagnostic) {
	// 'range' already advanced.
	rangeTk := p.previous

	args, diag := p.parseTypeArguments(); if diag != nil {
		return types.Type{}, diag
	}

	if len(args) != 1 {
		return types.Type{}, p.makeInvalidTypeArgumentsLengthDiagnostic(1, len(args))
	}

	return types.Type{
		Token: rangeTk,

		Data: types.TypeRange{
			Inside: args[0],
		},
	}, nil
}

func (p *Parser) parseParameterTypes() ([]types.Type, diagnostic.Diagnostic) {
	return p.parseTypeList(token.TokenLeftParen, token.TokenRightParen)
}

func (p *Parser) parseTypeArguments() ([]types.Type, diagnostic.Diagnostic) {
	return p.parseTypeList(token.TokenLess, token.TokenGreater)
}

func (p *Parser) parseTypeList(left, right token.TokenKind) ([]types.Type, diagnostic.Diagnostic) {
	_, diag := p.expectToken(left); if diag != nil {
		return nil, diag
	}

	params := []types.Type{}

	// Handle empty list '<>'
	if p.check(right) {
		p.advance()
		return params, nil
	}

	for {
		type_, diag := p.parseType(); if diag != nil {
			return []types.Type{}, diag
		}

		params = append(params, type_)

		if p.match(token.TokenComma) {
			continue
		} else if p.match(right) {
			break
		} else {
			return nil, p.makeExpectedTokenDiagnostic(right) // or comma
		}
	}

	return params, nil
}
