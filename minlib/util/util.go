package util

import (
	"hash/crc32"
	"encoding/binary"
	"fmt"
	"minlib/token"
	"path/filepath"
	"strconv"
	"strings"
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

func Error(pos token.Position, length int, message string, fileData *FileData) {
	fmt.Printf("[-] Error: %s\n", message)
	fmt.Printf(" | %s [-] %s (%d, %d)\n", strings.Repeat(" ", len(strconv.Itoa(pos.Line + 1))), fileData.Name, pos.Line + 1, pos.Col + 1)
	fmt.Printf(" |  %d | %s\n", pos.Line+1, fileData.Lines[pos.Line])
	fmt.Printf(" | %s  | %s%s\n", strings.Repeat(" ", len(strconv.Itoa(pos.Line+1))), strings.Repeat(" ", pos.Col), strings.Repeat("^", length))
	fmt.Printf(" | %s [-]\n", strings.Repeat(" ", len(strconv.Itoa(pos.Line + 1))))
	fmt.Printf("[-]\n\n")
}

