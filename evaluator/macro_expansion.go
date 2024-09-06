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

func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}
		macro, ok := isMacroCall(callExpression, env)
		if !ok {
			return node
		}
		args := quoteArgs(callExpression)
		evalEnv := extendMacroEnv(macro, args)

		evaluated := Eval(macro.Body, evalEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-nodes from macros")
		}
		return quote.Node
	})
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	extended := object.NewEnclosedEnvironment(macro.Env)
	for paramIdx, param := range macro.Parameters {
		extended.Set(param.Value, args[paramIdx])
	}
	return extended
}

func quoteArgs(expression *ast.CallExpression) []*object.Quote {
	var args []*object.Quote
	for _, a := range expression.Arguments {
		args = append(args, &object.Quote{Node: a})
	}
	return args
}

func isMacroCall(exp *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	identifier, ok := exp.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}
	obj, ok := env.Get(identifier.Value)
	if !ok {
		return nil, false
	}
	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}
	return macro, true

}

func isMacroDefinition(node ast.Statement) bool {
	letStatement, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStatement.Value.(*ast.MacroLiteral)
	return ok
}
