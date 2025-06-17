package analyzer

import (
	"minc/ast"
	"minc/diagnostic"
	"minc/tast"
	"minc/types"
)

func (a *Analyzer) fnDecl(
	d ast.Statement,
	decl ast.FnDeclaration,
	generatedTast *tast.Tast,
	newStmt func(tast.StmtData) tast.Statement,
	inferredType **types.TypeData,
	mode BlockAnalyzeMode,
) diagnostic.Diagnostic {

}
