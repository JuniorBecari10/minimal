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
