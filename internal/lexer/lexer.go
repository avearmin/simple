package lexer

import (
	"fmt"
	"github.com/avearmin/simple/internal/token"
)

type Lexer struct {
	input   string
	pos     int
	nextPos int
	char    byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextPos >= len(l.input) {
		l.char = 0
	} else {
		l.pos = l.nextPos
		l.nextPos = l.pos + 1
		l.char = l.input[l.pos]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	fmt.Println("|" + string(l.char) + "|")
	switch l.char {
	case '(':
		tok = token.NewFromByte(token.LParen, l.char)
	case ')':
		tok = token.NewFromByte(token.RParen, l.char)
	case '+':
		tok = token.NewFromByte(token.Add, l.char)
	case '=':
		tok = token.NewFromByte(token.Reassign, l.char)
	case ':':
		if l.input[l.nextPos] == '=' {
			pos := l.pos
			l.readChar()
			l.readChar()
			op := l.input[pos:l.pos]
			tok = token.NewFromString(token.Assign, op)
			return tok
		}
	case 0:
		tok = token.NewFromString(token.EOF, "")
	default:
		if isWhitespace(l.char) {
			delimiter := l.readWhitespaces()
			tok = token.NewFromString(token.Delimiter, delimiter)
			return tok
		}

		ident := l.readIdent()
		if isIdentInt(ident) {
			tok = token.NewFromString(token.Int, ident)
		} else if isIdentValid(ident) {
			tok = token.NewFromString(token.Ident, ident)
		} else {
			tok = token.NewFromString(token.Illegal, ident)
		}
		return tok
	}

	l.readChar()
	return tok
}

func (l *Lexer) readWhitespaces() string {
	pos := l.pos
	for isWhitespace(l.char) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readIdent() string {
	pos := l.pos
	for !isWhitespace(l.char) && l.char != ')' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isIdentInt(ident string) bool {
	for _, c := range []byte(ident) {
		if !isDigit(c) {
			return false
		}
	}
	return true
}

func isIdentValid(ident string) bool {
	for _, c := range []byte(ident) {
		if !isLetter(c) {
			return false
		}
	}
	return true
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
