package run

import (
	"fmt"
	"os"
	"strings"

	"minc/compiler"
	"minc/lexer"
	"minc/parser"
	"minlib/util"
	"minlib/value"
)

func Compile(sourcePath, outputPath string) {
	sourceContent, err := util.ReadSourceFile(sourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading source file '%s'.\n", sourcePath)
		os.Exit(1)
	}

	source := string(sourceContent)

	fileData := util.FileData{
		Name:  sourcePath,
		Lines: strings.Split(source, "\n"),
	}

	chunk, hadError := compileSource(source, &fileData)
	if hadError {
		os.Exit(1)
	}

	if err := util.WriteBytecodeFile(outputPath, chunk.Serialize()); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to output file '%s': %v\n", outputPath, err)
		os.Exit(1)
	}
}

func compileSource(source string, fileData *util.FileData) (value.Chunk, bool) {
	lex := lexer.NewLexer(source, fileData)
	tokens, hadError := lex.Lex()
	
	if hadError {
		return value.Chunk{}, true
	}

	parse := parser.NewParser(tokens, fileData)
	ast, hadError := parse.Parse()
	
	if hadError {
		return value.Chunk{}, true
	}

	comp := compiler.NewCompiler(ast, fileData)
	chunk, hadError := comp.Compile()
	
	if hadError {
		return value.Chunk{}, true
	}

	return chunk, false
}

