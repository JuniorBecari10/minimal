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
			_, diag := p.advance()
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

func newType(tok token.Token, data types.TypeData) types.Type {
	return types.Type{
		Token: tok,
		Data: data,
	}
}
