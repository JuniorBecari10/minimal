package run

import (
	"fmt"
	"os"
	"strings"
	"vm-go/compiler"
	"vm-go/disassembler"
	"vm-go/lexer"
	"vm-go/parser"
	"vm-go/util"
	"vm-go/value"
)

type RunMode int

const (
	ModeCompile RunMode = iota
	ModeDisassemble
	ModeDeserialize
)

// 'output' is optional, but when mode is ModeCompile it must be set
func Run(source, file string, output *string, mode RunMode) {

	var chunk value.Chunk
	var hadError bool = false
	var fileData util.FileData

	fileName := util.GetFileName(file)

	if mode == ModeDeserialize {
		fileData = util.FileData{
			Name: fileName,
			Lines: []string{},
		}
		
		chunk = value.Deserialize([]byte(source))
	} else {
		fileData = util.FileData{
			Name: fileName,
			Lines: strings.Split(source, "\n"),
		}
		
		chunk, hadError = compile(source, &fileData)
	}

	if hadError {
		return
	}
	
	switch mode {
		case ModeCompile: {
			outputFile := *output
			
			file, err := os.Create(outputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot write to '%s'.\n", outputFile)
				os.Exit(1)
			}

			defer func() {
				if err := file.Close(); err != nil {
					panic(err)
				}
			}()
			
			if _, err := file.Write(chunk.Serialize()); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot write to '%s'.\n", outputFile)
				os.Exit(1)
			}
		}

		case ModeDisassemble, ModeDeserialize: {
			diss := disassembler.NewDisassembler(chunk)
			diss.Disassemble()
		}
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
