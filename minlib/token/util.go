package token

import "fmt"

func StartToken() Token {
	return Token{
		Kind: TokenStart,
	}
}

func EndToken() Token {
	return Token{
		Kind: TokenEnd,
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
		case TokenInt, TokenFloat, TokenChar, TokenString, TokenIdentifier:
			return fmt.Sprintf("%s '%s'", t.Kind, t.Lexeme)
		
		default:
			return string(t.Kind)
	}
}
