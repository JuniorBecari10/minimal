package start

import (
	"fmt"
	"minc/analyzer"
	"minc/parser"
	"os"
	"strings"

	"minlib/file"
	"minlib/value"
)

func Compile(sourcePath, outputPath string) {
	sourceContent, err := file.ReadSourceFile(sourcePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading source file '%s'.\n", sourcePath)
		fmt.Fprintln(os.Stderr, "Please check if the name is typed correctly.")
		fmt.Fprintf(os.Stderr, "\nError: '%s'.\n", err.Error())
		os.Exit(1)
	}

	source := string(sourceContent)

	fileData := file.FileData{
		Name:  file.GetFileName(sourcePath),
		Lines: strings.Split(source, "\n"),
	}

	chunk, hadError := compileSource(source, &fileData)
	
	if hadError {
		os.Exit(1)
	}

	if err := file.WriteBytecodeFile(outputPath, chunk.Serialize()); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to output file '%s': %v\n", outputPath, err)
		fmt.Fprintf(os.Stderr, "Error: '%s'.\n", err.Error())
		os.Exit(1)
	}
}

func compileSource(source string, fileData *file.FileData) (value.Chunk, bool) {
	ast, res := parser.New(source, fileData).Parse()
	
	if res == parser.RES_ERROR {
		os.Exit(1)
	}

	// fmt.Printf("%#v\n\n", ast)

	tast, res := analyzer.New(ast, fileData).Analyze()
	
	if res == analyzer.RES_ERROR {
		os.Exit(1)
	}

	fmt.Printf("%#v\n", tast)
	// fmt.Sprintln(tast)

	return value.Chunk{}, true
}

