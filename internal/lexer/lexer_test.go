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
		"Assign/Reassign Int Ident with arithmatic/comparison operators": {
			input: `(:= foo 1)
(= foo (+ foo 1))
(= foo (- foo 1))
(= foo (* foo 2))
(= foo (/ foo 2))
(= foo (% foo 2))
(:= isBar (== foo 2))
(= isBar (!= foo 3))
(= isBar (<= foo 4))
(= isBar (>= foo 5))
(= isBar (< foo 3))
(= isBar (> foo 3))`,
			want: []token.Token{
				{token.LParen, "(", 1, 0},
				{token.Assign, ":=", 1, 1},
				{token.Delimiter, "", 1, 3},
				{token.Ident, "foo", 1, 4},
				{token.Delimiter, "", 1, 7},
				{token.Int, "1", 1, 8},
				{token.RParen, ")", 1, 9},
				{token.Delimiter, "", 1, 10},
				{token.LParen, "(", 2, 0},
				{token.Reassign, "=", 2, 1},
				{token.Delimiter, "", 2, 2},
				{token.Ident, "foo", 2, 3},
				{token.Delimiter, "", 2, 6},
				{token.LParen, "(", 2, 7},
				{token.Add, "+", 2, 8},
				{token.Delimiter, "", 2, 9},
				{token.Ident, "foo", 2, 10},
				{token.Delimiter, "", 2, 13},
				{token.Int, "1", 2, 14},
				{token.RParen, ")", 2, 15},
				{token.RParen, ")", 2, 16},
				{token.Delimiter, "", 2, 17},
				{token.LParen, "(", 3, 0},
				{token.Reassign, "=", 3, 1},
				{token.Delimiter, "", 3, 2},
				{token.Ident, "foo", 3, 3},
				{token.Delimiter, "", 3, 6},
				{token.LParen, "(", 3, 7},
				{token.Subtract, "-", 3, 8},
				{token.Delimiter, "", 3, 9},
				{token.Ident, "foo", 3, 10},
				{token.Delimiter, "", 3, 13},
				{token.Int, "1", 3, 14},
				{token.RParen, ")", 3, 15},
				{token.RParen, ")", 3, 16},
				{token.Delimiter, "", 3, 17},
				{token.LParen, "(", 4, 0},
				{token.Reassign, "=", 4, 1},
				{token.Delimiter, "", 4, 2},
				{token.Ident, "foo", 4, 3},
				{token.Delimiter, "", 4, 6},
				{token.LParen, "(", 4, 7},
				{token.Multiply, "*", 4, 8},
				{token.Delimiter, "", 4, 9},
				{token.Ident, "foo", 4, 10},
				{token.Delimiter, "", 4, 13},
				{token.Int, "2", 4, 14},
				{token.RParen, ")", 4, 15},
				{token.RParen, ")", 4, 16},
				{token.Delimiter, "", 4, 17},
				{token.LParen, "(", 5, 0},
				{token.Reassign, "=", 5, 1},
				{token.Delimiter, "", 5, 2},
				{token.Ident, "foo", 5, 3},
				{token.Delimiter, "", 5, 6},
				{token.LParen, "(", 5, 7},
				{token.Divide, "/", 5, 8},
				{token.Delimiter, "", 5, 9},
				{token.Ident, "foo", 5, 10},
				{token.Delimiter, "", 5, 13},
				{token.Int, "2", 5, 14},
				{token.RParen, ")", 5, 15},
				{token.RParen, ")", 5, 16},
				{token.Delimiter, "", 5, 17},
				{token.LParen, "(", 6, 0},
				{token.Reassign, "=", 6, 1},
				{token.Delimiter, "", 6, 2},
				{token.Ident, "foo", 6, 3},
				{token.Delimiter, "", 6, 6},
				{token.LParen, "(", 6, 7},
				{token.Modulo, "%", 6, 8},
				{token.Delimiter, "", 6, 9},
				{token.Ident, "foo", 6, 10},
				{token.Delimiter, "", 6, 13},
				{token.Int, "2", 6, 14},
				{token.RParen, ")", 6, 15},
				{token.RParen, ")", 6, 16},
				{token.Delimiter, "", 6, 17},
				{token.LParen, "(", 7, 0},
				{token.Assign, ":=", 7, 1},
				{token.Delimiter, "", 7, 3},
				{token.Ident, "isBar", 7, 4},
				{token.Delimiter, "", 7, 9},
				{token.LParen, "(", 7, 10},
				{token.Equals, "==", 7, 11},
				{token.Delimiter, "", 7, 13},
				{token.Ident, "foo", 7, 14},
				{token.Delimiter, "", 7, 17},
				{token.Int, "2", 7, 18},
				{token.RParen, ")", 7, 19},
				{token.RParen, ")", 7, 20},
				{token.Delimiter, "", 7, 21},
				{token.LParen, "(", 8, 0},
				{token.Reassign, "=", 8, 1},
				{token.Delimiter, "", 8, 2},
				{token.Ident, "isBar", 8, 3},
				{token.Delimiter, "", 8, 8},
				{token.LParen, "(", 8, 9},
				{token.NotEquals, "!=", 8, 10},
				{token.Delimiter, "", 8, 12},
				{token.Ident, "foo", 8, 13},
				{token.Delimiter, "", 8, 16},
				{token.Int, "3", 8, 17},
				{token.RParen, ")", 8, 18},
				{token.RParen, ")", 8, 19},
				{token.Delimiter, "", 8, 20},
				{token.LParen, "(", 9, 0},
				{token.Reassign, "=", 9, 1},
				{token.Delimiter, "", 9, 2},
				{token.Ident, "isBar", 9, 3},
				{token.Delimiter, "", 9, 8},
				{token.LParen, "(", 9, 9},
				{token.LessThanOrEquals, "<=", 9, 10},
				{token.Delimiter, "", 9, 12},
				{token.Ident, "foo", 9, 13},
				{token.Delimiter, "", 9, 16},
				{token.Int, "4", 9, 17},
				{token.RParen, ")", 9, 18},
				{token.RParen, ")", 9, 19},
				{token.Delimiter, "", 9, 20},
				{token.LParen, "(", 10, 0},
				{token.Reassign, "=", 10, 1},
				{token.Delimiter, "", 10, 2},
				{token.Ident, "isBar", 10, 3},
				{token.Delimiter, "", 10, 8},
				{token.LParen, "(", 10, 9},
				{token.GreaterThanOrEquals, ">=", 10, 10},
				{token.Delimiter, "", 10, 12},
				{token.Ident, "foo", 10, 13},
				{token.Delimiter, "", 10, 16},
				{token.Int, "5", 10, 17},
				{token.RParen, ")", 10, 18},
				{token.RParen, ")", 10, 19},
				{token.Delimiter, "", 10, 20},
				{token.LParen, "(", 11, 0},
				{token.Reassign, "=", 11, 1},
				{token.Delimiter, "", 11, 2},
				{token.Ident, "isBar", 11, 3},
				{token.Delimiter, "", 11, 8},
				{token.LParen, "(", 11, 9},
				{token.LessThan, "<", 11, 10},
				{token.Delimiter, "", 11, 11},
				{token.Ident, "foo", 11, 12},
				{token.Delimiter, "", 11, 15},
				{token.Int, "3", 11, 16},
				{token.RParen, ")", 11, 17},
				{token.RParen, ")", 11, 18},
				{token.Delimiter, "", 11, 19},
				{token.LParen, "(", 12, 0},
				{token.Reassign, "=", 12, 1},
				{token.Delimiter, "", 12, 2},
				{token.Ident, "isBar", 12, 3},
				{token.Delimiter, "", 12, 8},
				{token.LParen, "(", 12, 9},
				{token.GreaterThan, ">", 12, 10},
				{token.Delimiter, "", 12, 11},
				{token.Ident, "foo", 12, 12},
				{token.Delimiter, "", 12, 15},
				{token.Int, "3", 12, 16},
				{token.RParen, ")", 12, 17},
				{token.RParen, ")", 12, 18},
				{token.EOF, "", 12, 18},
			},
		},
		"Illegal token": {
			input: `(:= foo 5555xxxx)`,
			want: []token.Token{
				{token.LParen, "(", 1, 0},
				{token.Assign, ":=", 1, 1},
				{token.Delimiter, "", 1, 3},
				{token.Ident, "foo", 1, 4},
				{token.Delimiter, "", 1, 7},
				{token.Illegal, "5555xxxx", 1, 8},
				{token.RParen, ")", 1, 16},
				{token.EOF, "", 1, 16},
			},
		},
		"boolean assign with logical operators": {
			input: `(:= foo false)
(! foo)
(&& foo true)
(|| foo false)`,
			want: []token.Token{
				{token.LParen, "(", 1, 0},
				{token.Assign, ":=", 1, 1},
				{token.Delimiter, "", 1, 3},
				{token.Ident, "foo", 1, 4},
				{token.Delimiter, "", 1, 7},
				{token.Bool, "false", 1, 8},
				{token.RParen, ")", 1, 13},
				{token.Delimiter, "", 1, 14},
				{token.LParen, "(", 2, 0},
				{token.Not, "!", 2, 1},
				{token.Delimiter, "", 2, 2},
				{token.Ident, "foo", 2, 3},
				{token.RParen, ")", 2, 6},
				{token.Delimiter, "", 2, 7},
				{token.LParen, "(", 3, 0},
				{token.And, "&&", 3, 1},
				{token.Delimiter, "", 3, 3},
				{token.Ident, "foo", 3, 4},
				{token.Delimiter, "", 3, 7},
				{token.Bool, "true", 3, 8},
				{token.RParen, ")", 3, 12},
				{token.Delimiter, "", 3, 13},
				{token.LParen, "(", 4, 0},
				{token.Or, "||", 4, 1},
				{token.Delimiter, "", 4, 3},
				{token.Ident, "foo", 4, 4},
				{token.Delimiter, "", 4, 7},
				{token.Bool, "false", 4, 8},
				{token.RParen, ")", 4, 13},
				{token.EOF, "", 4, 13},
			},
		},
		"control flow keywords": {
			input: "(if elif else)",
			want: []token.Token{
				{token.LParen, "(", 1, 0},
				{token.If, "if", 1, 1},
				{token.Delimiter, "", 1, 3},
				{token.Elif, "elif", 1, 4},
				{token.Delimiter, "", 1, 8},
				{token.Else, "else", 1, 9},
				{token.RParen, ")", 1, 13},
				{token.EOF, "", 1, 13},
			},
		},
		"multiline if-elif-else": {
			input: `(if a
    (= x true)
elif b
    (= x false)
else
    (= y true)
)`,
			want: []token.Token{
				{Type: token.LParen, Literal: "(", Line: 1, Col: 0},
				{Type: token.If, Literal: "if", Line: 1, Col: 1},
				{Type: token.Delimiter, Literal: "", Line: 1, Col: 3},
				{Type: token.Ident, Literal: "a", Line: 1, Col: 4},
				{Type: token.Delimiter, Literal: "", Line: 1, Col: 5},
				{Type: token.LParen, Literal: "(", Line: 2, Col: 4},
				{Type: token.Reassign, Literal: "=", Line: 2, Col: 5},
				{Type: token.Delimiter, Literal: "", Line: 2, Col: 6},
				{Type: token.Ident, Literal: "x", Line: 2, Col: 7},
				{Type: token.Delimiter, Literal: "", Line: 2, Col: 8},
				{Type: token.Bool, Literal: "true", Line: 2, Col: 9},
				{Type: token.RParen, Literal: ")", Line: 2, Col: 13},
				{Type: token.Delimiter, Literal: "", Line: 2, Col: 14},
				{Type: token.Elif, Literal: "elif", Line: 3, Col: 0},
				{Type: token.Delimiter, Literal: "", Line: 3, Col: 4},
				{Type: token.Ident, Literal: "b", Line: 3, Col: 5},
				{Type: token.Delimiter, Literal: "", Line: 3, Col: 6},
				{Type: token.LParen, Literal: "(", Line: 4, Col: 4},
				{Type: token.Reassign, Literal: "=", Line: 4, Col: 5},
				{Type: token.Delimiter, Literal: "", Line: 4, Col: 6},
				{Type: token.Ident, Literal: "x", Line: 4, Col: 7},
				{Type: token.Delimiter, Literal: "", Line: 4, Col: 8},
				{Type: token.Bool, Literal: "false", Line: 4, Col: 9},
				{Type: token.RParen, Literal: ")", Line: 4, Col: 14},
				{Type: token.Delimiter, Literal: "", Line: 4, Col: 15},
				{Type: token.Else, Literal: "else", Line: 5, Col: 0},
				{Type: token.Delimiter, Literal: "", Line: 5, Col: 4},
				{Type: token.LParen, Literal: "(", Line: 6, Col: 4},
				{Type: token.Reassign, Literal: "=", Line: 6, Col: 5},
				{Type: token.Delimiter, Literal: "", Line: 6, Col: 6},
				{Type: token.Ident, Literal: "y", Line: 6, Col: 7},
				{Type: token.Delimiter, Literal: "", Line: 6, Col: 8},
				{Type: token.Bool, Literal: "true", Line: 6, Col: 9},
				{Type: token.RParen, Literal: ")", Line: 6, Col: 13},
				{Type: token.Delimiter, Literal: "", Line: 6, Col: 14},
				{Type: token.RParen, Literal: ")", Line: 7, Col: 0},
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
