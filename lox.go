package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	hadError        bool
	hadRuntimeError bool
	expressionValue interface{}
)

type Lox struct {
	interpreter *Interpreter[any]
}

func NewLox() *Lox {
	return &Lox{
		interpreter: NewInterpreter(),
	}
}

func (l *Lox) RunFile(path string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	l.Run(string(src))
	if hadError {
		os.Exit(65)
	}

	if hadRuntimeError {
		os.Exit(70)
	}

	return nil
}

func (l *Lox) RunPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("could not read input: %+v", err)
			}
			break
		}

		l.Run(scanner.Text())
		hadError = false

		if expressionValue != nil {
			fmt.Printf("%v\n", expressionValue)
		}
	}

	return nil
}

func (l *Lox) Run(src string) {
	if src != "" {
		s := newScanner(src)
		tokens := s.Scan()
		p := newParser(tokens)
		stmts := p.Parse()
		l.interpreter.Interpret(stmts)
	}
}

func ReportError(err ParseError) {
	hadError = true
	fmt.Printf("%s\n", err)
}

func ReportRuntimeError(err RuntimeError) {
	hadRuntimeError = true
	expressionValue = nil
	fmt.Printf("%s\n", err)
}
