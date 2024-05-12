package main

import (
	"fmt"
	"maps"
)

type loxFunction[T any] struct {
	declaration *Function
	closure     *Environment
}

func newLoxFunction(declaration *Function, env *Environment) *loxFunction[any] {
	closure := NewEnvironment(nil)
	maps.Copy(closure.values, env.values)
	return &loxFunction[any]{
		declaration: declaration,
		closure:     env,
	}
}

func (f *loxFunction[T]) call(i *Interpreter[T], args []any) (retVal T) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(*ReturnValue); ok {
				retVal = v.Value.(T)
				return
			} else {
				panic(r)
			}
		}
		retVal = any(&NilT{}).(T)
	}()

	env := NewEnvironment(f.closure)
	for i, p := range f.declaration.Params {
		env.Define(p, args[i])
	}

	i.executeBlock(f.declaration.Body, env)

	return retVal
}

func (f *loxFunction[T]) arity() int {
	return len(f.declaration.Params)
}

func (f *loxFunction[T]) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}
