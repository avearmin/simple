package token

type Type string

const (
	Illegal = "ILLEGAL"
	EOF     = "EOF"

	Delimiter = "DELIMITER"

	LParen = "("
	RParen = ")"

	Int = "INT"

	Assign   = ":="
	Reassign = "="

	Add      = "+"
	Subtract = "-"
	Divide   = "/"
	Multiply = "*"
	Modulo   = "%"

	Ident = "IDENT"
)

type Token struct {
	Type    Type
	Literal string
	Line    int
	Col     int
}

func NewFromByte(tokType Type, char byte, line, col int) Token {
	return Token{Type: tokType, Literal: string(char), Line: line, Col: col}
}

func NewFromString(tokType Type, str string, line, col int) Token {
	return Token{Type: tokType, Literal: str, Line: line, Col: col}
}
