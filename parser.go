package main

import (
	"errors"
	"log"
)

type Parser struct {
	tokens  []*Token
	current int
}

func newParser(tokens []*Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() []Stmt {
	stmts := make([]Stmt, 0, 0)

	for !p.isEOF() {
		stmt, err := p.declaration()
		if err != nil {
			log.Print(err)
			return nil
		}

		stmts = append(stmts, stmt)
	}

	return stmts
}

func (p *Parser) declaration() (Stmt, error) {
	var stmt Stmt
	var err error

	switch true {
	case p.match(VAR):
		stmt, err = p.varDeclaration()
	default:
		stmt, err = p.statement()
	}

	var parseErr ParseError
	if errors.As(err, &parseErr) {
		ReportError(parseErr)
		p.synchronize()
		return nil, nil
	}

	return stmt, err
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err = p.consume(SEMICOLON, "Expect ';' after variable declaration."); err != nil {
		return nil, err
	}

	return &Var{
		Initializer: initializer,
		Name:        name,
	}, nil
}

func (p *Parser) statement() (Stmt, error) {
	switch true {
	case p.match(PRINT):
		return p.printStatement()
	case p.match(LEFT_BRACE):
		return p.blockStatement()
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) printStatement() (Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err = p.consume(SEMICOLON, "Expect ';' after value."); err != nil {
		return nil, err
	}

	return &Print{Expression: val}, nil
}

func (p *Parser) blockStatement() (Stmt, error) {
	stmts := make([]Stmt, 0)
	for !p.check(RIGHT_BRACE) && !p.isEOF() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	if _, err := p.consume(RIGHT_BRACE, "Expect '}' after block."); err != nil {
		return nil, err
	}

	return &Block{Statements: stmts}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err = p.consume(SEMICOLON, "Expect ';' after expression."); err != nil {
		return nil, err
	}

	return &Expression{Expression: expr}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		varExpr, ok := expr.(*Variable)
		if !ok {
			return nil, NewParseError(equals, "Invalid assignment target.")
		}

		return &Assign{Name: varExpr.Name, Value: value}, nil
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isEOF() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) advance() *Token {
	if !p.isEOF() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isEOF() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &Literal{Value: false}, nil
	}

	if p.match(TRUE) {
		return &Literal{Value: true}, nil
	}

	if p.match(NIL) {
		return &Literal{Value: nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &Literal{Value: p.previous().Literal}, nil
	}

	if p.match(IDENTIFIER) {
		return &Variable{Name: p.previous()}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err = p.consume(RIGHT_PAREN, "expect ')' after expression."); err != nil {
			return nil, err
		}

		return &Grouping{Expression: expr}, nil
	}

	return nil, NewParseError(p.peek(), "expect expression.")
}

func (p *Parser) consume(t TokenType, message string) (*Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return nil, NewParseError(p.peek(), message)
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isEOF() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		default:
			p.advance()
		}
	}
}
