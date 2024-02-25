package main

import (
	"fmt"
)

type ParseError struct {
	Token   *Token
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf(
		"[%d]: error at %s %s: %s",
		e.Token.Line,
		e.Token.Type,
		e.Token.Lexeme,
		e.Message,
	)
}
