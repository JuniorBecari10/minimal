package util

import (
	"encoding/binary"
	"fmt"
	"strconv"
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

// This modifies in-place.
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func Remove[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

// Just like 'map'
func CopyList[T any](list []T, copier func(T) T) []T {
	res := make([]T, 0, len(list)) // pre-allocate with enough capacity

	for _, e := range list {
		res = append(res, copier(e))
	}

	return res
}

// Removes the last element from the supplied list and returns it.
func PopList[T any](list *[]T) T {
	lastIndex := len(*list) - 1
	
	topElement := (*list)[lastIndex]
	*list = (*list)[:lastIndex]

	return topElement
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

func Error(pos token.Position, length int, message string, fileData *FileData) {
	fmt.Printf("[-] Error at %s (%d, %d): %s\n", fileData.Name, pos.Line+1, pos.Col+1, message)
	fmt.Printf(" |  %d | %s\n", pos.Line+1, fileData.Lines[pos.Line])
	fmt.Printf(" | %s    %s%s\n", strings.Repeat(" ", len(strconv.Itoa(pos.Line+1))), strings.Repeat(" ", pos.Col), strings.Repeat("^", length))
	fmt.Printf("[-]\n\n")
}
