package parser

import (
	"dot/ast"
	"dot/token"
	"strconv"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// current token is LET
	if p.peekToken.Type != token.IDENTIFIER {
		return nil
	}
	identToken := p.nextToken()
	if p.peekToken.Type != token.ASSIGN {
		return nil
	}
	// current token is the starting of expression
	p.nextToken()
	identNode := ast.Identifier{Value: identToken.Literal}
	valueNode := p.parseExpression(LOWEST)
	stmt := &ast.LetStatement{Identifier: identNode, Value: valueNode}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	p.nextToken()
	value := p.parseExpression(LOWEST)
	return &ast.ReturnStatement{
		ReturnValue: value,
	}
}

func (p *Parser) parseExpression(priority int) ast.Expression {
	switch p.currentToken.Type {
	case token.INTEGER:
		value, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
		if err != nil {
			return nil
		}
		return &ast.Integer{
			Token: p.currentToken,
			Value: value,
		}
	case token.BANG, token.MINUS, token.PLUS:
		return &ast.PrefixExpression{
			Operator: p.currentToken,
			Right:    p.parseExpression(PREFIX),
		}
	default:
		return nil
	}
}

// func parsePrefixExpression(operator token.TokenType) ast.Expression {

// }
