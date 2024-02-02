package parser

import (
	"dot/ast"
	"dot/lexer"
	"strings"
	"testing"
)

func newParser(input string) *Parser {
	l := lexer.NewLexer(input)
	p := NewParser(l)
	return p
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		// TODO: add test for string literal
		// if strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`) {
		// 	return testStringLiteral(t, exp, v)
		// }
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

// TODO: add test for string literal
// func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
// 	str, ok := exp.(*ast.String)
// 	if !ok {
// 		t.Errorf("exp not *ast.String. got=%T", exp)
// 		return false
// 	}

// 	if str.Value != strings.Trim(value, `"`) {
// 		t.Errorf("str.Value not %s. got=%s", value, str.Value)
// 		return false
// 	}

// 	return true
// }

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.Integer)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	return true
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let y = true;", "y", true},
		{"let x = 5;", "x", 5},
		{"let foobar = y;", "foobar", "y"},
		//TODO: add test for string literal
		// {`let foobar = "hello world";`, "foobar", `hello world`},
	}
	for _, tt := range tests {
		p := newParser(tt.input)
		stmts := p.ParseProgram()
		for _, e := range p.errors {
			t.Error("PARSER ERROR: " + e)
		}
		stmt := stmts.Statements[0].(*ast.LetStatement)
		if stmt.Identifier.Value != tt.expectedIdentifier {
			t.Fatalf("wrong identifier. got=%+v. want=%+v", stmt.Identifier.Value, tt.expectedIdentifier)
		}
		if !testLiteralExpression(t, stmt.Value, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
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
		if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
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

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
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

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.InfixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// 0
		{
			"-a * b;",
			"((-a) * b)",
		},
		// 1
		{
			"!-a;",
			"(!(-a))",
		},
		// 2
		{
			"a + b + c;",
			"((a + b) + c)",
		},
		// 3
		{
			"a + b - c;",
			"((a + b) - c)",
		},
		// 4
		{
			"a * b * c;",
			"((a * b) * c)",
		},
		// 5
		{
			"a * b / c;",
			"((a * b) / c)",
		},
		// 6
		{
			"a + b / c;",
			"(a + (b / c))",
		},
		// 7
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		// 8
		{
			"3 + 4; -5 * 5;",
			"(3 + 4)\n((-5) * 5)",
		},
		// 9
		{
			"5 > 4 == 3 < 4;",
			"((5 > 4) == (3 < 4))",
		},
		// 10
		{
			"5 < 4 != 3 > 4;",
			"((5 < 4) != (3 > 4))",
		},
		// 11
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		// 12
		{
			"true;",
			"true",
		},
		// 13
		{
			"false;",
			"false",
		},
		// 14
		{
			"3 > 5 == false;",
			"((3 > 5) == false)",
		},
		// 15
		{
			"3 < 5 == true;",
			"((3 < 5) == true)",
		},
		// 16
		{
			"1 + (2 + 3) + 4;5+5",
			"((1 + (2 + 3)) + 4)\n(5 + 5)",
		},
		// 17
		{
			"(5 + 5) * 2;",
			"((5 + 5) * 2)",
		},
		// 18
		{
			"2 / (5 + 5);",
			"(2 / (5 + 5))",
		},
		// 19
		{
			"(5 + 5) * 2 * (5 + 5);",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		// 20
		{
			"-(5 + 5);",
			"(-(5 + 5))",
		},
		// 21
		{
			"!(true == true);",
			"(!(true == true))",
		},
		// {
		// 	"a + add(b * c) + d;",
		// 	"((a + add((b * c))) + d)",
		// },
		// {
		// 	"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8));",
		// 	"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		// },
		// {
		// 	"add(a + b + c * d / f + g);",
		// 	"add((((a + b) + ((c * d) / f)) + g))",
		// },
		// {
		// 	"a * [1, 2, 3, 4][b * c] * d;",
		// 	"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		// },
		// {
		// 	"add(a * b[2], b[1], 2 * [1, 2][1]);",
		// 	"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		// },
	}

	for i, tt := range tests {
		p := newParser(tt.input)
		program := p.ParseProgram()
		for _, e := range p.errors {
			t.Errorf("tests[%d] PARSER ERROR: "+e, i)
		}
		actual := program.String()
		if strings.TrimSpace(actual) != tt.expected {
			t.Errorf("tests[%d] expected=%q, got=%q", i, tt.expected, strings.TrimSpace(actual))
		} else {
			t.Logf("tests[%d] success", i)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"
	p := newParser(input)
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
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if exp.String() != "if ((x < y)) {\n  x\n}" {
		t.Errorf("exp.String() is not 'if ((x < y)) {\n  x\n}'. got=%q", exp.String())
	}
}
