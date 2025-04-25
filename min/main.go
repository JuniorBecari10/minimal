package main

import (
	"fmt"
	"min/commands"
	"os"
)

func main() {
	if len(os.Args) < 3 || len(os.Args) > 4 {
		usage()
	}

	command := os.Args[1]
	arg1 := os.Args[2]

	switch command {
		case "build": {
			if len(os.Args) < 4 {
				fmt.Fprintln(os.Stderr, "Usage: min build <source> <output>")
				os.Exit(1)
			}

			arg2 := os.Args[3]
			commands.Build(arg1, arg2)
		}

		case "disasm": commands.Disasm(arg1)
		case "disasmb": commands.Disasmb(arg1)
		case "execute": commands.Execute(arg1)
		case "run": commands.Run(arg1)

		default: {
			fmt.Fprintf(os.Stderr, "Invalid command: '%s'.\n", command)
			usage()
		}
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: min <command> <arguments...>")
	os.Exit(1)
}

