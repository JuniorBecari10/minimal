package commands

import (
	"bytes"
	"fmt"
	"min/disassembler"
	"minlib/util"
	"minlib/value"
	"os"
	"os/exec"
	"time"
)

func Build(source, output string) {
	// just run the compiler
	minc := getMinc()
	
	cmd := exec.Command(minc, source, output)
	
	cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

	fmt.Print(" [..] Compiling... ")
	start := time.Now()

    err := cmd.Run()

	fmt.Printf("completed in %d ms.\n", time.Since(start).Milliseconds())
    
	if err != nil {
        os.Exit(1)
    }
}

func Disasm(source string) {
	// run compiler and disassemble
	minc := getMinc()
	
	cmd := exec.Command(minc, source, "*stdout")
    stdout, err := cmd.Output()
    
	if err != nil {
        os.Exit(1)
    }

	data, err := util.ReadBytecode(*bytes.NewBuffer(stdout))

	if err != nil {
		os.Exit(1)
	}

	chunk := value.Deserialize(data)
	disassembler.NewDisassembler(chunk).Disassemble()
}

func Disasmb(bytecodePath string) {
	// just disassemble
	bytecode, err := util.ReadBytecodeFile(bytecodePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read bytecode file '%s': %s\n", bytecodePath, err.Error())
		os.Exit(1)
	}

	chunk := value.Deserialize([]byte(bytecode))
	disassembler.NewDisassembler(chunk).Disassemble()
}

func Execute(bytecode string) {
	// just run the bytecode
	minvm := getMinvm()
	
	cmd := exec.Command(minvm, bytecode)
	
	cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err := cmd.Run()
    
	if err != nil {
        os.Exit(1)
    }
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

