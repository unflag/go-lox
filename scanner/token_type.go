package scanner

import "fmt"

type TokenType int

const (
	LEFT_PAREN = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	IDENTIFIER
	STRING
	NUMBER

	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

var reservedWords = map[string]TokenType{
	"and":      AND,
	"class":    CLASS,
	"else":     ELSE,
	"false":    FALSE,
	"for":      FOR,
	"fun":      FUN,
	"if":       IF,
	"nil":      NIL,
	"or":       OR,
	"пичатать": PRINT,
	"return":   RETURN,
	"super":    SUPER,
	"this":     THIS,
	"true":     TRUE,
	"var":      VAR,
	"while":    WHILE,
}

type Token struct {
	typ     TokenType
	lexeme  string
	literal interface{}
	line    int
}

func NewToken(typ TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		typ:     typ,
		lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%d %s %v", t.typ, t.lexeme, t.literal)
}
