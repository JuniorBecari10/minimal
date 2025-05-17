package commands

import (
	"bytes"
	"fmt"
	"min/disassembler"
	"minlib/util"
	"minlib/value"
	"os"
	"os/exec"
	"path/filepath"
)

// TODO: add flag silent to remove these 'min' log messages.
func Build(source, output string) {
	// just run the compiler
	logMeasure("Compiling", func() {
		minc := getMinc()
		
		cmd := exec.Command(minc, source, output)
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		err := cmd.Run()

		if err != nil {
			log("Compiling phase failed.")
			os.Exit(1)
		}
	})
}

func Disasm(source string) {
	// run compiler and disassemble
	var data []byte

	logMeasure("Compiling", func() {
		minc := getMinc()
		
		cmd := exec.Command(minc, source, "*stdout")
		cmd.Stderr = os.Stderr

		stdout, err := cmd.Output()

		if err != nil {
			log("Compiling phase failed.")
			os.Exit(1)
		}

		data, err = util.ReadBytecode(*bytes.NewBuffer(stdout))

		if err != nil {
			log("Compiling phase failed.")
			os.Exit(1)
		}
	})

	logNewline("Disassembling", func() {
		chunk := value.Deserialize(data)
		disassembler.NewDisassembler(chunk).Disassemble()
	})
}

func Disasmb(bytecodePath string) {
	// just disassemble
	logNewline("Disassembling", func() {
		bytecode, err := util.ReadBytecodeFile(bytecodePath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read bytecode file '%s'.\n", bytecodePath)
			log("Disassembling phase failed.")
			
			os.Exit(1)
		}

		chunk := value.Deserialize([]byte(bytecode))
		disassembler.NewDisassembler(chunk).Disassemble()
	})
}

func Execute(bytecode string) {
	// just run the bytecode
	logNewline("Running", func() {
		minvm := getMinvm()
		
		cmd := exec.Command(minvm, bytecode)
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		
		if err != nil {
			log("Running phase failed.")
			os.Exit(1)
		}
	})
}

func Run(source string) {
	tempFile, err := os.CreateTemp("", "temp_*.mnb")
	
	if err != nil {
		log("Error generating temporary file for compilation.")
		os.Exit(1)
	}

	defer tempFile.Close()
	
	tempFileName, err := filepath.Abs(tempFile.Name())
	
	if err != nil {
		log("Error getting temporary file name for compilation.")
		os.Exit(1)
	}

	// compile and run
	logMeasure("Compiling", func() {
		minc := getMinc()
		
		cmd := exec.Command(minc, source, tempFileName)
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		err := cmd.Run()

		if err != nil {
			log("Compiling phase failed.")
			os.Exit(1)
		}
	})

	logNewline("Running", func() {
		minvm := getMinvm()
		cmd := exec.Command(minvm, tempFileName)
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		
		if err != nil {
			log("Running phase failed.")
			os.Exit(1)
		}
	})
}

// ---

// These functions end the program if the executable is not found.
// They also print a message.

func getCommand(command string) string {
	path, err := exec.LookPath(command)
	if err != nil {
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

