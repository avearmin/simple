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
		"Assign/Reassign Int Ident": {
			input: `(:= foo 1)
(= foo (+ foo 1))`,
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
				{token.EOF, "", 2, 16},
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
