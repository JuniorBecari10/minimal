package main

import (
	"fmt"
	"os"
	"vm-go/run"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Usage: vm <source> [(-o <output>) | (-d | --dissassemble) | (-b | --bytecode)]")
		os.Exit(1)
	}

	c, err := os.ReadFile(os.Args[1])
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read file: '%s'\n", os.Args[1])
		os.Exit(1)
	}

	mode := run.ModeCompile
	var output *string = nil

	if len(os.Args) >= 3 && (os.Args[2] == "-d" || os.Args[2] == "--dissassemble") {
		mode = run.ModeDisassemble
	} else if len(os.Args) >= 3 && (os.Args[2] == "-b" || os.Args[2] == "--bytecode") {
		mode = run.ModeDeserialize
	}

	if mode == run.ModeCompile {
		if len(os.Args) >= 4 && os.Args[2] == "-o" {
			// must have specified the output file
			output = &os.Args[3]
		} else {
			fmt.Fprintln(os.Stderr, "Please specify the output file with '-o <output>'")
			os.Exit(1)
		}
	}

	run.Run(string(c), os.Args[1], output, mode)
}
