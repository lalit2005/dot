package parser

import (
	"dot/ast"
	"dot/lexer"
	"dot/token"
	"log"
)

const (
	_ = iota
	LOWEST
	EQUALS
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

type Parser struct {
	currentToken token.Token
	peekToken    token.Token
	lexer        *lexer.Lexer
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt == nil {
			return nil
		}
		program.Statements = append(program.Statements, stmt)
		l := log.Default()
		l.Printf("STMT:: %+v", stmt)
	}
	return program
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:        lexer,
		currentToken: lexer.NextToken(),
		peekToken:    lexer.NextToken(),
	}
	return p
}

func (p *Parser) nextToken() token.Token {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
	return p.currentToken
}
