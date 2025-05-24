package parser

import (
	"minc/ast"
	"minc/diagnostic"
)

func (p *Parser) declaration(allowStats bool) (ast.Statement, diagnostic.Diagnostic) {
	if p.panicMode {
		p.synchronize()
	}
}
