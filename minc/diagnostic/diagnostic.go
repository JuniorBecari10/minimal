package diagnostic

import "minlib/file"

type Diagnostic interface {
	String() string
}

type DiagnosticBase struct {
	Message string
	Span Span
	FileData *file.FileData
}

type SimpleDiagnostic struct {
	DiagnosticBase
}

type HelpDiagnostic struct {
	DiagnosticBase
	Help string
}

