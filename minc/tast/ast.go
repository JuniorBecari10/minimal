package tast

import (
	"minc/ast"
	"minc/types"
	"minlib/token"
)

type Tast = []Statement

type AstBase ast.AstBase
type Parameter = ast.Parameter
type Field = ast.Field

// returned in resolvedVariable in the analyzer, containing the fields that Local and Global have in common.
type Variable struct {
	Name token.Token
	VarType types.Type

	Immutable bool
	Modified *bool // pointer to the original data for mutability
	Used *bool
}
