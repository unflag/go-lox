package main

import "time"

type clock[T any] struct{}

func (c *clock[T]) arity() int {
	return 0
}

func (c *clock[T]) call(i *Interpreter[T], args []any) T {
	return any(time.Now().Second()).(T)
}

func (c *clock[T]) String() string {
	return "<native fn clock>"
}
