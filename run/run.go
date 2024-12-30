package run

import (
	"fmt"
	"strings"
	"vm-go/compiler"
	"vm-go/disassembler"
	"vm-go/lexer"
	"vm-go/parser"
	"vm-go/util"
	"vm-go/value"
	"vm-go/vm"
)

type RunMode int

const (
	ModeRun RunMode = iota
	ModeDisassemble
)

func Run(source, fileName string, mode RunMode) {
	fileData := util.FileData{
		Name: fileName,
		Lines: strings.Split(source, "\n"),
	}

	chunk, hadError := compile(source, &fileData)

	if hadError {
		return
	}
	
	switch mode {
		case ModeRun: {
			vm_ := vm.NewVM(chunk, &fileData)
			result := vm_.Run()

			if result != vm.STATUS_OK {
				fmt.Printf("Exited with status %d.\n", result)
			}
		}

		case ModeDisassemble: {
			diss := disassembler.NewDisassembler(chunk, &fileData)
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
