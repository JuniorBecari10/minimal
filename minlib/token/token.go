package token

type TokenKind string

const (
	TokenIntLiteral    = "int literal"
	TokenFloatLiteral  = "float literal"
	TokenCharLiteral   = "char literal"
	TokenStringLiteral = "string literal"
	TokenIdentifier    = "identifier"

	TokenPlus    = "'+'"
	TokenMinus   = "'-'"
	TokenStar    = "'*'"
	TokenSlash   = "'/'"
	TokenPercent = "'%'"

	TokenPlusEqual    = "'+='"
	TokenMinusEqual   = "'-='"
	TokenStarEqual    = "'*='"
	TokenSlashEqual   = "'/='"
	TokenPercentEqual = "'%='"

	TokenLeftParen  = "'('"
	TokenRightParen = "')'"

	TokenLeftBrace  = "'{'"
	TokenRightBrace = "'}'"

	TokenEqual       = "'='"
	TokenDoubleEqual = "'=='"
	TokenBangEqual   = "'!='" 

	TokenSemicolon = "';'"
	TokenComma     = "','"
	TokenColon     = "':'"
	TokenQuestion  = "'?'"

	TokenDot       = "'.'"
	TokenDoubleDot = "'..'"

	TokenArrow = "'->'"

	TokenGreater      = "'>'"
	TokenGreaterEqual = "'>='"

	TokenLess      = "'<'"
	TokenLessEqual = "'<='"

	TokenAsKw       = "'as' keyword"
	TokenIfKw       = "'if' keyword"
	TokenElseKw     = "'else' keyword"
	TokenWhileKw    = "'while' keyword"
	TokenForKw      = "'for' keyword"
	TokenLoopKw     = "'loop' keyword"
	TokenLetKw      = "'let' keyword"
	TokenVarKw      = "'var' keyword"
	TokenFnKw       = "'fn' keyword"
	TokenBreakKw    = "'break' keyword"
	TokenContinueKw = "'continue' keyword"
	TokenInKw       = "'in' keyword"
	TokenSelfKw     = "'self' keyword"
	TokenRecordKw   = "'record' keyword"
	TokenReturnKw   = "'return' keyword"
	TokenOutKw      = "'out' keyword"

	TokenAndKw = "'and' keyword"
	TokenOrKw  = "'or' keyword"
	TokenNotKw = "'not' keyword"

	TokenTrueKw  = "'true' keyword"
	TokenFalseKw = "'false' keyword"
	TokenNilKw   = "'nil' keyword"
	TokenVoidKw  = "'void' keyword"

	TokenStart = "start"
	TokenEnd = "end"
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
