package commands

import (
	"fmt"
	"minlib/util"
	"minlib/value"
	"min/disassembler"
	"os"
)

func Build(source string) {

}

func Disasm(sourcePath string) {
	source, err := util.ReadSourceFile(sourcePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read source file '%s'.", sourcePath)
		os.Exit(1)
	}

	fmt.Println(source)
}

func Disasmb(bytecodePath string) {
	bytecode, err := util.ReadSourceFile(bytecodePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read bytecode file '%s'.", bytecode)
		os.Exit(1)
	}

	chunk := value.Deserialize([]byte(bytecode))
	disassembler.NewDisassembler(chunk).Disassemble()
}

func Execute(bytecode string) {

}

func Run(source string) {

}

