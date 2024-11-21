package run

import (
	"vm-go/compiler"
	"vm-go/disassembler"
	"vm-go/lexer"
	"vm-go/parser"
	"vm-go/value"
	"vm-go/vm"
)

type RunMode = int

const (
	ModeRun = iota
	ModeDisassemble
)

func Run(source string, mode RunMode) {
	bytecode, constants, hadError := compile(source)

	if hadError {
		return
	}
	
	switch mode {
		case ModeRun: {
			vm_ := vm.NewVM(bytecode, constants)
			vm_.Run()
		}

		case ModeDisassemble: {
			diss := disassembler.NewDisassembler(bytecode, constants)
			diss.Disassemble()
		}
	}
}

func compile(source string) ([]byte, []value.Value, bool) {
	lexer := lexer.NewLexer(source)
	tokens, hadError := lexer.Lex()

	if hadError {
		return nil, nil, true
	}

	parser := parser.NewParser(tokens)
	ast, hadError := parser.Parse()

	if hadError {
		return nil, nil, true
	}

	compiler := compiler.NewCompiler(ast)
	bytecode, constants, hadError := compiler.Compile()

	if hadError {
		return nil, nil, true
	}

	return bytecode, constants, false
}
