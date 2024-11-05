package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"vm-go/token"
	"vm-go/util"
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
		case '/': l.addToken(token.TokenSlash)

		case '(': l.addToken(token.TokenLeftParen)
		case ')': l.addToken(token.TokenRightParen)

		case '{': l.addToken(token.TokenLeftBrace)
		case '}': l.addToken(token.TokenRightBrace)

		case '=': l.addToken(token.TokenEqual)
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

// ---

func (l *Lexer) number() {
	for unicode.IsDigit(rune(l.peek())) {
		l.advance()
	}

	if l.match('.') && unicode.IsDigit(rune(l.peekNext())) {
		for unicode.IsDigit(rune(l.peek())) {
			l.advance()
		}
	}

	l.addToken(token.TokenNumber)
}

func (l *Lexer) identifier() {
	for unicode.IsLetter(rune(l.peek())) || l.peek() == '_' {
		l.advance()
	}

	l.addToken(l.checkKeyword())
}

func (l *Lexer) checkKeyword() token.TokenKind {
	switch l.source[l.start:l.current] {
		case "var": return token.TokenVarKw
		case "print": return token.TokenPrintKw

		default: return token.TokenIdentifier
	}
}

// ---

func (l *Lexer) match(c byte) bool {
	if l.isAtEnd() {
		return false
	}

	if l.source[l.current] == c {
		l.advance()
		return true
	}

	return false
}

func (l *Lexer) advance() byte {
	peek := l.peek()
	l.current += 1

	if peek == '\n' {
		l.increaseLine()
	} else {
		l.currentPos.Col += 1
	}

	return peek
}

func (l *Lexer) peek() byte {
	if (l.isAtEnd()) {
		return 0
	}

	return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
	if (l.isAtEnd()) {
		return 0
	}

	return l.source[l.current + 1]
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) increaseLine() {
	l.currentPos.Line += 1
	l.currentPos.Col = 0
}

func (l *Lexer) error(message string) {
	util.Error(l.startPos, message)
	l.hadError = true
}

// ---

func (l *Lexer) addToken(kind token.TokenKind) {
	l.tokens = append(l.tokens, token.Token{
		Kind: kind,
		Lexeme: l.source[l.start:l.current],
		Pos: l.startPos,
	})
}
