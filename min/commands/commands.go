package commands

import (
	"fmt"
	"min/disassembler"
	"minlib/util"
	"minlib/value"
	"os"
)

func Build(source string) {

}

func Disasm(sourcePath string) {
	source, err := util.ReadSourceFile(sourcePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read source file '%s'.", sourcePath)
	}

	
}

func Disasmb(bytecodePath string) {
	bytecode, err := util.ReadSourceFile(bytecodePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read bytecode file '%s'.", bytecode)
	}

	chunk := value.Deserialize([]byte(bytecode))
	disassembler.NewDisassembler(chunk).Disassemble()
}

func Execute(bytecode string) {

}

func Run(source string) {

}

