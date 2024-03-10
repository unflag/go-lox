package main

import (
	"fmt"
)

type ParseError struct {
	Token   *Token
	Message string
}

var _ error = ParseError{}

func NewParseError(token *Token, message string) error {
	return ParseError{
		Token:   token,
		Message: message,
	}
}

func (e ParseError) Error() string {
	return fmt.Sprintf(
		"[%d]: parse error at %s %s: %s",
		e.Token.Line,
		e.Token.Type,
		e.Token.Lexeme,
		e.Message,
	)
}

type RuntimeError struct {
	Token   *Token
	Message string
}

var _ error = ParseError{}

func NewRuntimeError(token *Token, message string) error {
	return RuntimeError{
		Token:   token,
		Message: message,
	}
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf(
		"[%d]: runtime error at %s %s: %s",
		e.Token.Line,
		e.Token.Type,
		e.Token.Lexeme,
		e.Message,
	)
}
