package ast

import (
	"bytes"
	"fmt"
	"strings"
)

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
	return e.Expression.String() + ";\n"
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
	Value float64
}

func (i *Integer) expressionNode() {}

func (i *Integer) String() string {
	return fmt.Sprintf("%g", i.Value)
}

type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string {
	return i.Value
}

type String struct {
	Value string
}

func (i *String) expressionNode() {}

func (i *String) String() string {
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
	return fmt.Sprintf("let %s = %s;\n", l.Identifier.String(), l.Value.String())
}

type ReturnStatement struct {
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) String() string {
	return fmt.Sprintf("return %s;\n", r.ReturnValue.String())
}

type WhileStatement struct {
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileStatement) statementNode() {}

func (w *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while (")
	out.WriteString(w.Condition.String())
	out.WriteString(") {\n")
	out.WriteString(w.Body.String())
	out.WriteString("}")
	return out.String()
}

type ForStatement struct {
	Initializer Statement
	Condition   Expression
	Incrementer Statement
	Body        *BlockStatement
}

func (f *ForStatement) statementNode() {}

func (f *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for (")
	out.WriteString(f.Initializer.String())
	out.WriteString(f.Condition.String())
	out.WriteString("; ")
	out.WriteString(f.Incrementer.String())
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("}")
	return out.String()
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

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) statementNode() {}

func (b *BlockStatement) String() string {
	var out string
	for _, s := range b.Statements {
		out += "  " + s.String()
	}
	return out
}

type IfExpression struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(" (" + ie.Condition.String() + ") ")
	out.WriteString("{\n")
	out.WriteString(ie.Consequence.String())
	out.WriteString("}")
	if ie.Alternative != nil {
		out.WriteString(" else {\n")
		out.WriteString(ie.Alternative.String())
		out.WriteString("}")
	}
	return out.String()
}

type Function struct {
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *Function) expressionNode() {}

func (f *Function) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("}")
	return out.String()
}

type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type ArrayLiteral struct {
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type IndexExpression struct {
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left.String(), ie.Index.String())
}
