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

// Run compiles the source file and writes the compiled output to the specified file.
func Run(sourcePath, outputPath string) {
	sourceContent, err := readFile(sourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading source file '%s': %v\n", sourcePath, err)
		os.Exit(1)
	}

	fileData := util.FileData{
		Name:  sourcePath,
		Lines: strings.Split(sourceContent, "\n"),
	}

	chunk, hadError := compileSource(sourceContent, &fileData)
	if hadError {
		os.Exit(1)
	}

	if err := writeOutputFile(outputPath, chunk.Serialize()); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to output file '%s': %v\n", outputPath, err)
		os.Exit(1)
	}
}

// readFile reads the source file content as a string.
func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// writeOutputFile writes serialized bytecode to a file.
func writeOutputFile(path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

// compileSource runs the lexer, parser, and compiler phases.
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

