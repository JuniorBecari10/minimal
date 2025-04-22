package util

import (
	"io"
	"os"
	"strings"
)

const (
	STDIN = "*stdin"
	STDOUT = "*stdout"
	STDERR = "*stderr"
)

// reads the source file, which also treats '*stdin' as a special value for stdin.
func ReadSourceFile(path string) ([]byte, error) {
	// check for special '*stdin'
	if path == STDIN {
		input, err := io.ReadAll(os.Stdin)
		
		if err != nil {
			return []byte{}, err
		}

		return input, nil
	}

	data, err := os.ReadFile(path)
	
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

// writes to the output file, which also treats '*stdout' as a special value for stdout, and '*stderr' for stderr.
func WriteOutputFile(path string, data []byte) error {
	var out io.Writer

	// check for special '*stdout' and '*stderr'
	switch strings.ToLower(path) {
		case "*stdout":
			out = os.Stdout
		
		case "*stderr":
			out = os.Stderr
		
		default: {
			file, err := os.Create(path)
			
			if err != nil {
				return err
			}
			
			defer file.Close()
			out = file
		}
	}

	_, err := out.Write(data)
	return err
}
