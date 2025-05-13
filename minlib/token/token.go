package token

type TokenKind string

const (
	TokenInt        = "int"
	TokenFloat      = "float"
	TokenChar       = "char"
	TokenString     = "string"
	TokenIdentifier = "identifier"

	TokenPlus    = "+"
	TokenMinus   = "-"
	TokenStar  	 = "*"
	TokenSlash   = "/"
	TokenPercent = "%"

	TokenPlusEqual    = "+="
	TokenMinusEqual   = "-="
	TokenStarEqual    = "*="
	TokenSlashEqual   = "/="
	TokenPercentEqual = "%="

	TokenLeftParen  = "("
	TokenRightParen = ")"

	TokenLeftBrace  = "{"
	TokenRightBrace = "}"

	TokenEqual       = "="
	TokenDoubleEqual = "=="
	TokenBangEqual   = "!="

	TokenSemicolon = ";"
	TokenComma     = ","
	TokenColon     = ":"
	TokenDot       = "."
	TokenDoubleDot = ".."

	TokenArrow = "->"

	TokenGreater      = ">"
	TokenGreaterEqual = ">="

	TokenLess      = "<"
	TokenLessEqual = "<="

	TokenIfKw       = "if keyword"
	TokenElseKw     = "else keyword"
	TokenWhileKw    = "while keyword"
	TokenForKw      = "for keyword"
	TokenLoopKw     = "loop keyword"
	TokenVarKw      = "var keyword"
	TokenFnKw       = "fn keyword"
	TokenBreakKw    = "break keyword"
	TokenContinueKw = "continue keyword"
	TokenInKw       = "in keyword"
	TokenSelfKw     = "self keyword"
	TokenRecordKw   = "record keyword"
	TokenReturnKw   = "return keyword"

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
	Line uint32
	Col  uint32
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
