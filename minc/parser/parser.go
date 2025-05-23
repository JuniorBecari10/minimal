package parser

import (
	"minc/ast"
	"minc/lexer"
	"minlib/token"
	"minlib/file"
)

type Parser struct {
	lexer *lexer.Lexer

	current token.Token
	next token.Token

	prefixMap map[token.TokenKind]func() ast.Expression
	infixMap map[token.TokenKind]func(ast.Expression, token.Position) ast.Expression
	precedenceMap map[token.TokenKind]int

	panicMode bool
	fileData *file.FileData
}

func New(source string, fileData *file.FileData) *Parser {
	lexer := lexer.New(source, fileData)

	p := &Parser{
		lexer: lexer,

		// current and next set when starting to parse

		panicMode: false,
		fileData: fileData,
	}

	// TODO: set maps

	return p
}
