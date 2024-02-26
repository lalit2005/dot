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
	WHILE      = "WHILE"
	FOR        = "FOR"

	PLUS        = "+"
	MINUS       = "-"
	SLASH       = "/"
	ASTERISK    = "*"
	EQUAL       = "=="
	PLUS_EQUAL  = "+="
	MINUS_EQUAL = "-="
	MULT_EQUAL  = "*="
	DIV_EQUAL   = "/="
	NOT_EQUAL   = "!="
	ASSIGN      = "="
	LPAREN      = "("
	RPAREN      = ")"
	LBRACE      = "{"
	RBRACE      = "}"
	SEMICOLON   = ";"
	COMMA       = ","
	LT          = "<"
	GT          = ">"
	LTE         = "<="
	GTE         = ">="
	BANG        = "!"
	COLON       = ":"
	LBRACKET    = "["
	RBRACKET    = "]"
	COMMENT     = "//"
	AND         = "&&"
	OR          = "||"

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
	"while":  WHILE,
	"for":    FOR,
}
