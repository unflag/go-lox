package main

import (
	"fmt"
	"strings"
)

type Printer[T string] struct{}

func newPrinter() *Printer[string] {
	return &Printer[string]{}
}

func (p *Printer[T]) Print(e Expr) T {
	return Accept[T](e, p)
}

func (p *Printer[T]) VisitBinary(e *Binary) T {
	return p.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (p *Printer[T]) VisitGrouping(e *Grouping) T {
	return p.parenthesize("group", e.Expression)
}

func (p *Printer[T]) VisitLiteral(e *Literal) T {
	return T(fmt.Sprintf("%v", e.Value))
}

func (p *Printer[T]) VisitUnary(e *Unary) T {
	return p.parenthesize(e.Operator.Lexeme, e.Right)
}

func (p *Printer[T]) parenthesize(name string, exprs ...Expr) T {
	expression := make([]string, 0, len(exprs))
	for _, e := range exprs {
		expression = append(expression, fmt.Sprintf("%v", Accept[T](e, p)))
	}
	return T(fmt.Sprintf("%s%s %s%s", "(", name, strings.Join(expression, " "), ")"))
}
