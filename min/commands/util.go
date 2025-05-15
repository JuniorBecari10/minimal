package commands

import (
	"fmt"
	"time"
)

func measure(message string, operation func()) {
	fmt.Printf(" [..] %s... \n", message)
	
	start := time.Now()
	operation()
	
	timeAmount := time.Since(start).Milliseconds()
	fmt.Printf("      completed in %d ms.\n", timeAmount)
}

func measureNewline(message string, operation func()) {
	fmt.Printf(" [..] %s...\n", message)
	operation()
}

