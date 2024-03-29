package lexer

import (
	"dot/token"
	"strings"
)

type Lexer struct {
	input           string
	currentPosition int
	peekPosition    int
	currentChar     byte
	peekChar        byte
	line            int
	column          int
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	var tok token.Token
	switch l.currentChar {
	case '+':
		if l.peekChar == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.PLUS_EQUAL, Literal: "+="}
		}
		tok = newToken(token.PLUS, l.currentChar)
	case '-':
		if l.peekChar == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.MINUS_EQUAL, Literal: "-="}
		}
		tok = newToken(token.MINUS, l.currentChar)
	case '/':
		if l.peekChar == '/' {
			initialPosition := l.currentPosition
			for l.currentChar != '\n' && l.currentChar != 0 {
				l.readChar()
			}
			tok = token.Token{Type: token.COMMENT, Literal: l.input[initialPosition:l.currentPosition]}
			return tok
		}
		if l.peekChar == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.DIV_EQUAL, Literal: "/="}
		}
		tok = newToken(token.SLASH, l.currentChar)
	case '*':
		if l.peekChar == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.MULT_EQUAL, Literal: "*="}
		}
		tok = newToken(token.ASTERISK, l.currentChar)
	case ';':
		tok = newToken(token.SEMICOLON, l.currentChar)
	case ',':
		tok = newToken(token.COMMA, l.currentChar)
	case '(':
		tok = newToken(token.LPAREN, l.currentChar)
	case ')':
		tok = newToken(token.RPAREN, l.currentChar)
	case '{':
		tok = newToken(token.LBRACE, l.currentChar)
	case '}':
		tok = newToken(token.RBRACE, l.currentChar)
	case '[':
		tok = newToken(token.LBRACKET, l.currentChar)
	case ']':
		tok = newToken(token.RBRACKET, l.currentChar)
	case ':':
		tok = newToken(token.COLON, l.currentChar)
	case '=':
		if l.peekChar == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.EQUAL, Literal: "=="}
		}
		tok = newToken(token.ASSIGN, l.currentChar)
	case '!':
		if l.peekChar == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.NOT_EQUAL, Literal: "!="}
		}
		tok = newToken(token.BANG, l.currentChar)
	case '<':
		if l.peekChar == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.LTE, Literal: "<="}
		}
		tok = newToken(token.LT, l.currentChar)
	case '>':
		if l.peekChar == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.GTE, Literal: ">="}
		}
		tok = newToken(token.GT, l.currentChar)
	case '&':
		if l.peekChar == '&' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.AND, Literal: "&&"}
		}
		tok = newToken(token.UNKNOWN, l.currentChar)
	case '|':
		if l.peekChar == '|' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.OR, Literal: "||"}
		}
	case '"', '\'':
		quoteType := l.currentChar
		l.readChar()
		initialPosition := l.currentPosition
		for {
			if l.currentChar == quoteType || l.currentChar == 0 {
				break
			}
			l.readChar()
		}
		tok = token.Token{Type: token.STRING, Literal: l.input[initialPosition:l.currentPosition]}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		previousChar := l.currentChar
		if isAlphabet(l.currentChar) {
			initialPosition := l.currentPosition
			for isAlphabet(l.currentChar) {
				l.readChar()
			}
			sequence := l.input[initialPosition:l.currentPosition]
			tokType, ok := token.Keywords[sequence]
			if !ok {
				return token.Token{Type: token.IDENTIFIER, Literal: sequence}
			} else {
				return token.Token{Type: tokType, Literal: sequence}
			}
		} else if isDigit(l.currentChar, previousChar, l.peekChar) {
			initialPosition := l.currentPosition
			for isDigit(l.currentChar, previousChar, l.peekChar) {
				l.readChar()
			}
			sequence := l.input[initialPosition:l.currentPosition]
			tokType, ok := token.Keywords[sequence]
			if !ok {
				return token.Token{Type: token.INTEGER, Literal: sequence}
			} else {
				return token.Token{Type: tokType, Literal: sequence}
			}
		} else {
			tok = newToken(token.UNKNOWN, l.currentChar)
		}
	}
	l.readChar()
	return tok
}

func NewLexer(input string) *Lexer {
	lexer := &Lexer{
		input:           strings.TrimSpace(input),
		currentPosition: 0,
		peekPosition:    0,
		currentChar:     input[0],
		peekChar:        0,
		line:            1,
		column:          0,
	}
	if len(input) > 1 {
		lexer.peekChar = input[1]
		lexer.peekPosition = 1
	}
	return lexer
}

func (l *Lexer) readChar() {
	if l.currentChar == '\n' {
		l.line += 1
		l.column = 0
	} else {
		l.column += 1
	}
	l.currentPosition = l.peekPosition
	if l.peekPosition > len(l.input)-1 {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.peekPosition]
	}
	l.peekPosition += 1
	if l.peekPosition > len(l.input)-1 {
		l.peekChar = 0
	} else {
		l.peekChar = l.input[l.peekPosition]
	}
}
