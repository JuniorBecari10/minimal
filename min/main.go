package main

import (
	"fmt"
	"min/commands"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	command := os.Args[1]
	argument := os.Args[2]

	switch command {
		case "build": commands.Build(argument)
		case "disasm": commands.Disasm(argument)
		case "disasmb": commands.Disasmb(argument)
		case "execute": commands.Execute(argument)
		case "run": commands.Run(argument)

		default: usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: min <command> [argument]")
	os.Exit(1)
}

