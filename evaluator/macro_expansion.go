package evaluator

import (
	"jaluik.com/monkey/ast"
	"jaluik.com/monkey/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	var definitions []int

	for i, statement := range program.Statements {
		if isMacroDefinition(statement) {
			addMacro(statement, env)
			definitions = append(definitions, i)
		}
	}
	for i := len(definitions) - 1; i >= 0; i-- {
		definitionIndex := definitions[i]
		program.Statements = append(program.Statements[:definitionIndex], program.Statements[definitionIndex+1:]...)
	}

}

func addMacro(statement ast.Statement, env *object.Environment) {
	letStatement := statement.(*ast.LetStatement)
	macroLiteral := letStatement.Value.(*ast.MacroLiteral)
	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Body:       macroLiteral.Body,
		Env:        env,
	}
	env.Set(letStatement.Name.Value, macro)
}

func isMacroDefinition(node ast.Statement) bool {
	letStatement, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStatement.Value.(*ast.MacroLiteral)
	return ok
}
