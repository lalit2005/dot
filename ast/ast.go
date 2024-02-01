package ast

type Node interface{}

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

type Program struct {
	Statements []Statement
}

type Integer struct {
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
	Operator string
	Right    Expression
}

func (i PrefixExpression) expressionNode() {}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (i InfixExpression) expressionNode() {}
