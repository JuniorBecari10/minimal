package util

import (
	"encoding/binary"
	"fmt"
	"strings"
	"vm-go/token"
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

func Error(pos token.Position, length int, message string, fileData *FileData) {
	fmt.Printf("[-] Error at %s (%d, %d): %s\n", fileData.Name, pos.Line + 1, pos.Col + 1, message)
	fmt.Printf(" | %s\n", fileData.Lines[pos.Line])
	fmt.Printf(" | %s%s\n", strings.Repeat(" ", pos.Col), strings.Repeat("^", length))
	fmt.Printf("[-]\n\n")
}