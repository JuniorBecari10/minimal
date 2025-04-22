package main

import (
	"fmt"
	"minc/run"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: minc <source> <output>")
		os.Exit(1)
	}

	run.Compile(os.Args[1], os.Args[2])
}

