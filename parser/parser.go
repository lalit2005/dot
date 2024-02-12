package parser

import (
	"dot/ast"
	"dot/lexer"
	"dot/token"
	"strconv"
)

const (
	_ = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

type Parser struct {
	currentToken  token.Token
	peekToken     token.Token
	lexer         *lexer.Lexer
	errors        []string
	prefixParsers map[token.TokenType]prefixParser
	infixParsers  map[token.TokenType]infixParser
}

type (
	prefixParser func() ast.Expression
	infixParser  func(ast.Expression) ast.Expression
)

var priority = map[token.TokenType]int{
	token.EQUAL:     EQUALS,
	token.NOT_EQUAL: EQUALS,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.SLASH:     PRODUCT,
	token.ASTERISK:  PRODUCT,
	token.LPAREN:    CALL,
	token.LBRACKET:  INDEX,
	token.BANG:      PREFIX,
}

func NewParser(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:         lexer,
		currentToken:  lexer.NextToken(),
		peekToken:     lexer.NextToken(),
		errors:        []string{},
		prefixParsers: make(map[token.TokenType]prefixParser),
		infixParsers:  make(map[token.TokenType]infixParser),
	}

	// if the first token is a comment, skip it
	if parser.currentToken.Type == token.COMMENT {
		parser.currentToken = parser.peekToken
		parser.peekToken = lexer.NextToken()
	}

	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.STRING, parser.parseString)
	parser.registerPrefix(token.INTEGER, parser.parseInteger)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunction)
	parser.registerPrefix(token.LBRACKET, parser.parseArrayLiteral)

	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQUAL, parser.parseInfixExpression)
	parser.registerInfix(token.EQUAL, parser.parseInfixExpression)
	parser.registerInfix(token.LPAREN, parser.parseCallExpression)
	parser.registerInfix(token.LBRACKET, parser.parseIndexExpression)

	return parser
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}
	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		// this sets current token to the first token of the next statement
		if p.currentToken.Type == token.COMMENT {
			p.nextToken()
		}
		program.Statements = append(program.Statements, statement)
	}
	return program
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Operator: p.currentToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// current token: first token of expression
	prefix := p.prefixParsers[p.currentToken.Type]
	if prefix == nil {
		p.newError("no prefix parser for '" + string(p.currentToken.Type) + "'")
		return nil
	}
	leftExp := prefix()
	for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		p.nextToken()
		infix := p.infixParsers[p.currentToken.Type]
		if infix == nil {
			return leftExp
		}
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseStatement() ast.Statement {
	// current token: first token of statement
	// the current token after each statement is passed goes to next statement's first token
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	if p.currentToken.Type != token.RETURN {
		p.newError("expected 'return'")
		return nil
	}
	p.nextToken()
	expr := &ast.ReturnStatement{
		ReturnValue: p.parseExpression(LOWEST),
	}
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	p.nextToken()
	return expr
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	expression := p.parseExpression(LOWEST)
	p.nextToken()
	if p.currentToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	// current token: first token of next statement
	return &ast.ExpressionStatement{Expression: expression}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	p.nextToken()
	if p.currentToken.Type != token.IDENTIFIER {
		p.newError("expected identifier after 'let'")
		return nil
	}
	identifier := ast.Identifier{Value: p.currentToken.Literal}
	p.nextToken()
	if p.currentToken.Type != token.ASSIGN {
		p.newError("expected '=' after identifier")
		return nil
	}
	p.nextToken()
	value := p.parseExpression(LOWEST)
	p.nextToken()
	if p.currentToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return &ast.LetStatement{Identifier: identifier, Value: value}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Value: p.currentToken.Literal}
}

func (p *Parser) parseInteger() ast.Expression {
	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		p.newError("could not parse '" + p.currentToken.Literal + "' as integer")
		return nil
	}
	return &ast.Integer{Value: value}
}

func (p *Parser) parseString() ast.Expression {
	return &ast.String{Value: p.currentToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Value: p.currentToken.Type == token.TRUE}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// current token: operator
	expression := &ast.InfixExpression{
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	// current token: right expression's first token
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// current token: '('
	p.nextToken()
	expression := p.parseExpression(LOWEST)
	if p.peekToken.Type != token.RPAREN {
		p.newError("expected ')'")
		return nil
	}
	p.nextToken()
	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	// current token: 'if'
	expression := &ast.IfExpression{}
	p.nextToken()
	if p.currentToken.Type != token.LPAREN {
		p.newError("expected '('")
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	p.nextToken()
	// current token: )
	if p.currentToken.Type != token.RPAREN {
		p.newError("expected ')'")
		return nil
	}
	p.nextToken()
	// current token: {
	if p.currentToken.Type != token.LBRACE {
		p.newError("expected '{'")
		return nil
	}
	p.nextToken()
	expression.Consequence = p.parseBlockStatement()
	if p.currentToken.Type != token.RBRACE {
		p.newError("expected '}'")
		return nil
	}
	if p.peekToken.Type == token.ELSE {
		p.nextToken()
		p.nextToken()
		if p.currentToken.Type != token.LBRACE {
			p.newError("expected '{'")
			return nil
		}
		p.nextToken()
		expression.Alternative = p.parseBlockStatement()
		if p.currentToken.Type != token.RBRACE {
			p.newError("expected '}'")
			return nil
		}
	}
	return expression
}

// after this function, current token is the first token of the next statement
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// current token: first token of first statemetent inside block
	block := &ast.BlockStatement{
		Statements: []ast.Statement{},
	}
	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		block.Statements = append(block.Statements, statement)
	}
	// current token is semicolon as p.parseStatement() calls p.nextToken() at the end
	if p.currentToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	// current token: first token of next statement
	return block
}

func (p *Parser) parseFunction() ast.Expression {
	// current token: 'fn'
	function := &ast.Function{
		Parameters: []*ast.Identifier{},
	}
	p.nextToken()
	if p.currentToken.Type != token.LPAREN {
		p.newError("expected '('")
		return nil
	}
	p.nextToken()
	function.Parameters = p.parseFunctionParameters()
	p.nextToken()
	if p.currentToken.Type != token.LBRACE {
		p.newError("expected '{'")
		return nil
	}
	p.nextToken()
	function.Body = p.parseBlockStatement()
	return function
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	// current token: first parameter
	parseFunctionParameters := []*ast.Identifier{}
	for p.currentToken.Type != token.RPAREN {
		if p.currentToken.Type == token.IDENTIFIER {
			identifier := &ast.Identifier{Value: p.currentToken.Literal}
			parseFunctionParameters = append(parseFunctionParameters, identifier)
		}
		p.nextToken()
		// current token: ',' or ')'
		if p.currentToken.Type == token.COMMA {
			p.nextToken()
		}
	}
	return parseFunctionParameters
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	// current token: '('
	call := &ast.CallExpression{
		Function: function,
	}
	p.nextToken()
	call.Arguments = p.parseCallArguments()
	return call
}

func (p *Parser) parseCallArguments() []ast.Expression {
	// current token: first argument
	arguments := []ast.Expression{}
	for p.currentToken.Type != token.RPAREN {
		argument := p.parseExpression(LOWEST)
		arguments = append(arguments, argument)
		p.nextToken()
		// current token: ',' or ')'
		if p.currentToken.Type == token.COMMA {
			p.nextToken()
		}
	}
	return arguments
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	// current token: '['
	array := &ast.ArrayLiteral{
		Elements: []ast.Expression{},
	}
	for p.peekToken.Type != token.RBRACKET {
		p.nextToken()
		element := p.parseExpression(LOWEST)
		array.Elements = append(array.Elements, element)
		if p.peekToken.Type == token.COMMA {
			p.nextToken()
		}
	}
	p.nextToken()
	return array
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	// current token: '['
	index := &ast.IndexExpression{
		Left: left,
	}
	p.nextToken()
	index.Index = p.parseExpression(LOWEST)
	if p.peekToken.Type != token.RBRACKET {
		p.newError("expected ']'")
		return nil
	}
	p.nextToken()
	return index
}
