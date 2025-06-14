package diagnostic

import (
	"fmt"
	"minlib/file"
	"os"
	"strconv"
	"strings"
)

type DiagnosticType string

const (
	TYPE_ERROR DiagnosticType = "Error"
	TYPE_WARNING = "Warning"
)

type Diagnostic interface {
	PrintDiagnostic()
	diagnosticType() DiagnosticType
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

type WarningDiagnostic struct {
	DiagnosticBase
}

// ---

// does not implement Diagnostic.
func (b *DiagnosticBase) printDiagnostic(diagType DiagnosticType) {
	lineNum := int(b.Span.Pos.Line + 1)
	colNum := int(b.Span.Pos.Col + 1)
	
	lineStr := strconv.Itoa(lineNum)
	padding := strings.Repeat(" ", len(lineStr))
	carets := strings.Repeat("^", b.Span.Length)

	eprintln("")
	eprintf("[-] %s: %s\n", diagType, b.Message)
	eprintf(" | %s [-] %s (line %d, col %d)\n", padding, b.FileData.Name, lineNum, colNum)
	eprintf(" |  %d | %s\n", lineNum, b.FileData.Lines[b.Span.Pos.Line])
	eprintf(" | %s  | %s%s\n", padding, strings.Repeat(" ", int(b.Span.Pos.Col)), carets)
	eprintf(" | %s [-]\n", padding)
}

// ---

func (s SimpleDiagnostic) PrintDiagnostic() {
	s.DiagnosticBase.printDiagnostic(s.diagnosticType())
	eprintln("[-]")
}

func (s SimpleDiagnostic) diagnosticType() DiagnosticType {
	return TYPE_ERROR
}

// ---

func (h HelpDiagnostic) PrintDiagnostic() {
	padding := strings.Repeat(" ", len(strconv.Itoa(int(h.Span.Pos.Line + 1))))

	h.DiagnosticBase.printDiagnostic(h.diagnosticType())
	eprintf(" | %s [-] Help\n", padding)
	
	for _, line := range h.Help {
		eprintf(" | %s  | %s\n", padding, line)
	}

	eprintf(" | %s [-]\n", padding)
	eprintln("[-]")
}

func (s HelpDiagnostic) diagnosticType() DiagnosticType {
	return TYPE_ERROR
}

// ---

func (w WarningDiagnostic) PrintDiagnostic() {
	w.DiagnosticBase.printDiagnostic(w.diagnosticType())
	eprintln("[-]")
}

func (w WarningDiagnostic) diagnosticType() DiagnosticType {
	return TYPE_WARNING
}

// ---

func eprintln(s string) {
	fmt.Fprintln(os.Stderr, s)
}

func eprintf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}
