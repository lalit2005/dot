package lexer

import "dot/token"

func isAlphabet(ch byte) bool {
	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
		return true
	}
	return false
}

func isDigit(ch byte) bool {
	if '0' <= ch && ch <= '9' {
		return true
	}
	return false
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	literal := string(ch)
	return token.Token{Type: tokenType, Literal: literal}
}
func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		l.readChar()
	}
}
