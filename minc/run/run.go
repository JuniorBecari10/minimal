package run

import (
	"fmt"
	"minc/compiler"
	"minc/lexer"
	"minc/parser"
	"minlib/util"
	"minlib/value"
	"os"
	"strings"
)

// 'output' is optional, but when mode is ModeCompile it must be set
func Run(source, output string) {
	fileData := util.FileData{
		Name: source,
		Lines: strings.Split(source, "\n"),
	}
	
	chunk, hadError := compile(source, &fileData)

	if hadError {
		return
	}
			
	file, err := os.Create(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot write to '%s'.\n", output)
		os.Exit(1)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	
	if _, err := file.Write(chunk.Serialize()); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot write to '%s'.\n", output)
		os.Exit(1)
	}
}

func compile(source string, fileData *util.FileData) (value.Chunk, bool) {
	lexer := lexer.NewLexer(source, fileData)
	tokens, hadError := lexer.Lex()

	if hadError {
		return value.Chunk{}, true
	}

	parser := parser.NewParser(tokens, fileData)
	ast, hadError := parser.Parse()

	if hadError {
		return value.Chunk{}, true
	}

	compiler := compiler.NewCompiler(ast, fileData)
	chunk_, hadError := compiler.Compile()

	if hadError {
		return value.Chunk{}, true
	}

	return chunk_, false
}
