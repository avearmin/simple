package parser

import (
	"testing"

	"github.com/avearmin/simple/internal/ast"
	"github.com/avearmin/simple/internal/lexer"
	"github.com/avearmin/simple/internal/token"
)

func TestParseProgram(t *testing.T) {
	tests := map[string]struct {
		input string
		want  *ast.Program
	}{
		"simple program": {
			input: `(:= foo (+ 1 2))`,
			want: &ast.Program{
				Statements: []ast.Statement{
					ast.AssignStatement{
						Token: token.Token{Type: token.Assign, Literal: ":="},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "foo"},
							Value: "foo",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.Add, Literal: "+"},
							First: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "1"},
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "2"},
							},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			l := lexer.New(test.input)
			p := New(l)

			program := p.ParseProgram()
			if !isEqualPrograms(program, test.want) {
				t.Fail()
			}
		})
	}
}

func isEqualPrograms(first, second *ast.Program) bool {
	if len(first.Statements) != len(second.Statements) {
		return false
	}

	for i := range first.Statements {
		if !isEqualStatements(first.Statements[i], second.Statements[i]) {
			return false
		}
	}

	return true
}

func isEqualStatements(first, second ast.Statement) bool {
	switch stmtOne := first.(type) {
	case ast.AssignStatement:
		stmtTwo, ok := second.(ast.AssignStatement)
		if !ok {
			return false
		}
		return isEqualAssignStatement(stmtOne, stmtTwo)
	}
	return false
}

func isEqualExpressions(first, second ast.Expression) bool {
	switch expOne := first.(type) {
	case ast.BinaryExpression:
		expTwo, ok := second.(ast.BinaryExpression)
		if !ok {
			return false
		}
		return isEqualBinaryExpressions(expOne, expTwo)
	case ast.Atom:
		expTwo, ok := second.(ast.Atom)
		if !ok {
			return false
		}
		return isEqualAtoms(expOne, expTwo)
	}

	return false
}

func isEqualAssignStatement(first, second ast.AssignStatement) bool {
	if !isEqualTokens(first.Token, second.Token) {
		return false
	}
	if !isEqualAtoms(first.Name, second.Name) {
		return false
	}
	if !isEqualExpressions(first.Value, second.Value) {
		return false
	}
	return true
}

func isEqualAtoms(first, second ast.Atom) bool {
	return isEqualTokens(first.Token, second.Token) && first.Value == second.Value
}

func isEqualBinaryExpressions(first, second ast.BinaryExpression) bool {
	return isEqualTokens(first.Token, second.Token) && isEqualExpressions(first.First, second.First)
}

func isEqualTokens(first, second token.Token) bool {
	if first.Type != second.Type {
		return false
	}
	if first.Literal != second.Literal {
		return false
	}
	return true
}
