package parser

import (
	"log"

	"github.com/unflag/go-lox/scanner"
)

type Parser struct {
	tokens  []*scanner.Token
	current int
}

func New(tokens []*scanner.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() Expr {
	expr, err := p.expression()
	if err != nil {
		log.Print(err)
		return nil
	}

	return expr
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
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

func (p *Parser) match(tokenTypes ...scanner.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	if p.isEOF() {
		return false
	}

	return p.peek().Type() == tokenType
}

func (p *Parser) advance() *scanner.Token {
	if !p.isEOF() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isEOF() bool {
	return p.peek().Type() == scanner.EOF
}

func (p *Parser) peek() *scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *scanner.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
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

	for p.match(scanner.MINUS, scanner.PLUS) {
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

	for p.match(scanner.SLASH, scanner.STAR) {
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
	if p.match(scanner.BANG, scanner.MINUS) {
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
	if p.match(scanner.FALSE) {
		return &Literal{Value: false}, nil
	}

	if p.match(scanner.TRUE) {
		return &Literal{Value: true}, nil
	}

	if p.match(scanner.NIL) {
		return &Literal{Value: nil}, nil
	}

	if p.match(scanner.NUMBER, scanner.STRING) {
		return &Literal{Value: p.previous().Literal()}, nil
	}

	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(scanner.RIGHT_PAREN, "expect ')' after expression."); err != nil {
			return nil, err
		}

		return &Grouping{Expression: expr}, nil
	}

	return nil, ParseError{Token: p.peek(), Message: "expect expression."}
}

func (p *Parser) consume(t scanner.TokenType, message string) (*scanner.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return nil, ParseError{Token: p.peek(), Message: message}
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isEOF() {
		if p.previous().Type() == scanner.SEMICOLON {
			return
		}

		switch p.peek().Type() {
		case scanner.CLASS, scanner.FUN, scanner.VAR, scanner.FOR, scanner.IF, scanner.WHILE, scanner.PRINT, scanner.RETURN:
			return
		default:
			p.advance()
		}
	}
}
