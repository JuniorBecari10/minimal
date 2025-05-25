package parser

import (
	"minc/diagnostic"
	"minc/types"
	"minlib/token"
)

func (p *Parser) parseTypeAnnotation() (types.Type, diagnostic.Diagnostic) {
	_, diag := p.expectToken(token.TokenColon); if diag != nil {
		return nil, diag
	}

	return p.parseType()
}

func (p *Parser) parseType() (types.Type, diagnostic.Diagnostic) {
	typeToken, diag := p.advance(); if diag != nil {
		return nil, diag
	}

	switch typeToken.Lexeme {
		case "int": return types.TypeInt{}, nil
		case "float": return types.TypeFloat{}, nil
		case "str": return types.TypeStr{}, nil
		case "char": return types.TypeChar{}, nil
		case "bool": return types.TypeBool{}, nil

		// 'untyped nil' should not be written, only inferred.

		// for generics, like: ok<_, int>
		case "_": return types.TypeUnknown{}, nil
		case "void": return types.TypeVoid{}, nil

		case "fn": return p.parseFnType()
		case "range": return p.parseRangeType()

		default: {
			parsedType, diag := p.parseType(); if diag != nil {
				return nil, diag
			}

			if p.match(token.TokenQuestion) {
				return types.TypeOptional{
					Inside: parsedType,
				}, nil
			}

			// TODO: see what we can do with user types
		}
	}
}

func (p *Parser) parseFnType() (types.Type, diagnostic.Diagnostic) {
	// 'fn' keyword already advanced.

	params, diag := p.parseParameterTypes(); if diag != nil {
		return nil, diag
	}

	var returnType types.Type = types.TypeVoid{}
	if p.check(token.TokenColon) {
		returnType, diag = p.parseTypeAnnotation(); if diag != nil {
			return nil, diag
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
		return nil, diag
	}

	if len(args) != 1 {
		return nil, p.makeInvalidTypeArgumentsLengthDiagnostic(1, len(args))
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
		return []types.Type{}, diag
	}

	params := []types.Type{}

	for !p.match(right) {
		type_, diag := p.parseType(); if diag != nil {
			return []types.Type{}, diag
		}

		params = append(params, type_)
	}

	_, diag = p.expectToken(right); if diag != nil {
		return []types.Type{}, diag
	}

	return params, nil
}
