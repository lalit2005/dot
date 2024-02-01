package ast

import "fmt"

type Node interface {
	String() string
}

type Expression interface {
	expressionNode()
	Node
}

type Statement interface {
	statementNode()
	Node
}

type ExpressionStatement struct {
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) String() string {
	return e.Expression.String()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out string
	for _, s := range p.Statements {
		out += s.String()
	}
	return out
}

type Integer struct {
	Value int64
}

func (i *Integer) expressionNode() {}

func (i *Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string {
	return i.Value
}

type Boolean struct {
	Value bool
}

func (i *Boolean) expressionNode() {}

func (i *Boolean) String() string {
	return fmt.Sprintf("%t", i.Value)
}

type LetStatement struct {
	Identifier Identifier
	Value      Expression
}

func (l *LetStatement) statementNode() {}

func (l *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s;", l.Identifier.String(), l.Value.String())
}

type ReturnStatement struct {
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) String() string {
	return fmt.Sprintf("return %s;", r.ReturnValue.String())
}

type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (i PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", i.Operator, i.Right)
}

func (i PrefixExpression) expressionNode() {}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (i InfixExpression) expressionNode() {}

func (i InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Operator, i.Right.String())
}
