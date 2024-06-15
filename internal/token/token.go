package token

type Type string

const (
	Illegal = "ILLEGAL"
	EOF     = "EOF"

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

func New(tokType Type, char byte) Token {
	return Token{Type: tokType, Literal: string(char)}
}
