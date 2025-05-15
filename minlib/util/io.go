package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
		return nil, err
	}
	
	return ReadBytecode(*bytes.NewBuffer(data))
}

func ReadBytecode(buffer bytes.Buffer) ([]byte, error) {
	data := buffer.Bytes()

	if len(data) < 8 {
		return nil, fmt.Errorf("File is too small to contain header and checksum")
	}

	// check header "MNML"
	header := binary.LittleEndian.Uint32(data[:4])
	expectedHeader := binary.LittleEndian.Uint32([]byte(BYTECODE_HEADER))

	if header != expectedHeader {
		return nil, fmt.Errorf("Invalid bytecode header")
	}

	// check checksum
	content := data[:len(data) - 4]

	checksum := binary.LittleEndian.Uint32(data[len(data) - 4:])
	expectedChecksum := computeChecksum(content)

	if checksum != expectedChecksum {
		return nil, fmt.Errorf("Incorrect checksum")
	}

	return data[4:len(data) - 4], nil
}

// writes 'data' into 'out'.
func writeOutput(out io.Writer, data []byte) error {
	_, err := out.Write(data)
	return err
}

// checks 'path' and returns a writer to either stdin, stdout or a file.
// the caller must close the stream.
func getOutputChannel(path string) (io.WriteCloser, error) {
	var out io.WriteCloser

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
			
			out = file
		}
	}

	return out, nil
}

// writes to the output file, which also treats '*stdout' as a special value for stdout, and '*stderr' for stderr.
func WriteOutputFile(path string, data []byte) error {
	out, err := getOutputChannel(path)
	defer out.Close()

	if err != nil {
		return err
	}

	return writeOutput(out, data)
}

// writes to the output bytecode file, which also treats '*stdout' as a special value for stdout, and '*stderr' for stderr, and adds some specific things about the bytecode.
func WriteBytecodeFile(path string, data []byte) error {
	out, err := getOutputChannel(path)
	defer out.Close()

	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)

	buffer.Write([]byte(BYTECODE_HEADER))
	err = writeOutput(buffer, data)
	
	if err != nil {
		return err
	}

	_, err = buffer.Write(IntToBytes(int(computeChecksum(buffer.Bytes()))))
	
	if err != nil {
		return err
	}

	return writeOutput(out, buffer.Bytes())
}
