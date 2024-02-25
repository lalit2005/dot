package eval

import (
	"dot/ast"
	"dot/lexer"
	"dot/object"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment, lexer lexer.Lexer) object.Object {
	switch node := node.(type) {
	case *ast.Integer:
		return &object.Integer{Value: node.Value}
	case *ast.Identifier:
		if fn, ok := builtins[node.Value]; ok {
			return fn
		}
		val, ok := env.Get(node.Value)
		if !ok {
			return newError("identifier not found: "+node.Value, lexer.Line(), lexer.Column())
		}
		return val
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env, lexer)
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
	case *ast.String:
		return &object.String{Value: node.Value}
	case *ast.LetStatement:
		val := Eval(node.Value, env, lexer)
		if val == nil {
			return nil
		}
		env.Set(node.Identifier.Value, val)
		return val
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env, lexer)
		if val == nil {
			return nil
		}
		return &object.ReturnValue{Value: val}
	case *ast.PrefixExpression:
		switch node.Operator {
		case "!":
			right, ok := Eval(node.Right, env, lexer).(*object.Boolean)
			if !ok {
				return newError("invalid operation: "+node.String(), lexer.Line(), lexer.Column())
			}
			return &object.Boolean{Value: !right.Value}
		case "-":
			right, ok := Eval(node.Right, env, lexer).(*object.Integer)
			if !ok {
				return newError("invalid operation: "+node.String(), lexer.Line(), lexer.Column())
			}
			return &object.Integer{Value: -right.Value}
		default:
			return newError("unknown operator: "+node.Operator, lexer.Line(), lexer.Column())
		}
	case *ast.InfixExpression:
		// reassigning the value of an element in an array
		if node.Operator == "=" {
			left := node.Left.(*ast.IndexExpression)
			val := Eval(node.Right, env, lexer)
			if val == nil {
				return nil
			}
			arrayObj, ok := env.Get(left.Left.String())
			if !ok {
				return newError("identifier not found: "+left.Left.String(), lexer.Line(), lexer.Column())
			}
			array := arrayObj.(*object.Array)
			index := int(Eval(left.Index, env, lexer).(*object.Integer).Value)
			if index < 0 || index >= len(array.Elements) {
				return newError("index out of range", lexer.Line(), lexer.Column())
			}
			array.Elements[index] = val
			return val
		}
		left := Eval(node.Left, env, lexer)
		right := Eval(node.Right, env, lexer)
		if left == nil || right == nil {
			if left == nil {
				return newError("left operand is nil", lexer.Line(), lexer.Column())
			} else {
				return newError("right operand is nil", lexer.Line(), lexer.Column())
			}
		}
		switch {
		case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
			return evalIntegerInfixOperation(node.Operator, left, right, lexer)
		case node.Operator == "==":
			return getBooleanObject(left.String() == right.String())
		case node.Operator == "&&":
			if left.Type() != object.BOOLEAN_OBJ || right.Type() != object.BOOLEAN_OBJ {
				return newError(fmt.Sprintf("invalid operation: %s %s %s", left.Type(), node.Operator, right.Type()), lexer.Line(), lexer.Column())
			}
			return getBooleanObject(left.(*object.Boolean).Value && right.(*object.Boolean).Value)
		case node.Operator == "||":
			if left.Type() != object.BOOLEAN_OBJ || right.Type() != object.BOOLEAN_OBJ {
				return newError(fmt.Sprintf("invalid operation: %s %s %s", left.Type(), node.Operator, right.Type()), lexer.Line(), lexer.Column())
			}
			return getBooleanObject(left.(*object.Boolean).Value || right.(*object.Boolean).Value)
		case node.Operator == "!=":
			return getBooleanObject(left != right)
		case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
			if node.Operator != "+" {
				return newError(fmt.Sprintf("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type()), lexer.Line(), lexer.Column())
			}
			return &object.String{Value: left.(*object.String).Value + right.(*object.String).Value}
		case left.Type() != right.Type():
			return newError(fmt.Sprintf("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type()), lexer.Line(), lexer.Column())
		}
	case *ast.IfExpression:
		condition := Eval(node.Condition, env, lexer)
		if condition == nil {
			return nil
		}
		if condition.String() == "true" {
			return Eval(node.Consequence, env, lexer)
		} else if node.Alternative != nil {
			return Eval(node.Alternative, env, lexer)
		} else {
			return NULL
		}
	case *ast.Function:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env, lexer)
		if function.Type() == object.ERROR_OBJ {
			return function
		}
		args := evalExpressions(node.Arguments, env, lexer)
		if len(args) == 1 && args[0].Type() == object.ERROR_OBJ {
			return args[0]
		}
		return applyFunction(function, args, lexer)
	case *ast.BlockStatement:
		var result object.Object
		for _, statement := range node.Statements {
			result = Eval(statement, env, lexer)
			if result != nil {
				rt := result.Type()
				if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
					return result
				}
			}
		}
		return result
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env, lexer)
		if len(elements) == 1 && elements[0].Type() == object.ERROR_OBJ {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env, lexer)
		if left.Type() == object.ERROR_OBJ {
			return left
		}
		index := Eval(node.Index, env, lexer)
		if index.Type() == object.ERROR_OBJ {
			return index
		}
		return evalIndexExpression(left, index, lexer)
	case *ast.WhileStatement:
		condition := Eval(node.Condition, env, lexer)
		// log.Println(condition)
		if condition == nil {
			return nil
		}
		for condition.String() == "true" {
			result := Eval(node.Body, env, lexer)
			if result != nil {
				rt := result.Type()
				if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
					return result
				}
			}
			condition = Eval(node.Condition, env, lexer)
			if condition == nil {
				return nil
			}
		}
		return NULL
	case *ast.Program:
		var result object.Object
		for _, statement := range node.Statements {
			result = Eval(statement, env, lexer)
			switch result := result.(type) {
			case *object.ReturnValue:
				return result.Value
			case *object.Error:
				return result
			}
		}
		return result
	}
	return newError("unknown node type: "+node.String(), lexer.Line(), lexer.Column())
}

func newError(msg string, line int, column int) *object.Error {
	return &object.Error{Message: msg + " - " + fmt.Sprintf("at line %d, column %d", line, column)}
}

func evalIntegerInfixOperation(operator string, l object.Object, r object.Object, lexer lexer.Lexer) object.Object {
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
	case "<=":
		return getBooleanObject(left <= right)
	case ">=":
		return getBooleanObject(left >= right)
	case "==":
		return getBooleanObject(left == right)
	case "!=":
		return getBooleanObject(left != right)
	default:
		return newError(fmt.Sprintf("unknown operator: %s %s %s", l.Type(), operator, r.Type()), lexer.Line(), lexer.Column())
	}
}

func getBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func evalExpressions(exps []ast.Expression, env *object.Environment, lexer lexer.Lexer) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env, lexer)
		if evaluated == nil {
			return nil
		}
		if evaluated.Type() == object.ERROR_OBJ {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object, lexer lexer.Lexer) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv, lexer)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: "+string(fn.Type()), lexer.Line(), lexer.Column())
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

func evalIndexExpression(left object.Object, index object.Object, lexer lexer.Lexer) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError((fmt.Sprintf("index operator not supported: %s", left.Type())), lexer.Line(), lexer.Column())
	}
}

func evalArrayIndexExpression(array object.Object, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := float64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[int(idx)]
}
