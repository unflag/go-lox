package main

import (
	"errors"
	"fmt"
)

type loxCallable[T any] interface {
	call(i *Interpreter[T], args []any) T
	arity() int
}

type ReturnValue struct {
	Value Expr
}

type Interpreter[T any] struct {
	globals *Environment
	env     *Environment
}

func NewInterpreter() *Interpreter[any] {
	globals := NewEnvironment(nil)
	globals.Define(&Token{Lexeme: "clock"}, &clock[any]{})
	return &Interpreter[any]{
		globals: globals,
		env:     globals,
	}
}

func (i *Interpreter[T]) Interpret(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			var runtimeErr RuntimeError
			if errors.As(r.(error), &runtimeErr) {
				ReportRuntimeError(runtimeErr)
			} else {
				panic(r)
			}
		}
	}()

	for n, s := range statements {
		i.execute(s)
		if n == len(statements)-1 {
			if _, ok := s.(*Expression); !ok {
				expressionValue = nil
			}
		}
	}
}

func (i *Interpreter[T]) evaluate(e Expr) T {
	return AcceptExprVisitor[T](e, i)
}

func (i *Interpreter[T]) execute(s Stmt) {
	AcceptStmtVisitor[T](s, i)
}

func (i *Interpreter[T]) VisitExpressionStmt(s *Expression) {
	expressionValue = i.evaluate(s.Expression)
}

func (i *Interpreter[T]) VisitFunctionStmt(s *Function) {
	fn := newLoxFunction(s, i.env)
	i.env.Define(s.Name, fn)
}

func (i *Interpreter[T]) VisitIfStmt(s *If) {
	if toBool(i.evaluate(s.Expression)) {
		i.execute(s.ThenBranch)
	} else {
		i.execute(s.ElseBranch)
	}
}

func (i *Interpreter[T]) VisitPrintStmt(s *Print) {
	value := i.evaluate(s.Expression)
	fmt.Printf("%v\n", value)
}

func (i *Interpreter[T]) VisitReturnStmt(s *Return) {
	var value Expr = &NilT{}
	if _, ok := s.Value.(*NilT); !ok {
		value = i.evaluate(s.Value)
	}
	panic(&ReturnValue{Value: value})
}

func (i *Interpreter[T]) VisitWhileStmt(s *While) {
	for toBool(i.evaluate(s.Condition)) {
		i.execute(s.Body)
	}
}

func (i *Interpreter[T]) VisitVarStmt(s *Var) {
	var value interface{}
	if s.Initializer != nil {
		value = i.evaluate(s.Initializer)
	}

	i.env.Define(s.Name, value)
}

func (i *Interpreter[T]) VisitBlockStmt(s *Block) {
	i.executeBlock(s.Statements, i.env)
}

func (i *Interpreter[T]) executeBlock(stmts []Stmt, env *Environment) {
	previousEnv := i.env
	defer func() {
		i.env = previousEnv
	}()

	i.env = env
	for _, stmt := range stmts {
		i.execute(stmt)
	}
}

func (i *Interpreter[T]) VisitLogicalExpr(e *Logical) T {
	left := any(i.evaluate(e.Left))

	if e.Operator.Type == OR {
		if toBool(left) {
			return left.(T)
		}
	} else {
		if !toBool(left) {
			return left.(T)
		}
	}

	return any(i.evaluate(e.Right)).(T)
}

func (i *Interpreter[T]) VisitBinaryExpr(e *Binary) T {
	var v any

	left := any(i.evaluate(e.Left))
	right := any(i.evaluate(e.Right))

	lStr, lok := left.(string)
	rStr, rok := right.(string)
	if lok && rok {
		return i.visitBinaryStringExpr(e.Operator, lStr, rStr)
	}

	lFl, lok := left.(float64)
	rFl, rok := right.(float64)
	if lok && rok {
		return i.visitBinaryFloat64Expr(e.Operator, lFl, rFl)
	}

	if !lok || !rok {
		panic(NewRuntimeError(e.Operator, fmt.Sprintf("Unsupported operands: %v %v", left, right)))
	}

	return v.(T)
}

func (i *Interpreter[T]) visitBinaryStringExpr(op *Token, l, r string) T {
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
		panic(NewRuntimeError(op, fmt.Sprintf("Unsupported operands: %v %v", l, r)))
	}
	return v.(T)
}

func (i *Interpreter[T]) visitBinaryFloat64Expr(op *Token, l, r float64) T {
	var v any
	switch op.Type {
	case MINUS:
		v = l - r
	case SLASH:
		if r == 0 {
			panic(NewRuntimeError(op, "Division by zero"))
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
		panic(NewRuntimeError(op, fmt.Sprintf("Unsupported operands: %v %v", l, r)))
	}
	return v.(T)
}

func (i *Interpreter[T]) VisitCallExpr(e *Call) T {
	callee := i.evaluate(e.Callee)

	args := make([]interface{}, 0, len(e.Args))
	for _, arg := range e.Args {
		args = append(args, i.evaluate(arg))
	}

	f, ok := any(callee).(loxCallable[T])
	if !ok {
		panic(NewRuntimeError(e.Paren, "Can only call functions and classes."))
	}

	if f.arity() != len(args) {
		panic(NewRuntimeError(e.Paren, fmt.Sprintf("Expected %d arguments but got %d.", f.arity(), len(args))))
	}

	return f.call(i, args)
}

func (i *Interpreter[T]) VisitGroupingExpr(e *Grouping) T {
	return i.evaluate(e.Expression)
}

func (i *Interpreter[T]) VisitLiteralExpr(e *Literal) T {
	if e.Value == nil {
		var v interface{} = "nil"
		return v.(T)
	}
	return e.Value.(T)
}

func (i *Interpreter[T]) VisitUnaryExpr(e *Unary) T {
	var v any
	right := i.evaluate(e.Right)
	switch e.Operator.Type {
	case MINUS:
		if _, ok := any(right).(float64); !ok {
			panic(NewRuntimeError(e.Operator, fmt.Sprintf("Cannot negate %T", right)))
		}
		v = -(any(right).(float64))
	case BANG:
		v = any(!toBool(right)).(T)
	default:
	}
	return v.(T)
}

func (i *Interpreter[T]) VisitVariableExpr(e *Variable) T {
	val := i.env.Get(e.Name)
	if val == nil {
		panic(NewRuntimeError(e.Name, "Uninitialized variable"))
	}

	return val.(T)
}

func (i *Interpreter[T]) VisitAssignExpr(e *Assign) T {
	value := i.evaluate(e.Value)
	i.env.Assign(e.Name, value)

	return value
}

func toBool(obj any) bool {
	if obj == nil {
		return false
	}

	if _, ok := obj.(NilT); ok {
		return false
	}

	if b, ok := obj.(bool); ok {
		return b
	}

	return true
}
