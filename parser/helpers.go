package parser

import (
	"dot/token"
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParser) {
	p.prefixParsers[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParser) {
	p.infixParsers[tokenType] = fn
}

func (p *Parser) newError(error string) {
	p.errors = append(p.errors, error)
}
