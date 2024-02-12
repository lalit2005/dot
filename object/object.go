package object

import (
	"dot/ast"
	"fmt"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	STRING_OBJ       = "STRING"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	ERROR_OBJ        = "ERROR"
	ARRAY_OBJ        = "ARRAY"
)

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) String() string   { return "ERROR: " + e.Message }

type Object interface {
	Type() ObjectType
	String() string
}

type Integer struct {
	Value float64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) String() string {
	return fmt.Sprintf("%g", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

type Null struct{}

func (n *Null) String() string { return NULL_OBJ }

type String struct{ Value string }

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) String() string {
	return s.Value
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

func (rv *ReturnValue) String() string {
	return rv.Value.String()
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

func (f *Function) String() string {
	return "fn"
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }

func (a *Array) String() string {
	var out string
	for i, e := range a.Elements {
		if i == 0 {
			out += e.String()
		} else {
			out += ", " + e.String()
		}
	}
	return "[" + out + "]"
}
