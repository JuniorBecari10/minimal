package main

import (
	"fmt"
	"minc/start"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: minc <source> <output>")
		os.Exit(1)
	}

	start.Compile(os.Args[1], os.Args[2])
}
