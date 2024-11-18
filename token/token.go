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

	TokenEqual       = "="
	TokenDoubleEqual = "=="

	TokenBangEqual = "!="
	TokenSemicolon = ";"

	TokenGreater      = ">"
	TokenGreaterEqual = ">="

	TokenLess      = "<"
	TokenLessEqual = "<="

	TokenIfKw    = "if"
	TokenElseKw  = "else"
	TokenWhileKw = "while"
	TokenVarKw   = "var"
	TokenPrintKw = "print"

	TokenAndKw = "and"
	TokenOrKw  = "or"
	TokenXorKw = "xor"
	TokenNotKw = "not"

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
