package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i] = Modify(statement, modifier).(Statement)
		}
	case *ExpressionStatement:
		node.Expression = Modify(node.Expression, modifier).(Expression)
	}
	return modifier(node)
}
