package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"vm-go/token"
)

type Lexer struct {
	source string

	start   int
	current int
	
	startPos   token.Position
	currentPos token.Position

	hadError bool
	tokens []token.Token
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source:  source,

		start:   0,
		current: 0,
		
		startPos:   token.Position{},
		currentPos: token.Position{},

		hadError: false,
		tokens:   []token.Token{},
	}
}

func (l *Lexer) Lex() ([]token.Token, bool) {
	for !l.isAtEnd() {
		l.scanToken()
	}

	return l.tokens, l.hadError
}

// ---

func (l *Lexer) scanToken() {
	for strings.IndexByte(" \r\t\n", l.peek()) != -1 {
		if (l.peek() == '\n') {
			l.current += 1
			l.increaseLine()

			continue
		}

		l.advance()
	}

	l.start = l.current
	l.startPos = l.currentPos

	c := l.advance()

	if c == 0 {
		return
	}

	switch c {
		case '+': l.addToken(token.TokenPlus)
		case '-': l.addToken(token.TokenMinus)
		case '*': l.addToken(token.TokenStar)
		case '/': {
			if l.match('/') {
				// A comment goes until the end of the line.
				for l.peek() != '\n' && !l.isAtEnd() {
					l.advance()
				}
			} else {
				l.addToken(token.TokenSlash)
			}
		}

		case '(': l.addToken(token.TokenLeftParen)
		case ')': l.addToken(token.TokenRightParen)

		case '{': l.addToken(token.TokenLeftBrace)
		case '}': l.addToken(token.TokenRightBrace)

		case '=': {
			if l.match('=') {
				l.addToken(token.TokenDoubleEqual)
			} else {
				l.addToken(token.TokenEqual)
			}
		}

		case '!': {
			if l.match('=') {
				l.addToken(token.TokenBangEqual)
			} else {
				l.error(fmt.Sprintf("Unknown character: '%c' (%d)", c, int(c)))
			}
		}

		case '>': {
			if l.match('=') {
				l.addToken(token.TokenGreaterEqual)
			} else {
				l.addToken(token.TokenGreater)
			}
		}

		case '<': {
			if l.match('=') {
				l.addToken(token.TokenLessEqual)
			} else {
				l.addToken(token.TokenLess)
			}
		}

		case ';': l.addToken(token.TokenSemicolon)

		default: {
			if unicode.IsDigit(rune(c)) {
				l.number()
			} else if unicode.IsLetter(rune(c)) || c == '_' {
				l.identifier()
			} else {
				l.error(fmt.Sprintf("Unknown character: '%c' (%d)", c, int(c)))
			}
		}
	}
}
