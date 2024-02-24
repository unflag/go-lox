package scanner

import (
	"fmt"
	"log"
	"strconv"
)

type Scanner struct {
	source []rune
	tokens []*Token

	start, current, line int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: []rune(source),
	}
}

func (s *Scanner) Scan() ([]*Token, error) {
	for !s.isEOF() {
		s.start = s.current
		if err := s.scanToken(); err != nil {
			log.Printf("could not scan token: %+v", err)
		}
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))

	return s.tokens, nil
}

func (s *Scanner) isEOF() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
	char := s.next()
	switch true {
	case char == '(':
		s.tokens = append(s.tokens, NewToken(LEFT_PAREN, string(s.source[s.start:s.current]), nil, s.line))
	case char == ')':
		s.tokens = append(s.tokens, NewToken(RIGHT_PAREN, string(s.source[s.start:s.current]), nil, s.line))
	case char == '{':
		s.tokens = append(s.tokens, NewToken(LEFT_BRACE, string(s.source[s.start:s.current]), nil, s.line))
	case char == '}':
		s.tokens = append(s.tokens, NewToken(RIGHT_BRACE, string(s.source[s.start:s.current]), nil, s.line))
	case char == ',':
		s.tokens = append(s.tokens, NewToken(COMMA, string(s.source[s.start:s.current]), nil, s.line))
	case char == '.':
		s.tokens = append(s.tokens, NewToken(DOT, string(s.source[s.start:s.current]), nil, s.line))
	case char == '-':
		s.tokens = append(s.tokens, NewToken(MINUS, string(s.source[s.start:s.current]), nil, s.line))
	case char == '+':
		s.tokens = append(s.tokens, NewToken(PLUS, string(s.source[s.start:s.current]), nil, s.line))
	case char == ';':
		s.tokens = append(s.tokens, NewToken(SEMICOLON, string(s.source[s.start:s.current]), nil, s.line))
	case char == '*':
		s.tokens = append(s.tokens, NewToken(STAR, string(s.source[s.start:s.current]), nil, s.line))
	case char == '!':
		var typ = BANG
		if s.nextMatch('=') {
			typ = BANG_EQUAL
		}
		s.tokens = append(s.tokens, NewToken(typ, string(s.source[s.start:s.current]), nil, s.line))
	case char == '=':
		var typ = EQUAL
		if s.nextMatch('=') {
			typ = EQUAL_EQUAL
		}
		s.tokens = append(s.tokens, NewToken(typ, string(s.source[s.start:s.current]), nil, s.line))
	case char == '<':
		var typ = LESS
		if s.nextMatch('=') {
			typ = LESS_EQUAL
		}
		s.tokens = append(s.tokens, NewToken(typ, string(s.source[s.start:s.current]), nil, s.line))
	case char == '>':
		var typ = GREATER
		if s.nextMatch('=') {
			typ = GREATER_EQUAL
		}
		s.tokens = append(s.tokens, NewToken(typ, string(s.source[s.start:s.current]), nil, s.line))
	case char == '/':
		c := s.peek(0)
		if s.nextMatch('/') || s.nextMatch('*') {
			s.readComment(c)
			break
		}
		s.tokens = append(s.tokens, NewToken(SLASH, string(s.source[s.start:s.current]), nil, s.line))
	case char == '"':
		s.readString()
	case isDigit(char):
		s.readNumber()
	case isAlpha(char):
		s.readIdentifier()
	case char == ' ' || char == '\r' || char == '\t':
	case char == '\n':
		s.line++
	default:
		return fmt.Errorf("unexpected character at position %d:%d: %s", s.line, s.current, string(char))
	}

	return nil
}

func (s *Scanner) next() rune {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) nextMatch(char rune) bool {
	if s.isEOF() {
		return false
	}

	if s.source[s.current] != char {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek(offset int) rune {
	if s.isEOF() || s.current+offset >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+offset]
}

func (s *Scanner) readString() {
	for s.peek(0) != '"' && !s.isEOF() {
		if s.peek(0) == '\n' {
			s.line++
		}
		s.next()
	}

	if s.isEOF() {
		log.Printf("Unterminated string at %d:%d.", s.line, s.start)
		return
	}

	// the closing "
	s.next()

	s.tokens = append(s.tokens, NewToken(STRING, string(s.source[s.start+1:s.current-1]), nil, s.line))
}

func (s *Scanner) readNumber() {
	for isDigit(s.peek(0)) {
		s.next()
	}

	if s.peek(0) == '.' && isDigit(s.peek(1)) {
		s.next()
		for isDigit(s.peek(0)) {
			s.next()
		}
	}

	number, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	if err != nil {
		log.Printf("could not parse number %s: %+v", string(s.source[s.start:s.current]), err)
		return
	}

	s.tokens = append(s.tokens, NewToken(NUMBER, "", number, s.line))
}

func (s *Scanner) readIdentifier() {
	for isAlphaNumeric(s.peek(0)) {
		s.next()
	}

	word := string(s.source[s.start:s.current])
	var typ TokenType = IDENTIFIER
	if t, ok := reservedWords[word]; ok {
		typ = t
	}

	s.tokens = append(s.tokens, NewToken(typ, word, nil, s.line))
}

func (s *Scanner) readComment(char rune) {
	switch char {
	case '/':
		for s.peek(0) != '\n' && !s.isEOF() {
			s.next()
		}
	case '*':
		for !s.isEOF() {
			if c := s.next(); c == '*' && s.nextMatch('/') {
				break
			}
		}
	}
}

func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func isAlpha(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}

func isAlphaNumeric(char rune) bool {
	return isAlpha(char) || isDigit(char)
}
