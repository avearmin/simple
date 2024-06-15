package lexer

import "github.com/avearmin/simple/internal/token"

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

	l.consumeWhitespaces()

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
			op := string(l.input[pos:l.pos])
			tok = token.NewFromString(token.Assign, op)
			return tok
		}
	case 0:
		tok = token.NewFromString(token.EOF, "")
	default:
		if isLetter(l.char) {
			ident := l.readIdent()
			return token.NewFromString(token.Ident, ident)
		} else if isDigit(l.char) {
			num := l.readNumber()
			return token.NewFromString(token.Int, num)
		} else {
			return token.NewFromByte(token.Illegal, l.char)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) consumeWhitespaces() {
	for l.char == ' ' || l.char == '\n' || l.char == '\t' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdent() string {
	pos := l.pos
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
