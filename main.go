package main

import (
	"fmt"
	"os"
	"vm-go/run"
)

func main() {
	if len(os.Args) == 1 || len(os.Args) > 3 {
		fmt.Println("Usage: vm <source> [-d | --dissassemble]")
		return
	}

	c, err := os.ReadFile(os.Args[1])
	
	if err != nil {
		fmt.Printf("Cannot read file: '%s'\n", os.Args[1])
		os.Exit(1)
	}

	mode := run.ModeRun
	if len(os.Args) == 3 && (os.Args[2] == "-d" || os.Args[2] == "--dissassemble") {
		mode = run.ModeDisassemble
	}

	run.Run(string(c), mode)
}
