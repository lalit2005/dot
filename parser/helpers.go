package parser

import (
	"dot/token"
)

func (p *Parser) nextToken() token.Token {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
	return p.currentToken
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParser) {
	p.prefixParsers[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParser) {
	p.infixParsers[tokenType] = fn
}

func (p *Parser) newError(error string) {
	p.errors = append(p.errors, error)
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := priority[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if precedence, ok := priority[p.currentToken.Type]; ok {
		return precedence
	}
	return LOWEST
}
