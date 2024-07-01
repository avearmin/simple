package lexer

import (
	"testing"

	"github.com/avearmin/simple/internal/token"
)

func TestIsIdentInt(t *testing.T) {
	tests := map[string]struct {
		input string
		want  bool
	}{
		"not int": {
			input: "5555xxxx",
			want:  false,
		},
		"is int": {
			input: "5555",
			want:  true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := isIdentInt(test.input); got != test.want {
				t.Errorf("isIdentInt(%q) = %v, want %v", test.input, got, test.want)
			}
		})
	}
}

func TestIsIdentValid(t *testing.T) {
	tests := map[string]struct {
		input string
		want  bool
	}{
		"not valid": {
			input: "5555xxxx",
			want:  false,
		},
		"valid": {
			input: "foo",
			want:  true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := isIdentValid(test.input); got != test.want {
				t.Errorf("isIdentInt(%q) = %v, want %v", test.input, got, test.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []token.Token
	}{
		"Assign/Reassign Int Ident with arithmatic operators": {
			input: `(:= foo 1)
(= foo (+ foo 1))
(= foo (- foo 1))
(= foo (* foo 2))
(= foo (/ foo 2))
(= foo (% foo 2))`,
			want: []token.Token{
				{token.LParen, "(", 1, 0},
				{token.Assign, ":=", 1, 1},
				{token.Delimiter, " ", 1, 3},
				{token.Ident, "foo", 1, 4},
				{token.Delimiter, " ", 1, 7},
				{token.Int, "1", 1, 8},
				{token.RParen, ")", 1, 9},
				{token.Delimiter, "\n", 1, 10},
				{token.LParen, "(", 2, 0},
				{token.Reassign, "=", 2, 1},
				{token.Delimiter, " ", 2, 2},
				{token.Ident, "foo", 2, 3},
				{token.Delimiter, " ", 2, 6},
				{token.LParen, "(", 2, 7},
				{token.Add, "+", 2, 8},
				{token.Delimiter, " ", 2, 9},
				{token.Ident, "foo", 2, 10},
				{token.Delimiter, " ", 2, 13},
				{token.Int, "1", 2, 14},
				{token.RParen, ")", 2, 15},
				{token.RParen, ")", 2, 16},
				{token.Delimiter, "\n", 2, 17},
				{token.LParen, "(", 3, 0},
				{token.Reassign, "=", 3, 1},
				{token.Delimiter, " ", 3, 2},
				{token.Ident, "foo", 3, 3},
				{token.Delimiter, " ", 3, 6},
				{token.LParen, "(", 3, 7},
				{token.Subtract, "-", 3, 8},
				{token.Delimiter, " ", 3, 9},
				{token.Ident, "foo", 3, 10},
				{token.Delimiter, " ", 3, 13},
				{token.Int, "1", 3, 14},
				{token.RParen, ")", 3, 15},
				{token.RParen, ")", 3, 16},
				{token.Delimiter, "\n", 3, 17},
				{token.LParen, "(", 4, 0},
				{token.Reassign, "=", 4, 1},
				{token.Delimiter, " ", 4, 2},
				{token.Ident, "foo", 4, 3},
				{token.Delimiter, " ", 4, 6},
				{token.LParen, "(", 4, 7},
				{token.Multiply, "*", 4, 8},
				{token.Delimiter, " ", 4, 9},
				{token.Ident, "foo", 4, 10},
				{token.Delimiter, " ", 4, 13},
				{token.Int, "2", 4, 14},
				{token.RParen, ")", 4, 15},
				{token.RParen, ")", 4, 16},
				{token.Delimiter, "\n", 4, 17},
				{token.LParen, "(", 5, 0},
				{token.Reassign, "=", 5, 1},
				{token.Delimiter, " ", 5, 2},
				{token.Ident, "foo", 5, 3},
				{token.Delimiter, " ", 5, 6},
				{token.LParen, "(", 5, 7},
				{token.Divide, "/", 5, 8},
				{token.Delimiter, " ", 5, 9},
				{token.Ident, "foo", 5, 10},
				{token.Delimiter, " ", 5, 13},
				{token.Int, "2", 5, 14},
				{token.RParen, ")", 5, 15},
				{token.RParen, ")", 5, 16},
				{token.Delimiter, "\n", 5, 17},
				{token.LParen, "(", 6, 0},
				{token.Reassign, "=", 6, 1},
				{token.Delimiter, " ", 6, 2},
				{token.Ident, "foo", 6, 3},
				{token.Delimiter, " ", 6, 6},
				{token.LParen, "(", 6, 7},
				{token.Modulo, "%", 6, 8},
				{token.Delimiter, " ", 6, 9},
				{token.Ident, "foo", 6, 10},
				{token.Delimiter, " ", 6, 13},
				{token.Int, "2", 6, 14},
				{token.RParen, ")", 6, 15},
				{token.RParen, ")", 6, 16},
				{token.EOF, "", 6, 16},
			},
		},
		"Illegal token": {
			input: `(:= foo 5555xxxx)`,
			want: []token.Token{
				{token.LParen, "(", 1, 0},
				{token.Assign, ":=", 1, 1},
				{token.Delimiter, " ", 1, 3},
				{token.Ident, "foo", 1, 4},
				{token.Delimiter, " ", 1, 7},
				{token.Illegal, "5555xxxx", 1, 8},
				{token.RParen, ")", 1, 16},
				{token.EOF, "", 1, 16},
			},
		},
		"boolean assign/reassign with logical operators": {
			input: `(:= foo false)
(= foo true)
(! foo)
(== foo true)
(!= foo true)`,
			want: []token.Token{
				{token.LParen, "(", 1, 0},
				{token.Assign, ":=", 1, 1},
				{token.Delimiter, " ", 1, 3},
				{token.Ident, "foo", 1, 4},
				{token.Delimiter, " ", 1, 7},
				{token.Bool, "false", 1, 8},
				{token.RParen, ")", 1, 13},
				{token.Delimiter, "\n", 1, 14},
				{token.LParen, "(", 2, 0},
				{token.Reassign, "=", 2, 1},
				{token.Delimiter, " ", 2, 2},
				{token.Ident, "foo", 2, 3},
				{token.Delimiter, " ", 2, 6},
				{token.Bool, "true", 2, 7},
				{token.RParen, ")", 2, 11},
				{token.Delimiter, "\n", 2, 12},
				{token.LParen, "(", 3, 0},
				{token.Not, "!", 3, 1},
				{token.Delimiter, " ", 3, 2},
				{token.Ident, "foo", 3, 3},
				{token.RParen, ")", 3, 6},
				{token.Delimiter, "\n", 3, 7},
				{token.LParen, "(", 4, 0},
				{token.Equals, "==", 4, 1},
				{token.Delimiter, " ", 4, 3},
				{token.Ident, "foo", 4, 4},
				{token.Delimiter, " ", 4, 7},
				{token.Bool, "true", 4, 8},
				{token.RParen, ")", 4, 12},
				{token.Delimiter, "\n", 4, 13},
				{token.LParen, "(", 5, 0},
				{token.NotEquals, "!=", 5, 1},
				{token.Delimiter, " ", 5, 3},
				{token.Ident, "foo", 5, 4},
				{token.Delimiter, " ", 5, 7},
				{token.Bool, "true", 5, 8},
				{token.RParen, ")", 5, 12},
				{token.EOF, "", 5, 12},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lexer := New(test.input)
			for i := range test.want {
				got := lexer.NextToken()
				if !isEqualTokens(got, test.want[i]) {
					t.Fatalf("%d: got=%+v, but want=%+v", i, got, test.want[i])
				}
			}
		})

	}
}

func isEqualTokens(tokenOne, tokenTwo token.Token) bool {
	return (tokenOne.Type == tokenTwo.Type) && (tokenOne.Literal == tokenTwo.Literal) && (tokenOne.Line == tokenTwo.Line) && (tokenOne.Col == tokenTwo.Col)
}
