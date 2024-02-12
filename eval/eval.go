package eval

import (
	"dot/ast"
	"dot/object"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Integer:
		return &object.Integer{Value: node.Value}
	case *ast.Identifier:
		val, ok := env.Get(node.Value)
		if !ok {
			return newError("identifier not found: " + node.Value)
		}
		return val
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
	case *ast.String:
		return &object.String{Value: node.Value}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if val == nil {
			return nil
		}
		env.Set(node.Identifier.Value, val)
		return val
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if val == nil {
			return nil
		}
		return &object.ReturnValue{Value: val}
	case *ast.PrefixExpression:
		switch node.Operator {
		case "!":
			right, ok := Eval(node.Right, env).(*object.Boolean)
			if !ok {
				return newError("invalid operation: " + node.String())
			}
			return &object.Boolean{Value: !right.Value}
		case "-":
			right, ok := Eval(node.Right, env).(*object.Integer)
			if !ok {
				return newError("invalid operation: " + node.String())
			}
			return &object.Integer{Value: -right.Value}
		default:
			return newError("unknown operator: " + node.Operator)
		}
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		if left == nil || right == nil {
			if left == nil {
				return newError("left operand is nil")
			} else {
				return newError("right operand is nil")
			}
		}
		switch {
		case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
			return evalIntegerInfixOperation(node.Operator, left, right)
		case node.Operator == "==":
			return getBooleanObject(left == right)
		case node.Operator == "!=":
			return getBooleanObject(left != right)
		case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
			if node.Operator != "+" {
				return newError(fmt.Sprintf("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type()))
			}
			return &object.String{Value: left.(*object.String).Value + right.(*object.String).Value}
		case left.Type() != right.Type():
			return newError(fmt.Sprintf("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type()))
		}
	case *ast.IfExpression:
		condition := Eval(node.Condition, env)
		if condition == nil {
			return nil
		}
		if condition == TRUE {
			return Eval(node.Consequence, env)
		} else {
			return Eval(node.Alternative, env)
		}
	case *ast.Function:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if function.Type() == object.ERROR_OBJ {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && args[0].Type() == object.ERROR_OBJ {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.BlockStatement:
		var result object.Object
		for _, statement := range node.Statements {
			result =
				Eval(statement, env)
			if result != nil {
				rt := result.Type()
				if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
					return result
				}
			}
		}
		return result
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && elements[0].Type() == object.ERROR_OBJ {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.Program:
		var result object.Object
		for _, statement := range node.Statements {
			result = Eval(statement, env)
			switch result := result.(type) {
			case *object.ReturnValue:
				return result.Value
			case *object.Error:
				return result
			}
		}
		return result
	}
	return newError("unknown node type: " + node.String())
}

func newError(msg string) *object.Error {
	return &object.Error{Message: fmt.Sprint(msg)}
}

func evalIntegerInfixOperation(operator string, l object.Object, r object.Object) object.Object {
	left := l.(*object.Integer).Value
	right := r.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: left + right}
	case "-":
		return &object.Integer{Value: left - right}
	case "*":
		return &object.Integer{Value: left * right}
	case "/":
		return &object.Integer{Value: left / right}
	case "<":
		return getBooleanObject(left < right)
	case ">":
		return getBooleanObject(left > right)
	default:
		return newError(fmt.Sprintf("unknown operator: %s %s %s", l.Type(), operator, r.Type()))
	}
}

func getBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if evaluated == nil {
			return nil
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	default:
		return newError("not a function: " + string(fn.Type()))
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
