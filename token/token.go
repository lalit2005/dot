package token

type TokenType string

const (
	IDENTIFIER = "IDENTIFIER"
	INTEGER    = "INTEGER"
	TRUE       = "TRUE"
	FALSE      = "FALSE"
	LET        = "LET"
	IF         = "IF"
	EOF        = "EOF"
	FUNCTION   = "FUNCTION"
	STRING     = "STRING"

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

	UNKNOWN = "UNKNOWN"
)

type Token struct {
	Type    TokenType
	Literal string
}

var Keywords = map[string]TokenType{
	"fn":    FUNCTION,
	"true":  TRUE,
	"false": FALSE,
	"let":   LET,
	"if":    IF,
}
