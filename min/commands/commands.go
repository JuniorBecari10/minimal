package commands

import (
	"bytes"
	"fmt"
	"io"
	"min/disassembler"
	"minlib/util"
	"minlib/value"
	"os"
	"os/exec"
)

func Build(source, output string) {
	// just run the compiler
	measure("Compiling", func() {
		minc := getMinc()
		
		cmd := exec.Command(minc, source, output)
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		err := cmd.Run()

		if err != nil {
			os.Exit(1)
		}
	})
}

func Disasm(source string) {
	// run compiler and disassemble
	var data []byte

	measure("Compiling", func() {
		minc := getMinc()
		
		cmd := exec.Command(minc, source, "*stdout")

		stdoutPipe, _ := cmd.StdoutPipe()
		cmd.Stderr = os.Stderr

		var stdoutBuf bytes.Buffer

		stdoutReader := io.TeeReader(stdoutPipe, os.Stdout)

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error while running compiler: %s\n", err.Error())
			os.Exit(1)
		}

		io.Copy(&stdoutBuf, stdoutReader)

		if err := cmd.Wait(); err != nil {
			os.Exit(1)
		}

		stdout := stdoutBuf.Bytes()
		
		var err error
		data, err = util.ReadBytecode(*bytes.NewBuffer(stdout))

		if err != nil {
			os.Exit(1)
		}
	})

	measureNewline("Disassembling", func() {
		chunk := value.Deserialize(data)
		disassembler.NewDisassembler(chunk).Disassemble()
	})
}

func Disasmb(bytecodePath string) {
	// just disassemble
	measureNewline("Disassembling", func() {
		bytecode, err := util.ReadBytecodeFile(bytecodePath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read bytecode file '%s': %s\n", bytecodePath, err.Error())
			os.Exit(1)
		}

		chunk := value.Deserialize([]byte(bytecode))
		disassembler.NewDisassembler(chunk).Disassemble()
	})
}

func Execute(bytecode string) {
	// just run the bytecode
	measureNewline("Running", func() {
		minvm := getMinvm()
		
		cmd := exec.Command(minvm, bytecode)
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		
		if err != nil {
			os.Exit(1)
		}
	})
}

func Run(source string) {
	// compile and run
	minc := getMinc()
	minvm := getMinvm()
	
	cmd := exec.Command(minc, source, "*stdout", "|", minvm, "*stdin")
	
	cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err := cmd.Run()
    
	if err != nil {
        os.Exit(1)
    }
}

// ---

// These functions end the program if the executable is not found.
// They also print a message.

func getCommand(command string) string {
	path, err := exec.LookPath(command)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprintf(os.Stderr, "The desired command requires '%s' to be in path.\n", command)
		fmt.Fprintln(os.Stderr, "Please place it in path before running this command again.")
		os.Exit(1)
	}

	return path
}

func getMinc() string {
	return getCommand("minc")
}

func getMinvm() string {
	return getCommand("minvm")
}

