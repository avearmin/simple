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

	Add = "+"

	Ident = "IDENT"
)

type Token struct {
	Type    Type
	Literal string
}

func NewFromByte(tokType Type, char byte) Token {
	return Token{Type: tokType, Literal: string(char)}
}

func NewFromString(tokType Type, str string) Token {
	return Token{Type: tokType, Literal: str}
}
