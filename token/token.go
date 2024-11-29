package token

type TokenKind string

const (
	TokenNumber     = "Number"
	TokenString     = "String"
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
	TokenPercent   = "%"

	TokenGreater      = ">"
	TokenGreaterEqual = ">="

	TokenLess      = "<"
	TokenLessEqual = "<="

	TokenIfKw    = "if"
	TokenElseKw  = "else"
	TokenWhileKw = "while"
	TokenForKw   = "for"
	TokenVarKw   = "var"
	TokenPrintKw = "print"

	TokenAndKw = "and"
	TokenOrKw  = "or"
	TokenXorKw = "xor"
	TokenNotKw = "not"

	TokenTrueKw  = "true"
	TokenFalseKw = "false"
	TokenNilKw   = "nil"

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
