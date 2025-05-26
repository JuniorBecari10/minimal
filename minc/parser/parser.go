package parser

import (
	"minc/ast"
	"minc/lexer"
	"minlib/token"
	"minlib/file"
)

type ParserResult int

const (
	RES_OK ParserResult = iota
	RES_ERROR
)

type Parser struct {
	lexer *lexer.Lexer

	previous token.Token
	current token.Token
	next token.Token

	// Turned on if the lexer the parser owns had an error. The parser may not continue parsing, because
	// the subsequent tokens may be incomplete and thus not suitable for parsing.
	hadLexerError bool
	fileData *file.FileData
}

func New(source string, fileData *file.FileData) *Parser {
	lexer := lexer.New(source, fileData)

	p := &Parser{
		lexer: lexer,

		previous: token.StartToken(),
		// current and next set when starting to parse

		hadLexerError: false,
		fileData: fileData,
	}

	return p
}

func (p *Parser) setInitialTokens() ParserResult {
	current, diag := p.lexer.Lex()
	
	if diag != nil {
		diag.PrintDiagnostic()
		return RES_ERROR
	}

	p.current = current

	// ---

	next, diag := p.lexer.Lex()
	
	if diag != nil {
		diag.PrintDiagnostic()
		return RES_ERROR
	}

	p.next = next
	return RES_OK
}

func (p *Parser) Parse() ([]ast.Statement, ParserResult) {
	stmts := []ast.Statement{}
	res := p.setInitialTokens()

	if res == RES_ERROR {
		return stmts, res
	}

	for !p.current.IsEnd() {
		decl, diag := p.declaration(false, true)

		if diag != nil {
			diag.PrintDiagnostic()
			res = RES_ERROR

			p.synchronize()
		}

		stmts = append(stmts, decl)

		// The parser cannot recover from a lexer error.
		if p.hadLexerError {
			return stmts, RES_ERROR
		}
	}

	return stmts, res
}
