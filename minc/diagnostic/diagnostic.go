package diagnostic

import (
	"fmt"
	"minlib/file"
	"os"
	"strconv"
	"strings"
)

type Diagnostic interface {
	PrintDiagnostic()
}

type DiagnosticBase struct {
	Message string
	Span Span
	FileData *file.FileData
}

// ---

type SimpleDiagnostic struct {
	DiagnosticBase
}

type HelpDiagnostic struct {
	DiagnosticBase
	Help []string
}

// ---

func (b *DiagnosticBase) PrintDiagnostic() {
	lineNum := int(b.Span.Pos.Line + 1)
	colNum := int(b.Span.Pos.Col + 1)
	
	lineStr := strconv.Itoa(lineNum)
	padding := strings.Repeat(" ", len(lineStr))
	carets := strings.Repeat("^", b.Span.Length)

	eprintln("")
	eprintf("[-] Error: %s\n", b.Message)
	eprintf(" | %s [-] %s (line %d, col %d)\n", padding, b.FileData.Name, lineNum, colNum)
	eprintf(" |  %d | %s\n", lineNum, b.FileData.Lines[b.Span.Pos.Line])
	eprintf(" | %s  | %s%s\n", padding, strings.Repeat(" ", int(b.Span.Pos.Col)), carets)
	eprintf(" | %s [-]\n", padding)
}

func (s *SimpleDiagnostic) PrintDiagnostic() {
	s.DiagnosticBase.PrintDiagnostic()
	eprintln("[-]")
}

func (h *HelpDiagnostic) PrintDiagnostic() {
	padding := strings.Repeat(" ", len(strconv.Itoa(int(h.Span.Pos.Line + 1))))

	h.DiagnosticBase.PrintDiagnostic()
	eprintf(" | %s [-] Help\n", padding)
	
	for _, line := range h.Help {
		eprintf(" | %s  | %s\n", padding, line)
	}

	eprintf(" | %s [-]\n", padding)
	eprintln("[-]")
}

// ---

func eprintln(s string) {
	fmt.Fprintln(os.Stderr, s)
}

func eprintf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}
