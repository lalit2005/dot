package eval

import (
	"dot/object"
	"fmt"
	"strconv"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)), 0, 0)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: float64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: float64(len(arg.Elements))}
			default:
				return newError(fmt.Sprintf("argument to `len` not supported, got %s", args[0].Type()), 0, 0)
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)), 0, 0)
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(fmt.Sprintf("argument to `first` must be ARRAY, got %s", args[0].Type()), 0, 0)
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)), 0, 0)
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(fmt.Sprintf("argument to `last` must be ARRAY, got %s", args[0].Type()), 0, 0)
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)), 0, 0)
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(fmt.Sprintf("argument to `rest` must be ARRAY, got %s", args[0].Type()), 0, 0)
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args)), 0, 0)
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(fmt.Sprintf("argument to `push` must be ARRAY, got %s", args[0].Type()), 0, 0)
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				println(arg.String())
			}
			return &object.String{Value: ""}
		},
	},
	"ask": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				print(arg.String())
			}
			var input string
			fmt.Scanln(&input)
			return &object.String{Value: input}
		},
	},
	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)), 0, 0)
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError(fmt.Sprintf("argument to `int` must be STRING, got %s", args[0].Type()), 0, 0)
			}

			str := args[0].(*object.String).Value
			value, err := strconv.Atoi(str)
			if err != nil {
				return newError(fmt.Sprintf("failed to convert string to integer: %s", err.Error()), 0, 0)
			}

			return &object.Integer{Value: float64(value)}
		},
	},
}
