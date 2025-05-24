package lexer

import (
	"minc/diagnostic"
	"minlib/file"
	"minlib/token"
	"strings"
	"unicode"
)

type Lexer struct {
	source string

	start int
	current int

	startPos token.Position
	currentPos token.Position

	fileData *file.FileData
}

func New(source string, fileData *file.FileData) *Lexer {
	return &Lexer{
		source: source,

		start: 0,
		current: 0,

		startPos: token.Position{},
		currentPos: token.Position{},

		fileData: fileData,
	}
}

func (l *Lexer) Lex() (token.Token, diagnostic.Diagnostic) {
	for strings.IndexByte(" \r\t\n", l.peek(0)) != -1 {
		l.advance()
	}

	l.start = l.current
	l.startPos = l.currentPos

	c := l.advance()
	
	if c == 0 {
		return token.EndToken(), nil
	}

	switch c {
		case '+': {
			if l.match('=') {
				return l.makeToken(token.TokenPlusEqual), nil
			} else {
				return l.makeToken(token.TokenPlus), nil
			}
		}

		case '-': {
			if l.match('>') {
				return l.makeToken(token.TokenArrow), nil
			} else if l.match('=') {
				return l.makeToken(token.TokenMinusEqual), nil
			} else {
				return l.makeToken(token.TokenMinus), nil
			}
		}

		case '*': {
			if l.match('=') {
				return l.makeToken(token.TokenStarEqual), nil
			} else {
				return l.makeToken(token.TokenStar), nil
			}
		}

		case '/': {
			if l.match('/') {
				// Comment. Skip to end of line
				for l.peek(0) != '\n' && !l.isAtEnd(0) {
					l.advance()
				}

				// Try to lex again another token
				return l.Lex()
			} else if l.match('=') {
				return l.makeToken(token.TokenSlashEqual), nil
			} else {
				return l.makeToken(token.TokenSlash), nil
			}
		}

		case '%': {
			if l.match('=') {
				return l.makeToken(token.TokenPercentEqual), nil
			} else {
				return l.makeToken(token.TokenPercent), nil
			}
		}

		case '(': return l.makeToken(token.TokenLeftParen), nil
		case ')': return l.makeToken(token.TokenRightParen), nil

		case '{': return l.makeToken(token.TokenLeftBrace), nil
		case '}': return l.makeToken(token.TokenRightBrace), nil

		case '=': {
			if l.match('=') {
				return l.makeToken(token.TokenDoubleEqual), nil
			} else {
				return l.makeToken(token.TokenEqual), nil
			}
		}

		case '!': {
			if l.match('=') {
				return l.makeToken(token.TokenBangEqual), nil
			} else {
				return token.Token{}, l.makeUnknownTokenDiagnostic(c)
			}
		}

		case '>': {
			if l.match('=') {
				return l.makeToken(token.TokenGreaterEqual), nil
			} else {
				return l.makeToken(token.TokenGreater), nil
			}
		}

		case '<': {
			if l.match('=') {
				return l.makeToken(token.TokenLessEqual), nil
			} else {
				return l.makeToken(token.TokenLess), nil
			}
		}

		case ';': return l.makeToken(token.TokenSemicolon), nil

		case '"': return l.string()
		case '\'': return l.char()

		case ',': return l.makeToken(token.TokenComma), nil
		case ':': return l.makeToken(token.TokenColon), nil

		case '.': {
			if l.match('.') {
				return l.makeToken(token.TokenDoubleDot), nil
			} else {
				return l.makeToken(token.TokenDot), nil
			}
		}

		default: {
			if unicode.IsDigit(rune(c)) {
				return l.number(), nil
			} else if unicode.IsLetter(rune(c)) || c == '_' {
				return l.identifier(), nil
			} else {
				return token.Token{}, l.makeUnknownTokenDiagnostic(c)
			}
		}
	}
}
