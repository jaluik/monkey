package evaluator

import (
	"fmt"
	"jaluik.com/monkey/ast"
	"jaluik.com/monkey/object"
	"jaluik.com/monkey/token"
)

func quote(node ast.Node, env *object.Environment) *object.Quote {
	node = evalUnquote(node, env)
	return &object.Quote{Node: node}
}

func evalUnquote(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}
		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}
		if len(call.Arguments) != 1 {
			return node
		}
		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToAstNode(unquoted)
	},
	)
}

func convertObjectToAstNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}
	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TURE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{
			Token: t,
			Value: obj.Value,
		}
	case *object.Quote:
		return obj.Node

	default:
		return nil
	}
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return callExpression.Function.TokenLiteral() == "unquote"
}
