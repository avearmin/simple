package lexer

import (
	"github.com/avearmin/simple/internal/token"
	"testing"
)

func TestParse(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []token.Token
	}{
		"Assign/Reassign Int Ident": {
			input: `
				(:= foo 1)
				(= foo (+ foo 1))
			`,
			want: []token.Token{
				{token.LParen, "("},
				{token.Assign, ":="},
				{token.Ident, "foo"},
				{token.Int, "1"},
				{token.RParen, ")"},
				{token.LParen, "("},
				{token.Reassign, "="},
				{token.Ident, "foo"},
				{token.LParen, "("},
				{token.Add, "+"},
				{token.Ident, "foo"},
				{token.Int, "1"},
				{token.RParen, ")"},
				{token.RParen, ")"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lexer := New(test.input)
			for i := range test.want {
				got := lexer.NextToken()
				if !isEqualTokens(got, test.want[i]) {
					t.Fatalf("got=%+v, but want=%+v", got, test.want[i])
				}
			}
		})

	}
}

func isEqualTokens(tokenOne, tokenTwo token.Token) bool {
	return (tokenOne.Type == tokenTwo.Type) && (tokenOne.Literal == tokenTwo.Literal)
}
