package parser

import (
	"fmt"
	"strings"
)

type Printer[T string] struct{}

func NewPrinter() *Printer[string] {
	return &Printer[string]{}
}

func (p *Printer[T]) Print(e Expr) T {
	return accept[T](e, p)
}

func (p *Printer[T]) visit(e Expr) T {
	switch e := e.(type) {
	case *Binary:
		return p.parenthesize(e.Operator.Lexeme(), e.Left, e.Right)
	case *Grouping:
		return p.parenthesize("group", e.Expression)
	case *Literal:
		return T(fmt.Sprintf("%v", e.Value))
	case *Unary:
		return p.parenthesize(e.Operator.Lexeme(), e.Right)
	default:
		return ""
	}
}

func (p *Printer[T]) parenthesize(name string, exprs ...Expr) T {
	expression := make([]string, 0, len(exprs))
	for _, e := range exprs {
		expression = append(expression, fmt.Sprintf("%v", accept[T](e, p)))
	}
	return T(fmt.Sprintf("%s%s %s%s", "(", name, strings.Join(expression, " "), ")"))
}
