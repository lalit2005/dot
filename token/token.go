package token

type TokenType string

const (
	IDENTIFIER = "IDENTIFIER"
	INTEGER    = "INTEGER"
	TRUE       = "TRUE"
	FALSE      = "FALSE"
	LET        = "LET"
	IF         = "IF"
	ELSE       = "ELSE"
	EOF        = "EOF"
	FUNCTION   = "FUNCTION"
	STRING     = "STRING"
	RETURN     = "RETURN"

	PLUS      = "+"
	MINUS     = "-"
	SLASH     = "/"
	ASTERISK  = "*"
	EQUAL     = "=="
	NOT_EQUAL = "!="
	ASSIGN    = "="
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	SEMICOLON = ";"
	COMMA     = ","
	LT        = "<"
	GT        = ">"
	BANG      = "!"
	COLON     = ":"
	LBRACKET  = "["
	RBRACKET  = "]"
	COMMENT   = "//"
	AND       = "&&"
	OR        = "||"

	UNKNOWN = "UNKNOWN"
)

type Token struct {
	Type    TokenType
	Literal string
}

var Keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
	"let":    LET,
	"if":     IF,
	"return": RETURN,
	"else":   ELSE,
}
