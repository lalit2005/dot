package parser

import (
	"dot/ast"
	"dot/lexer"
	"dot/token"
	"fmt"
	"log"
	"strconv"
)

var l = log.Default()

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
	token.EQUAL:    EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.BANG:     PREFIX,
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:        lexer,
		currentToken: lexer.NextToken(),
		peekToken:    lexer.NextToken(),
	}
	p.prefixParsers = make(map[token.TokenType]prefixParser)
	p.registerPrefix(token.INTEGER, p.parseInteger)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParsers = make(map[token.TokenType]infixParser)
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}
	// l.Printf("stmt 0: %+v", program.Statements[0])
	// l.Printf("stmt 1: %+v", program.Statements[1])
	return program
}

func (p *Parser) nextToken() token.Token {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
	return p.currentToken
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		// l.Printf("parseExpressionStatement")
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expression := p.parseExpression(LOWEST)
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return &ast.ExpressionStatement{Expression: expression}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	p.nextToken()
	// the current token is expression's starting token
	value := p.parseExpression(LOWEST)
	return &ast.ReturnStatement{
		ReturnValue: value,
	}
}

func (p *Parser) parseInteger() ast.Expression {
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		p.newError(fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal))
		return nil
	}
	return &ast.Integer{
		Value: value,
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Value: p.currentToken.Literal}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// current token is LET
	if p.peekToken.Type != token.IDENTIFIER {
		p.newError(fmt.Sprintf("expected next token to be IDENTIFIER, got %s instead", p.peekToken.Type))
		return nil
	}
	identToken := p.nextToken()
	if p.peekToken.Type != token.ASSIGN {
		p.newError(fmt.Sprintf("expected next token to be ASSIGN, got %s instead", p.peekToken.Type))
		return nil
	}
	p.nextToken()
	p.nextToken()
	// current token is the starting of expression
	// l.Printf("parseLetStatement currentToken: %+v", p.currentToken)
	identNode := ast.Identifier{Value: identToken.Literal}
	valueNode := p.parseExpression(LOWEST)
	stmt := &ast.LetStatement{Identifier: identNode, Value: valueNode}
	return stmt
}

func (p *Parser) parseExpression(pr int) ast.Expression {
	prefix := p.prefixParsers[p.currentToken.Type]
	if prefix == nil {
		p.newError(fmt.Sprintf("no prefix parser for %s found", p.currentToken.Type))
		return nil
	}
	leftExpression := prefix()
	return leftExpression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Operator: p.currentToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}
