package run

import (
	"fmt"
	"io"
	"os"
	"strings"

	"minc/compiler"
	"minc/lexer"
	"minc/parser"
	"minlib/util"
	"minlib/value"
)

const (
	STDIN = "*stdin"
	STDOUT = "*stdout"
	STDERR = "*stderr"
)

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

func readFile(path string) (string, error) {
	// check for special '*stdin'
	if path == STDIN {
		input, err := io.ReadAll(os.Stdin)
		
		if err != nil {
			return "", err
		}

		return string(input), nil
	}

	data, err := os.ReadFile(path)
	
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func writeOutputFile(path string, data []byte) error {
	var out io.Writer

	// check for special '*stdout' and '*stderr'
	switch strings.ToLower(path) {
		case "*stdout":
			out = os.Stdout
		
		case "*stderr":
			out = os.Stderr
		
		default: {
			file, err := os.Create(path)
			
			if err != nil {
				return err
			}
			
			defer file.Close()
			out = file
		}
	}

	_, err := out.Write(data)
	return err
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

