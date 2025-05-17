package commands

import (
	"fmt"
	"time"
)

func logMeasure(message string, operation func()) {
	fmt.Printf("[..] %s... \n", message)
	
	start := time.Now()
	operation()
	
	timeAmount := time.Since(start).Milliseconds()
	fmt.Printf("     completed in %d ms.\n\n", timeAmount)
}

func logNewline(message string, operation func()) {
	fmt.Printf("[..] %s...\n", message)
	operation()
}

func log(message string) {
	fmt.Printf("\n[..] %s\n", message)
}

