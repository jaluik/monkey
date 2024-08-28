package evaluator

import (
	"jaluik.com/monkey/ast"
	"jaluik.com/monkey/object"
)

func quote(node ast.Node) *object.Quote {
	return &object.Quote{Node: node}
}
