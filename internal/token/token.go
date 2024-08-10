package token

type Type string

const (
	Illegal = "ILLEGAL"
	EOF     = "EOF"

	Delimiter = "DELIMITER"

	LParen = "("
	RParen = ")"

	Int  = "INT"
	Bool = "BOOL"

	Assign   = ":="
	Reassign = "="

	Add      = "+"
	Subtract = "-"
	Divide   = "/"
	Multiply = "*"
	Modulo   = "%"

	Not = "!"
	And = "&&"
	Or  = "||"

	Equals              = "=="
	NotEquals           = "!="
	LessThan            = "<"
	GreaterThan         = ">"
	LessThanOrEquals    = "<="
	GreaterThanOrEquals = ">="

	If   = "IF"
	Elif = "ELIF"
	Else = "ELSE"

	Ident = "IDENT"
)

var identToType = map[string]Type{
	"if":   If,
	"elif": Elif,
	"else": Else,
}

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

func LookupIdent(ident string) (Type, bool) {
	if tokenType, ok := identToType[ident]; ok {
		return tokenType, true
	} else {
		return "", false
	}
}

func IsBoolToken(tokenType Type) bool {
	return tokenType == Equals ||
		tokenType == NotEquals ||
		tokenType == LessThan ||
		tokenType == GreaterThan ||
		tokenType == LessThanOrEquals ||
		tokenType == GreaterThanOrEquals ||
		tokenType == Not ||
		tokenType == And ||
		tokenType == Or
}
