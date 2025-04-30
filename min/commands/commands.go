package commands

import (
	"fmt"
	"minlib/util"
	"minlib/value"
	"min/disassembler"
	"os"
)

func Build(source, output string) {

}

func Disasm(sourcePath string) {
	source, err := util.ReadSourceFile(sourcePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read source file '%s': %s\n", sourcePath, err.Error())
		os.Exit(1)
	}

	fmt.Println(source)
}

func Disasmb(bytecodePath string) {
	bytecode, err := util.ReadBytecodeFile(bytecodePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read bytecode file '%s': %s\n", bytecodePath, err.Error())
		os.Exit(1)
	}

	chunk := value.Deserialize([]byte(bytecode))
	disassembler.NewDisassembler(chunk).Disassemble()
}

func Execute(bytecode string) {

}

func Run(source string) {

}

