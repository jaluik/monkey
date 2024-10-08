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
	case *IfExpression:
		node.Condition = Modify(node.Condition, modifier).(Expression)
		node.Consequence = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative = Modify(node.Alternative, modifier).(*BlockStatement)
		}
	case *BlockStatement:
		for i := range node.Statements {
			node.Statements[i] = Modify(node.Statements[i], modifier).(Statement)
		}
	case *ReturnStatement:
		node.ReturnValue = Modify(node.ReturnValue, modifier).(Expression)
	case *LetStatement:
		node.Value = Modify(node.Value, modifier).(Expression)
	case *FunctionLiteral:
		for i := range node.Parameters {
			node.Parameters[i] = Modify(node.Parameters[i], modifier).(*Identifier)
		}
		node.Body = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i := range node.Elements {
			node.Elements[i] = Modify(node.Elements[i], modifier).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, val := range node.Pairs {
			newKey := Modify(key, modifier).(Expression)
			newVal := Modify(val, modifier).(Expression)
			newPairs[newKey] = newVal
		}
		node.Pairs = newPairs

	case *ExpressionStatement:
		node.Expression = Modify(node.Expression, modifier).(Expression)
	}
	return modifier(node)
}
