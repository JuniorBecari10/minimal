package main

import (
	"fmt"
	"os"
	"vm-go/compiler"
	"vm-go/disassembler"
	"vm-go/lexer"
	"vm-go/parser"
	"vm-go/vm"
)

func main() {
	c, err := os.ReadFile("file.txt")
	
	if err != nil {
		fmt.Println("Cannot read file")
		os.Exit(1)
	}

	interpret(string(c))
}

func interpret(source string) {
	lexer := lexer.NewLexer(source)
	tokens, hadError := lexer.Lex()

	if hadError {
		return
	}

	fmt.Printf("%#v\n", tokens)

	parser := parser.NewParser(tokens)
	ast, hadError := parser.Parse()

	if hadError {
		return
	}
	
	compiler := compiler.NewCompiler(ast)
	instructions, constants, hadError := compiler.Compile()

	if hadError {
		return
	}
	
	d := disassembler.NewDisassembler(instructions, constants)
	d.Disassemble()
	fmt.Println()

	vm_ := vm.NewVM(instructions, constants)
	status := vm_.Run()

	if status != vm.STATUS_OK {
		return
	}
}
