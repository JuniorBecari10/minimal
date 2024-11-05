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

	parser := parser.NewParser(tokens)
	ast, hadError := parser.Parse()

	fmt.Printf("%+v\n", ast)

	if hadError {
		return
	}

	compiler := compiler.NewCompiler(ast)
	instructions, constants := compiler.Compile()

	d := disassembler.NewDisassembler(instructions, constants)
	d.Disassemble()

	vm_ := vm.NewVM(instructions, constants)
	status := vm_.Run()

	if status != vm.STATUS_OK {
		fmt.Println(status)
		return
	}
}
