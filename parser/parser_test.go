package parser

import (
	"dot/ast"
	"dot/lexer"
	"testing"
)

func newParser(input string) *Parser {
	l := lexer.NewLexer(input)
	p := NewParser(l)
	return p
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		// {"let y = true;", "y", true},
		{"let x = 5;", "x", 5},
		// {"let foobar = y;", "foobar", "y"},
	}
	for _, tt := range tests {
		p := newParser(tt.input)
		stmts := p.ParseProgram()
		for _, e := range p.errors {
			t.Error("PARSER ERROR: " + e)
		}
		t.Logf("%+v", stmts)
		stmt := stmts.Statements[0].(*ast.LetStatement)
		if stmt.Identifier.Value != tt.expectedIdentifier {
			t.Fatalf("wrong identifier. got=%q. want=%q", stmt.Identifier.Value, tt.expectedIdentifier)
		}
		if stmt.Value != tt.expectedValue {
			t.Fatalf("wrong value. got=%q. want=%q", stmt.Value, tt.expectedValue)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		// {"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		p := newParser(tt.input)
		program := p.ParseProgram()
		for _, e := range p.errors {
			t.Error("PARSER ERROR: " + e)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got=%T", stmt)
		}
		if returnStmt.ReturnValue != tt.expectedValue {
			t.Fatalf("returnStmt.ReturnValue not %q, got %q", tt.expectedValue,
				returnStmt.ReturnValue)
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		// {"!foobar;", "!", "foobar"},
		// {"-foobar;", "-", "foobar"},
		// {"!true;", "!", true},
		// {"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		p := newParser(tt.input)
		program := p.ParseProgram()
		for _, e := range p.errors {
			t.Error("PARSER ERROR: " + e)
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
	}
}
