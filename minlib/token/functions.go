package token

import (
	"fmt"
)

func StartToken() Token {
	return Token{
		Kind: TokenStart,
	}
}

func EndToken(lastToken Token) Token {
	// 1 more to keep a space between the tokens and guarantee they won't merge, if both are made of letters (keywords and identifiers).
	// let _
	//     ^ here
	lastToken.Pos.Col += uint32(len(lastToken.Lexeme) + 1)

	return Token{
		Kind: TokenEnd,
		Pos: lastToken.Pos,
		Lexeme: "",
	}
}

func (t Token) IsStart() bool {
	return t.Kind == TokenStart
}

func (t Token) IsEnd() bool {
	return t.Kind == TokenEnd
}

// ---

func (t Token) FormatError() string {
	switch t.Kind {
		case TokenIntLiteral, TokenFloatLiteral, TokenCharLiteral, TokenStringLiteral, TokenIdentifier:
			return fmt.Sprintf("%s '%s'", t.Kind, t.Lexeme)
		
		default:
			return string(t.Kind)
	}
}

func (t Token) Length() int {
	if t.Kind == TokenStringLiteral {
		return len(t.Lexeme) + 2
	} else {
		return len(t.Lexeme)
	}
}
