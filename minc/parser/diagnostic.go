package parser

import "minc/diagnostic"

func (p *Parser) makeStatementsNotAllowedDiagnostic() diagnostic.SimpleDiagnostic {
	return diagnostic.SimpleDiagnostic{
		DiagnosticBase: diagnostic.DiagnosticBase{
			Message: "Statements are not allowed at top-level.",
			Span: diagnostic.Span{
				Pos: p.current.Pos,
				Length: len(p.current.Lexeme),
			},
			FileData: p.fileData,
		},
	}
}
