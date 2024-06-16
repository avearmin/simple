package lexer

import (
	"github.com/avearmin/simple/internal/token"
	"testing"
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
				{token.LParen, "("},
				{token.Assign, ":="},
				{token.Delimiter, " "},
				{token.Ident, "foo"},
				{token.Delimiter, " "},
				{token.Int, "1"},
				{token.RParen, ")"},
				{token.Delimiter, "\n"},
				{token.LParen, "("},
				{token.Reassign, "="},
				{token.Delimiter, " "},
				{token.Ident, "foo"},
				{token.Delimiter, " "},
				{token.LParen, "("},
				{token.Add, "+"},
				{token.Delimiter, " "},
				{token.Ident, "foo"},
				{token.Delimiter, " "},
				{token.Int, "1"},
				{token.RParen, ")"},
				{token.RParen, ")"},
				{token.EOF, ""},
			},
		},
		"Illegal token": {
			input: `(:= foo 5555xxxx)`,
			want: []token.Token{
				{token.LParen, "("},
				{token.Assign, ":="},
				{token.Delimiter, " "},
				{token.Ident, "foo"},
				{token.Delimiter, " "},
				{token.Illegal, "5555xxxx"},
				{token.RParen, ")"},
				{token.EOF, ""},
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
	return (tokenOne.Type == tokenTwo.Type) && (tokenOne.Literal == tokenTwo.Literal)
}
