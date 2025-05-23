package lexer

import (
	"minlib/token"
	"minlib/file"
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
