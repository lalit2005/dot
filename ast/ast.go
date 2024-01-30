package ast

import (
	"dot/token"
)

type Node interface {
	// String()
}

type Expression interface {
	expressionNode()
	Node
}

type Statement interface {
	statementNode()
	Node
}

type ExpressionStatement interface {
	expressionStatement()
	Node
}

type Program struct {
	Statements []Statement
}

type Integer struct {
	Token token.Token
	Value int64
}

func (i *Integer) expressionNode() {}

type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}

type Boolean struct {
	Value bool
}

func (i *Boolean) expressionNode() {}

type LetStatement struct {
	Identifier Identifier
	Value      Expression
}

func (l *LetStatement) statementNode() {}

type ReturnStatement struct {
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}

type PrefixExpression struct {
	Operator token.Token
	Right    Expression
}

func (i PrefixExpression) expressionNode() {}
