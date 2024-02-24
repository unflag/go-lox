package parser

import (
	"fmt"

	"github.com/unflag/go-lox/scanner"
)

type ParseError struct {
	Token   *scanner.Token
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf(
		"[%d]: error at %s %s: %s",
		e.Token.Line(),
		e.Token.Type(),
		e.Token.Lexeme(),
		e.Message,
	)
}
