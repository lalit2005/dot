package parser

import (
	"dot/ast"
	"dot/lexer"
	"dot/token"
	"log"
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

	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.STRING, parser.parseString)
	parser.registerPrefix(token.INTEGER, parser.parseInteger)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfExpression)

	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQUAL, parser.parseInfixExpression)
	parser.registerInfix(token.EQUAL, parser.parseInfixExpression)

	return parser
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}
	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
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
	// p.nextToken()
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
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
		p.nextToken()
	}
	return &ast.LetStatement{Identifier: identifier, Value: value}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Value: p.currentToken.Literal}
}

func (p *Parser) parseInteger() ast.Expression {
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
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
	log.Printf("current token: %s", p.currentToken.Literal)
	expression.Consequence = p.parseBlockStatement()
	log.Printf("current token: %s", p.currentToken.Literal)
	if p.currentToken.Type != token.RBRACE {
		p.newError("expected '}'")
		return nil
	}
	// if p.currentToken.Type == token.SEMICOLON {
	// 	p.nextToken()
	// }
	return expression
}

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
	// p.nextToken()
	// current token: first token of next statement
	return block
}
