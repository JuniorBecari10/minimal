package parser

import (
	"minc/ast"
	"minc/lexer"
	"minlib/file"
	"minlib/token"
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

	hadError bool

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

		hadError: false,
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
		stmt := p.parseTopLevelDeclaration()

		if p.hadError {
			continue
		}

		// The parser cannot recover from a lexer error.
		if p.hadLexerError {
			return stmts, RES_ERROR
		}
		
		stmts = append(stmts, stmt)
	}

	return stmts, res
}
