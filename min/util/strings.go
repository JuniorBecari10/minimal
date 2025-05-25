package util

import "strings"

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
