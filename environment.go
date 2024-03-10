package main

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *Environment) Define(key *Token, value interface{}) {
	e.values[key.Lexeme] = value
}

func (e *Environment) Assign(key *Token, value interface{}) {
	_, ok := e.values[key.Lexeme]
	if !ok && e.enclosing == nil {
		panic(NewRuntimeError(key, "Undefined variable"))
	}

	if !ok {
		e.enclosing.Assign(key, value)
		return
	}

	e.values[key.Lexeme] = value
}

func (e *Environment) Get(key *Token) interface{} {
	val, ok := e.values[key.Lexeme]
	if !ok && e.enclosing == nil {
		panic(NewRuntimeError(key, "Undefined variable"))
	}

	if !ok {
		return e.enclosing.Get(key)
	}

	return val
}
