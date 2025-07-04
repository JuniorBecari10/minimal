package parser

import (
	"minc/diagnostic"
	"minc/types"
	"minlib/token"
)

type TypePrefixFn = func() (types.Type, diagnostic.Diagnostic)
type TypeInfixFn = func(types.Type) (types.Type, diagnostic.Diagnostic)

type TypePrecedence int
const (
	TYPE_PREC_LOWEST TypePrecedence = iota
	TYPE_PREC_POSTFIX    // ?
	// TYPE_PREC_PRIMARY // int, float, str (unused)
)

func (p *Parser) parseType() (types.Type, diagnostic.Diagnostic) {
	return p._type(TYPE_PREC_LOWEST)
}

func (p *Parser) _type(prec TypePrecedence) (types.Type, diagnostic.Diagnostic) {
	prefixFn := p.getTypePrefixFn(p.current)

	if prefixFn == nil {
		return types.Type{}, p.makeExpectedTypeDiagnostic()
	}

	left, diag := prefixFn(); if diag != nil {
		return types.Type{}, diag
	}

	for getTypePrecedenceFor(p.current.Kind) > prec {
		infixFn := p.getTypeInfixFn(p.current)

		if infixFn == nil {
			break
		}

		left, diag = infixFn(left); if diag != nil {
			return types.Type{}, diag
		}
	}

	return left, nil
}

// needs entire token because the lexeme will be analyzed as well.
func (p *Parser) getTypePrefixFn(tok token.Token) TypePrefixFn {
	newFn := func(t types.TypeData) TypePrefixFn {
        return func() (types.Type, diagnostic.Diagnostic) {
			_, diag := p.advance() // this advances the first token, making each handler not need to do it.
			return newType(tok, t), diag
		}
	}

	switch tok.Kind {
		case token.TokenIdentifier: {
			switch tok.Lexeme {
				case "int": return newFn(types.TypeInt{})
				case "float": return newFn(types.TypeFloat{})
				case "str": return newFn(types.TypeStr{})
				case "char": return newFn(types.TypeChar{})
				case "bool": return newFn(types.TypeBool{})
				case "_": return newFn(types.TypeUnknown{})
				case "void": return newFn(types.TypeVoid{})
				case "range": return p.parseRangeType

				case "(": return p.parseGroupType
			}
		}

		case token.TokenFnKw: return p.parseFnType
	}

	return func() (types.Type, diagnostic.Diagnostic) {
		return types.Type{}, p.makeExpectedTypeDiagnostic()
	}
}

func (p *Parser) getTypeInfixFn(tok token.Token) TypeInfixFn {
	switch tok.Kind {
		case token.TokenQuestion: {
			return func(t types.Type) (types.Type, diagnostic.Diagnostic) {
				_, diag := p.advance()

				return newType(tok, types.TypeOptional{
					Inside: t,
				}), diag
			}
		}
	}

	return func(types.Type) (types.Type, diagnostic.Diagnostic) {
		return types.Type{}, p.makeExpectedTokenDiagnostic(token.TokenQuestion)
	}
}

func getTypePrecedenceFor(kind token.TokenKind) TypePrecedence {
	switch kind {
		case token.TokenQuestion:
			return TYPE_PREC_POSTFIX
		
		default:
			return TYPE_PREC_LOWEST
	}
}

// ---

func (p *Parser) parseFnType() (types.Type, diagnostic.Diagnostic) {
	// 'fn' keyword is already advanced.
	fnKw := p.previous

	params, diag := p.parseParameterTypes(); if diag != nil {
		return types.Type{}, diag
	}

	// dummy token, because this token doesn't exist and when omitted, the type is void,
	// and I won't add another type for functions with an optional return type.
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

func (p *Parser) parseGroupType() (types.Type, diagnostic.Diagnostic) {
	// '(' already advanced.

	t, diag := p.parseType(); if diag != nil {
		return types.Type{}, diag
	}

	_, diag = p.expectToken(token.TokenRightParen)
	return t, diag
}

// ---

func newType(tok token.Token, data types.TypeData) types.Type {
	return types.Type{
		Token: tok,
		Data: data,
	}
}
