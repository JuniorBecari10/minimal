package analyzer

import (
	"fmt"
	"minc/ast"
	"minc/diagnostic"
	"minc/types"
)

// pre-declares all top-level declarations in order to properly do name resolution.
// it doesn't do all the required setup, it just adds the name to the list with its type, using shallow inference.
// TODO: make the top-level not get analyzed twice
func (a *Analyzer) hoistTopLevel() AnalyzerResult {
	res := RES_OK

	for _, d := range a.ast {
		global, diag := a.analyzeTopLevelDecl(d); if diag != nil {
			diag.PrintDiagnostic()

			if diag.DiagnosticType() == diagnostic.TYPE_WARNING {
				res = RES_ERROR
			}

			continue
		}

		a.globals = append(a.globals, global)
	}

	return res
}

func (a *Analyzer) analyzeTopLevelDecl(d ast.Statement) (Global, diagnostic.Diagnostic) {
		switch decl := d.Data.(type) {
			// In 'fn' statements we check only the declaration. All types must be explicitly annotated and concrete.
			case ast.FnDeclaration: {
				global, diag := a.topLevelFnDecl(decl); if diag != nil {
					return Global{}, diag
				}

				return global, nil
			}
			
			// In 'var'/'let' declarations we check the type and if omitted, we try to infer it shallowly.
			case ast.VarDeclaration: {
				global, diag := a.topLevelVarDecl(decl); if diag != nil {
					return Global{}, diag
				}

				return global, nil
			}

			// In records, all types must be explicitly annotated and concrete.
			case ast.RecordDeclaration: {
				// not for now. this is a dummy error.
				return Global{}, a.makeTypeAnnotationsNeeded(decl.Name)
			}

			// Should not reach here.
			default:
				panic(fmt.Sprintf("Unknown declaration %#v of type %T", decl, decl))
		}
}

// Adds native functions and variables to the global scope.
// These can have dummy types because they won't go in diagnostics.
func (a *Analyzer) addNatives() {
	// fn print()
	a.globals = append(a.globals, newNative("print", types.TypeFunction{
		Parameters: []types.Type{},
		Return: types.DummyType(types.TypeVoid{}),
	}))
	
	// fn println()
	a.globals = append(a.globals, newNative("println", types.TypeFunction{
		Parameters: []types.Type{},
		Return: types.DummyType(types.TypeVoid{}),
	}))
	
	// fn input(prompt: str): str
	a.globals = append(a.globals, newNative("print", types.TypeFunction{
		Parameters: []types.Type{ types.DummyType(types.TypeStr{}) },
		Return: types.DummyType(types.TypeStr{}),
	}))
	
	// fn time(): int
	a.globals = append(a.globals, newNative("time", types.TypeFunction{
		Parameters: []types.Type{},
		Return: types.DummyType(types.TypeVoid{}),
	}))
}

func (a *Analyzer) topLevelFnDecl(decl ast.FnDeclaration) (Global, diagnostic.Diagnostic) {
	// Check the return type. Maybe extract this in a different and reusable function.

	// the token can be a dummy one, since this won't be printed in the diagnostic, since void is concrete.
	// at least inside this one.
	returnType := types.DummyType(types.TypeVoid{})

	if decl.Return != nil {
		returnType = *decl.Return
	}

	if !typeIsConcrete(returnType.Data) {
		return Global{}, a.makeExpectedConcreteType(returnType)
	}
	
	paramTypes := []types.Type{}

	for _, param := range decl.Parameters {
		if param.Type == nil {
			return Global{}, a.makeTypeAnnotationsNeeded(param.Name)
		}

		paramTypes = append(paramTypes, *param.Type)
	}

	return Global{
		name: decl.Name,
		globalType: types.DummyType(types.TypeFunction{
			Parameters: paramTypes,
			Return: returnType,
		}),

		immutable: true,
		initialized: false,
		modified: false,
	}, nil
}

func (a *Analyzer) topLevelVarDecl(decl ast.VarDeclaration) (Global, diagnostic.Diagnostic) {
	if decl.Type == nil {
		// type isn't annotated. infer it shallowly.
		expr, diag := a.analyzeExpression(decl.Init, true); if diag != nil {
			return Global{}, diag
		}

		// type of expr will be the type of the variable, if applicable (checked later).
		// must be concrete; otherwise, it will require a type annotation.

		if !typeIsConcrete(expr.Data.Type()) {
			return Global{}, a.makeTypeAnnotationsNeeded(decl.Name)
		}
		
		// type must be dummy because it is inferred and therefore not in the source code.
		return Global{
			name: decl.Name,
			globalType: types.DummyType(expr.Data.Type()),

			immutable: decl.Immutable,
			initialized: false,
			modified: false,
		}, nil
	} else {
		// type is annotated; add the variable with its type.
		// actual type checking is done later.
		return Global{
			name: decl.Name,
			globalType: *decl.Type,

			immutable: decl.Immutable,
			initialized: false,
			modified: false,
		}, nil
	}
}
