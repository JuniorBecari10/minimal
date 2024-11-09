package token

type TokenKind string

const (
	TokenNumber     = "Number"
	TokenIdentifier = "Identifier"

	TokenPlus  = "+"
	TokenMinus = "-"
	TokenStar  = "*"
	TokenSlash = "/"

	TokenLeftParen  = "("
	TokenRightParen = ")"

	TokenLeftBrace  = "{"
	TokenRightBrace = "}"

	TokenEqual     = "="
	TokenSemicolon = ";"

	TokenIfKw    = "if"
	TokenElseKw  = "else"
	TokenWhileKw = "while"
	TokenVarKw   = "var"
	TokenPrintKw = "print"

	TokenAbsent = "Absent"
)

type Position struct {
	Line int
	Col  int
}

type Token struct {
	Kind   TokenKind
	Lexeme string
	Pos    Position
}

func (t Token) IsAbsent() bool {
	return t.Kind == TokenAbsent
}

func AbsentToken() Token {
	return Token{
		Kind: TokenAbsent,
	}
}
