package util

import (
	"hash/crc32"
	"encoding/binary"
	"fmt"
	"minlib/token"
	"path/filepath"
	"strconv"
	"strings"
	"os"
)

// It is little-endian
func IntToBytes(n int) []byte {
	bytes := make([]byte, 4)

	binary.LittleEndian.PutUint32(bytes, uint32(n))
	return bytes
}

func BytesToInt(b []byte) (int, error) {
	if len(b) != 4 {
		return 0, fmt.Errorf("input byte slice must be 4 bytes long")
	}

	return int(binary.LittleEndian.Uint32(b)), nil
}

// Removes the last element from the supplied list and returns it.
func PopList[T any](list *[]T) T {
	lastIndex := len(*list) - 1
	
	topElement := (*list)[lastIndex]
	*list = (*list)[:lastIndex]

	return topElement
}

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

// computeChecksum calculates the CRC32 checksum over the entire byte slice.
func computeChecksum(data []byte) uint32 {
    table := crc32.MakeTable(crc32.IEEE)
    return crc32.Checksum(data, table)
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
