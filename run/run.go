package run

import (
	"fmt"
	"strings"
	"vm-go/chunk"
	"vm-go/compiler"
	"vm-go/disassembler"
	"vm-go/lexer"
	"vm-go/parser"
	"vm-go/util"
	"vm-go/vm"
)

type RunMode = int

const (
	ModeRun = iota
	ModeDisassemble
)

func Run(source, fileName string, mode RunMode) {
	fileData := util.FileData{
		Name: fileName,
		Lines: strings.Split(source, "\n"),
	}

	chunk, hadError := compile(source, &fileData)
	for _, p := range chunk.Positions {
		fmt.Println(p)
	}

	if hadError {
		return
	}
	
	switch mode {
		case ModeRun: {
			vm_ := vm.NewVM(chunk, &fileData)
			vm_.Run()
		}

		case ModeDisassemble: {
			diss := disassembler.NewDisassembler(chunk, &fileData)
			diss.Disassemble()
		}
	}
}

func compile(source string, fileData *util.FileData) (chunk.Chunk, bool) {
	lexer := lexer.NewLexer(source, fileData)
	tokens, hadError := lexer.Lex()

	if hadError {
		return chunk.Chunk{}, true
	}

	parser := parser.NewParser(tokens, fileData)
	ast, hadError := parser.Parse()

	if hadError {
		return chunk.Chunk{}, true
	}

	compiler := compiler.NewCompiler(ast, fileData)
	chunk_, hadError := compiler.Compile()

	if hadError {
		return chunk.Chunk{}, true
	}

	return chunk_, false
}
