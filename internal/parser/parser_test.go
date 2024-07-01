package parser

import (
	"fmt"
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
			input: `(:= foo (+ 1 2))
(= foo (- foo 1))
(= foo (* foo 2))
(= foo (/ foo 2))
(= foo (% foo 2))`,
			want: &ast.Program{
				Statements: []ast.Statement{
					ast.AssignStatement{
						Token: token.Token{Type: token.Assign, Literal: ":=", Line: 1, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "foo", Line: 1, Col: 4},
							Value: "foo",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.Add, Literal: "+", Line: 1, Col: 9},
							First: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "1", Line: 1, Col: 11},
								Value: "1",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "2", Line: 1, Col: 13},
								Value: "2",
							},
						},
					},
					ast.ReassignStatement{
						Token: token.Token{Type: token.Reassign, Literal: "=", Line: 2, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "foo", Line: 2, Col: 3},
							Value: "foo",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.Subtract, Literal: "-", Line: 2, Col: 8},
							First: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "foo", Line: 2, Col: 10},
								Value: "foo",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "1", Line: 2, Col: 14},
								Value: "1",
							},
						},
					},
					ast.ReassignStatement{
						Token: token.Token{Type: token.Reassign, Literal: "=", Line: 3, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "foo", Line: 3, Col: 3},
							Value: "foo",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.Multiply, Literal: "*", Line: 3, Col: 8},
							First: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "foo", Line: 3, Col: 10},
								Value: "foo",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "2", Line: 3, Col: 14},
								Value: "2",
							},
						},
					},
					ast.ReassignStatement{
						Token: token.Token{Type: token.Reassign, Literal: "=", Line: 4, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "foo", Line: 4, Col: 3},
							Value: "foo",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.Divide, Literal: "/", Line: 4, Col: 8},
							First: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "foo", Line: 4, Col: 10},
								Value: "foo",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "2", Line: 4, Col: 14},
								Value: "2",
							},
						},
					},
					ast.ReassignStatement{
						Token: token.Token{Type: token.Reassign, Literal: "=", Line: 5, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "foo", Line: 5, Col: 3},
							Value: "foo",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.Modulo, Literal: "%", Line: 5, Col: 8},
							First: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "foo", Line: 5, Col: 10},
								Value: "foo",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "2", Line: 5, Col: 14},
								Value: "2",
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

			program, err := p.ParseProgram()
			if err != nil {
				t.Fatalf("Parsing failed with error: %s", err)
			}
			if err := isEqualPrograms(program, test.want); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func isEqualPrograms(first, second *ast.Program) error {
	if len(first.Statements) != len(second.Statements) {
		return fmt.Errorf("expected statements len=%d, but got=%d", len(second.Statements), len(first.Statements))
	}

	for i := range first.Statements {
		if !isEqualStatements(first.Statements[i], second.Statements[i]) {
			return fmt.Errorf("expected statement=%v, but got=%v", second.Statements[i], first.Statements[i])
		}
	}

	return nil
}

func isEqualStatements(first, second ast.Statement) bool {
	switch stmtOne := first.(type) {
	case ast.AssignStatement:
		stmtTwo, ok := second.(ast.AssignStatement)
		if !ok {
			return false
		}
		return isEqualAssignStatement(stmtOne, stmtTwo)
	case ast.ReassignStatement:
		stmtTwo, ok := second.(ast.ReassignStatement)
		if !ok {
			return false
		}
		return isEqualReassignStatement(stmtOne, stmtTwo)
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

func isEqualReassignStatement(first, second ast.ReassignStatement) bool {
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
