package file

import "path/filepath"

func GetFileName(path string) string {
	// Use filepath to normalize and extract the file name
	normalizedPath := filepath.FromSlash(path) // Converts / to \ or vice versa based on OS
	return filepath.Base(normalizedPath)
}
