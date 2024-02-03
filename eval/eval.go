package eval

import (
	"dot/ast"
)

func Eval(node ast.Node, env *Environment) interface{} {
	switch node := node.(type) {
	case *ast.Integer:
		return node.Value
	case *ast.Identifier:
		val, ok := env.Get(node.Value)
		if !ok {
			return nil
		}
		return val
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Boolean:
		return node.Value
	case *ast.String:
		return node.Value
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if val == nil {
			return nil
		}
		env.Set(node.Identifier.Value, val)
	case *ast.PrefixExpression:
		switch node.Operator {
		case "!":
			val, ok := Eval(node.Right, env).(bool)
			if !ok {
				return nil
			}
			return !val
		case "-":
			val, ok := Eval(node.Right, env).(int64)
			if !ok {
				return nil
			}
			return -val
		}
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		if left == nil || right == nil {
			return nil
		}
		switch {
		case left == nil || right == nil:
			return nil
		case node.Operator == "==":
			return left == right
		case node.Operator == "!=":
			return left != right
		case left != nil && right != nil:
			switch node.Operator {
			case "+":
				return left.(int64) + right.(int64)
			case "-":
				return left.(int64) - right.(int64)
			case "*":
				return left.(int64) * right.(int64)
			case "/":
				return left.(int64) / right.(int64)
			case ">":
				return left.(int64) > right.(int64)
			case "<":
				return left.(int64) < right.(int64)
			default:
				return nil
			}
		}
	case *ast.IfExpression:
		condition := Eval(node.Condition, env)
		if condition == nil {
			return nil
		}
		if condition.(bool) {
			return Eval(node.Consequence, env)
		}
		if node.Alternative != nil {
			return Eval(node.Alternative, env)
		}

	case *ast.BlockStatement:
		var result interface{}
		for _, statement := range node.Statements {
			result =
				Eval(statement, env)
		}
		return result
	case *ast.Program:
		var result interface{}
		for _, statement := range node.Statements {
			result =
				Eval(statement, env)
		}
		return result
	}
	return nil
}

type Environment struct {
	store map[string]interface{}
	Outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]interface{}), Outer: nil}
}

func (e *Environment) Get(name string) (interface{}, bool) {
	val, ok := e.store[name]
	return val, ok
}

func (e *Environment) Set(name string, val interface{}) interface{} {
	e.store[name] = val
	return val
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = outer
	return env
}
