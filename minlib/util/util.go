package util

import (
	"fmt"
	"minlib/token"
	"path/filepath"
	"strconv"
	"strings"
	"os"
)

func GetFileName(path string) string {
	// Use filepath to normalize and extract the file name
	normalizedPath := filepath.FromSlash(path) // Converts / to \ or vice versa based on OS
	return filepath.Base(normalizedPath)
}

func PadLeft(s string, length int, padChar string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(padChar, length-len(s))
	return padding + s
}

func PadRight(s string, length int, padChar string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(padChar, length-len(s))
	return s + padding
}

func Center(str string, width int, padChar string) string {
	// If the string is already wider than the desired width, return it as-is
	if len(str) >= width {
		return str
	}

	// Calculate the total padding required
	padding := width - len(str)

	// Split the padding into left and right (prefer left padding if it's odd)
	leftPadding := padding / 2
	rightPadding := padding - leftPadding

	// Pad the string on both sides
	return strings.Repeat(padChar, leftPadding) + str + strings.Repeat(" ", rightPadding)
}

func PrintError(pos token.Position, length int, message string, help *string, fileData *FileData) {
	lineNum := int(pos.Line + 1)
	lineStr := strconv.Itoa(lineNum)
	pad := strings.Repeat(" ", len(lineStr))

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, " [-] Error: %s\n", message)
	fmt.Fprintf(os.Stderr, "  | %s [-] %s (line %d, col %d)\n", pad, fileData.Name, lineNum, pos.Col+1)
	fmt.Fprintf(os.Stderr, "  |  %d | %s\n", lineNum, fileData.Lines[pos.Line])
	fmt.Fprintf(os.Stderr, "  | %s  | %s%s\n", pad,
		strings.Repeat(" ", int(pos.Col)), strings.Repeat("^", length))
	fmt.Fprintf(os.Stderr, "  | %s [-]\n", pad)

	if help != nil {
		fmt.Fprintln(os.Stderr, "  |")
		fmt.Fprintf(os.Stderr, "  | %s [-] Help: %s\n", pad, *help)
	}

	fmt.Fprintln(os.Stderr, " [-]")
}
