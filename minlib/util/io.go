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

const (
	BYTECODE_HEADER = "MNML"
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

// reads the bytecode file, which also treats '*stdin' as a special value for stdin, and performs some checks.
func ReadBytecodeFile(path string) ([]byte, error) {
	data, err := ReadSourceFile(path)

	if err != nil {
		return []byte{}, err
	}

	// TODO: add checks

	return data, nil
}

// writes 'data' into 'out'.
func writeOutput(out io.Writer, data []byte) error {
	_, err := out.Write(data)
	return err
}

// checks 'path' and returns a writer to either stdin, stdout or a file.
func getOutputChannel(path string) (io.Writer, error) {
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
				return nil, err
			}
			
			// TODO: handle closed file
			defer file.Close()
			out = file
		}
	}

	return out, nil
}

// writes to the output file, which also treats '*stdout' as a special value for stdout, and '*stderr' for stderr.
func WriteOutputFile(path string, data []byte) error {
	out, err := getOutputChannel(path)

	if err != nil {
		return err
	}

	return writeOutput(out, data)
}

// writes to the output bytecode file, which also treats '*stdout' as a special value for stdout, and '*stderr' for stderr, and adds some specific things about the bytecode.
func WriteBytecodeFile(path string, data []byte) error {
	out, err := getOutputChannel(path)

	if err != nil {
		return err
	}

	// TODO: add header writing

	return writeOutput(out, data)
}
