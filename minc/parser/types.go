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


func (p *Parser) parseType() (types.Type, diagnostic.Diagnostic) {
	// Grouping with parentheses
	if p.match(token.TokenLeftParen) {
		innerType, diag := p.parseType(); if diag != nil {
			return types.Type{}, diag
		}

		_, diag = p.expectToken(token.TokenRightParen); if diag != nil {
			return types.Type{}, diag
		}

		// Allow postfix optional for grouped type
		if p.match(token.TokenQuestion) {
			return types.Type{
				Token: innerType.Token,
				Data: types.TypeOptional{
					Inside: innerType,
				},
			}, nil
		}

		return innerType, nil
	}

	// Normal type parsing

	endDiag := p.makeExpectedTypeDiagnostic()

	typeToken, diag := p.advance(); if diag != nil {
		return types.Type{}, diag
	}

	var baseType types.Type

	newType := func(data types.TypeData) types.Type {
		return types.Type{
			Token: typeToken,
			Data: data,
		}
	}

	switch typeToken.Lexeme {
		case "int": baseType = newType(types.TypeInt{})
		case "float": baseType = newType(types.TypeFloat{})
		case "str": baseType = newType(types.TypeStr{})
		case "char": baseType = newType(types.TypeChar{})
		case "bool": baseType = newType(types.TypeBool{})
		
		// 'untyped nil' should not be written, only inferred.

		// for generics, like: ok<_, int>
		case "_": baseType = newType(types.TypeUnknown{})
		case "void": baseType = newType(types.TypeVoid{})
		case "fn": return p.parseFnType()
		case "range": return p.parseRangeType()

		default: {
			if typeToken.Kind == token.TokenIdentifier {
				baseType = types.TypeUserDefined{
					Name: typeToken.Lexeme,
				}
			} else {
				// diagnostic is already prepared before the token advances
				return types.Type{}, endDiag
			}
		}
	}

	// Apply optional modifier
	if p.match(token.TokenQuestion) {
		baseType = types.TypeOptional{
			Inside: baseType,
		}
	}

	return baseType, nil
}

func (p *Parser) parseFnType() (types.Type, diagnostic.Diagnostic) {
	// 'fn' keyword is already advanced.

	params, diag := p.parseParameterTypes(); if diag != nil {
		return types.Type{}, diag
	}

	var returnType types.Type = types.TypeVoid{}
	if p.check(token.TokenColon) {
		returnType, diag = p.parseTypeAnnotation(); if diag != nil {
			return types.Type{}, diag
		}
	}

	return types.TypeFunction{
		Parameters: params,
		Return: returnType,
	}, nil
}

func (p *Parser) parseRangeType() (types.Type, diagnostic.Diagnostic) {
	// 'range' already advanced.

	args, diag := p.parseTypeArguments(); if diag != nil {
		return types.Type{}, diag
	}

	if len(args) != 1 {
		return types.Type{}, p.makeInvalidTypeArgumentsLengthDiagnostic(1, len(args))
	}

	return types.TypeRange{
		Inside: args[0],
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
		return types.Type{}, diag
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
			return types.Type{}, p.makeExpectedTokenDiagnostic(right) // or comma
		}
	}

	return params, nil
}
