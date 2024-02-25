package parser

import (
	"dot/token"
	"fmt"
)

func (p *Parser) nextToken() token.Token {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
	if p.currentToken.Type == token.COMMENT {
		p.currentToken = p.peekToken
		p.peekToken = p.lexer.NextToken()
	}
	return p.currentToken
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParser) {
	p.prefixParsers[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParser) {
	p.infixParsers[tokenType] = fn
}

func (p *Parser) newError(error string, line int, column int) {
	p.errors = append(p.errors, error+" - "+fmt.Sprintf("at line %d, column %d", line, column))
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

func (p *Parser) PrintErrors() {
	if len(p.errors) > 0 {
		for _, e := range p.errors {
			fmt.Printf("PARSER ERROR: %s\n", e)
		}
	}
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.newError(fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type), p.lexer.Line(), p.lexer.Column())
	return false
}
