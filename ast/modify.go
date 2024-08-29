package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i] = Modify(statement, modifier).(Statement)
		}
	case *InfixExpression:
		node.Left = Modify(node.Left, modifier).(Expression)
		node.Right = Modify(node.Right, modifier).(Expression)
	case *PrefixExpression:
		node.Right = Modify(node.Right, modifier).(Expression)
	case *IndexExpression:
		node.Left = Modify(node.Left, modifier).(Expression)
		node.Index = Modify(node.Index, modifier).(Expression)
	case *ExpressionStatement:
		node.Expression = Modify(node.Expression, modifier).(Expression)
	}
	return modifier(node)
}
