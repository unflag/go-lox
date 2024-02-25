package main

import (
	"fmt"
)

type Interpreter[T any] struct{}

func newInterpreter() *Interpreter[any] {
	return &Interpreter[any]{}
}

func (i *Interpreter[T]) Evaluate(e Expr) T {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Error evaluating expression: %s\n", err)
		}
	}()

	return Accept[T](e, i)
}

func (i *Interpreter[T]) VisitBinary(e *Binary) T {
	var v any

	left := any(i.Evaluate(e.Left))
	right := any(i.Evaluate(e.Right))

	lStr, lok := left.(string)
	rStr, rok := right.(string)
	if lok && rok {
		return i.visitBinaryString(e.Operator, lStr, rStr)
	}

	lFl, lok := left.(float64)
	rFl, rok := right.(float64)
	if lok && rok {
		return i.visitBinaryFloat64(e.Operator, lFl, rFl)
	}

	if !lok || !rok {
		panic(fmt.Sprintf("[Line %d] unsupported operands: %v %v %v", e.Operator.Line, left, e.Operator.Type, right))
	}

	return v.(T)
}

func (i *Interpreter[T]) visitBinaryString(op *Token, l, r string) T {
	var v any
	switch op.Type {
	case PLUS:
		v = fmt.Sprintf("%s%s", l, r)
	case GREATER:
		v = l > r
	case GREATER_EQUAL:
		v = l >= r
	case LESS:
		v = l < r
	case LESS_EQUAL:
		v = l <= r
	case BANG_EQUAL:
		v = l != r
	case EQUAL_EQUAL:
		v = l == r
	default:
		panic(fmt.Sprintf("[Line %d] unsupported operands: %v %v %v", op.Line, l, op.Type, r))
	}
	return v.(T)
}

func (i *Interpreter[T]) visitBinaryFloat64(op *Token, l, r float64) T {
	var v any
	switch op.Type {
	case MINUS:
		v = l - r
	case SLASH:
		if r == 0 {
			panic("division by zero")
		}
		v = l / r
	case STAR:
		v = l * r
	case PLUS:
		v = l + r
	case GREATER:
		v = l > r
	case GREATER_EQUAL:
		v = l >= r
	case LESS:
		v = l < r
	case LESS_EQUAL:
		v = l <= r
	case BANG_EQUAL:
		v = l != r
	case EQUAL_EQUAL:
		v = l == r
	default:
		panic(fmt.Sprintf("[Line %d] unsupported operands: %v %v %v", op.Line, l, op.Type, r))
	}
	return v.(T)
}

func (i *Interpreter[T]) VisitGrouping(e *Grouping) T {
	return i.Evaluate(e.Expression)
}

func (i *Interpreter[T]) VisitLiteral(e *Literal) T {
	return e.Value.(T)
}

func (i *Interpreter[T]) VisitUnary(e *Unary) T {
	var v any
	right := i.Evaluate(e.Right)
	switch e.Operator.Type {
	case MINUS:
		if _, ok := any(right).(float64); !ok {
			panic(fmt.Sprintf("[Line %d] cannot negate %T", e.Operator.Line, right))
		}
		v = -(any(right).(float64))
	case BANG:
		v = any(!toBool(right)).(T)
	default:
	}
	return v.(T)
}

func toBool(obj any) bool {
	if obj == nil {
		return false
	}

	if b, ok := obj.(bool); ok {
		return b
	}

	return true
}
