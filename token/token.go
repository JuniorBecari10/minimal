package token

type TokenKind string

const (
	TokenNumber     = "number"
	TokenString     = "string"
	TokenIdentifier = "identifier"

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
	TokenComma     = ","
	TokenPercent   = "%"

	TokenGreater      = ">"
	TokenGreaterEqual = ">="

	TokenLess      = "<"
	TokenLessEqual = "<="

	TokenIfKw     = "if keyword"
	TokenElseKw   = "else keyword"
	TokenWhileKw  = "while keyword"
	TokenForKw    = "for keyword"
	TokenVarKw    = "var keyword"
	TokenFnKw     = "fn keyword"
	TokenReturnKw = "return keyword"

	TokenAndKw = "and keyword"
	TokenOrKw  = "or keyword"
	TokenNotKw = "not keyword"

	TokenTrueKw  = "true keyword"
	TokenFalseKw = "false keyword"
	TokenNilKw   = "nil keyword"
	TokenVoidKw  = "void keyword"

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
