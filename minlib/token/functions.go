package token

import (
	"fmt"
	"minlib/file"
)

func StartToken() Token {
	return Token{
		Kind: TokenStart,
	}
}

func EndToken(data *file.FileData) Token {
	return Token{
		Kind: TokenEnd,
		Pos: Position{
			Line: uint32(len(data.Lines) + 1),
			Col: uint32(len(data.Lines[len(data.Lines)-1]) + 2),
		},
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
