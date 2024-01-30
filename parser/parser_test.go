package parser

import (
	"dot/ast"
	"dot/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let y = true;", "y", true},
		{"let x = 5;", "x", 5},
		{"let foobar = y;", "foobar", "y"},
	}
	for _, tt := range tests {
		p := newParser(tt.input)
		stmts := p.ParseProgram()
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

func newParser(input string) *Parser {
	l := lexer.NewLexer(input)
	p := NewParser(l)
	return p
}
