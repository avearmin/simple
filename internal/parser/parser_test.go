package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/avearmin/simple/internal/ast"
	"github.com/avearmin/simple/internal/lexer"
	"github.com/avearmin/simple/internal/token"
)

type batchError struct {
	msgs []string
}

func newBatchError() *batchError {
	return &batchError{msgs: []string{}}
}

func (be *batchError) add(s string) {
	be.msgs = append(be.msgs, s)
}

func (be *batchError) Error() string {
	var builder strings.Builder

	for _, msg := range be.msgs {
		builder.WriteString(msg)
	}

	return builder.String()
}

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
(= foo (% foo 2))

(:= isBar (== foo 2))
(= isBar (!= foo 2))
(= isBar (<= foo 2))
(= isBar (>= foo 2))
(= isBar (> foo 2))
(= isBar (< foo 2))

(if isBar (= foo (+ foo 1))
elif (== foo 1) (= isBar false)
else (= isBar true))

(fn addThenDouble x y
    (:= z (+ x y))
    (return (* 2 z)))

(for (:= i 0) (< i 5) (= i (+ i 1))
    (addThenDouble i 2))`,
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
					ast.AssignStatement{
						Token: token.Token{Type: token.Assign, Literal: ":=", Line: 7, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "isBar", Line: 7, Col: 4},
							Value: "isBar",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.Equals, Literal: "==", Line: 7, Col: 11},
							First: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "foo", Line: 7, Col: 14},
								Value: "foo",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "2", Line: 7, Col: 18},
								Value: "2",
							},
						},
					},
					ast.ReassignStatement{
						Token: token.Token{Type: token.Reassign, Literal: "=", Line: 8, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "isBar", Line: 8, Col: 3},
							Value: "isBar",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.NotEquals, Literal: "!=", Line: 8, Col: 10},
							First: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "foo", Line: 8, Col: 13},
								Value: "foo",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "2", Line: 8, Col: 17},
								Value: "2",
							},
						},
					},
					ast.ReassignStatement{
						Token: token.Token{Type: token.Reassign, Literal: "=", Line: 9, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "isBar", Line: 9, Col: 3},
							Value: "isBar",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.LessThanOrEquals, Literal: "<=", Line: 9, Col: 10},
							First: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "foo", Line: 9, Col: 13},
								Value: "foo",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "2", Line: 9, Col: 17},
								Value: "2",
							},
						},
					},
					ast.ReassignStatement{
						Token: token.Token{Type: token.Reassign, Literal: "=", Line: 10, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "isBar", Line: 10, Col: 3},
							Value: "isBar",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.GreaterThanOrEquals, Literal: ">=", Line: 10, Col: 10},
							First: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "foo", Line: 10, Col: 13},
								Value: "foo",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "2", Line: 10, Col: 17},
								Value: "2",
							},
						},
					},
					ast.ReassignStatement{
						Token: token.Token{Type: token.Reassign, Literal: "=", Line: 11, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "isBar", Line: 11, Col: 3},
							Value: "isBar",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.GreaterThan, Literal: ">", Line: 11, Col: 10},
							First: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "foo", Line: 11, Col: 12},
								Value: "foo",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "2", Line: 11, Col: 16},
								Value: "2",
							},
						},
					},
					ast.ReassignStatement{
						Token: token.Token{Type: token.Reassign, Literal: "=", Line: 12, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "isBar", Line: 12, Col: 3},
							Value: "isBar",
						},
						Value: ast.BinaryExpression{
							Token: token.Token{Type: token.LessThan, Literal: "<", Line: 12, Col: 10},
							First: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "foo", Line: 12, Col: 12},
								Value: "foo",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "2", Line: 12, Col: 16},
								Value: "2",
							},
						},
					},
					ast.ConditionalStatement{
						Token: token.Token{Type: token.If, Literal: "if", Line: 14, Col: 1},
						IfCondition: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "isBar", Line: 14, Col: 4},
							Value: "isBar",
						},
						IfStatements: []ast.Statement{
							ast.ReassignStatement{
								Token: token.Token{Type: token.Reassign, Literal: "=", Line: 14, Col: 11},
								Name: ast.Atom{
									Token: token.Token{Type: token.Ident, Literal: "foo", Line: 14, Col: 13},
									Value: "foo",
								},
								Value: ast.BinaryExpression{
									Token: token.Token{Type: token.Add, Literal: "+", Line: 14, Col: 18},
									First: ast.Atom{
										Token: token.Token{Type: token.Ident, Literal: "foo", Line: 14, Col: 20},
										Value: "foo",
									},
									Second: ast.Atom{
										Token: token.Token{Type: token.Int, Literal: "1", Line: 14, Col: 24},
										Value: "1",
									},
								},
							},
						},
						ElifBlocks: []ast.ElifBlock{
							{
								Token: token.Token{Type: token.Elif, Literal: "elif", Line: 15, Col: 0},
								Condition: ast.BinaryExpression{
									Token: token.Token{Type: token.Equals, Literal: "==", Line: 15, Col: 6},
									First: ast.Atom{
										Token: token.Token{Type: token.Ident, Literal: "foo", Line: 15, Col: 9},
										Value: "foo",
									},
									Second: ast.Atom{
										Token: token.Token{Type: token.Int, Literal: "1", Line: 15, Col: 13},
										Value: "1",
									},
								},
								Statements: []ast.Statement{
									ast.ReassignStatement{
										Token: token.Token{Type: token.Reassign, Literal: "=", Line: 15, Col: 17},
										Name: ast.Atom{
											Token: token.Token{Type: token.Ident, Literal: "isBar", Line: 15, Col: 19},
											Value: "isBar",
										},
										Value: ast.Atom{
											Token: token.Token{Type: token.Bool, Literal: "false", Line: 15, Col: 25},
											Value: "false",
										},
									},
								},
							},
						},
						ElseBlock: ast.ElseBlock{
							Token: token.Token{Type: token.Else, Literal: "else", Line: 16, Col: 0},
							Statements: []ast.Statement{
								ast.ReassignStatement{
									Token: token.Token{Type: token.Reassign, Literal: "=", Line: 16, Col: 6},
									Name: ast.Atom{
										Token: token.Token{Type: token.Ident, Literal: "isBar", Line: 16, Col: 8},
										Value: "isBar",
									},
									Value: ast.Atom{
										Token: token.Token{Type: token.Bool, Literal: "true", Line: 16, Col: 14},
										Value: "true",
									},
								},
							},
						},
					},
					ast.FunctionAssignStatement{
						Token: token.Token{Type: token.Fn, Literal: "fn", Line: 18, Col: 1},
						Name: ast.Atom{
							Token: token.Token{Type: token.Ident, Literal: "addThenDouble", Line: 18, Col: 4},
							Value: "addThenDouble",
						},
						Params: []ast.Atom{
							{
								Token: token.Token{Type: token.Ident, Literal: "x", Line: 18, Col: 18},
								Value: "x",
							},
							{
								Token: token.Token{Type: token.Ident, Literal: "y", Line: 18, Col: 20},
								Value: "y",
							},
						},
						Statements: []ast.Statement{
							ast.AssignStatement{
								Token: token.Token{Type: token.Assign, Literal: ":=", Line: 19, Col: 5},
								Name: ast.Atom{
									Token: token.Token{Type: token.Ident, Literal: "z", Line: 19, Col: 8},
									Value: "z",
								},
								Value: ast.BinaryExpression{
									Token: token.Token{Type: token.Add, Literal: "+", Line: 19, Col: 11},
									First: ast.Atom{
										Token: token.Token{Type: token.Ident, Literal: "x", Line: 19, Col: 13},
										Value: "x",
									},
									Second: ast.Atom{
										Token: token.Token{Type: token.Ident, Literal: "y", Line: 19, Col: 15},
										Value: "y",
									},
								},
							},
							ast.ReturnStatement{
								Token: token.Token{Type: token.Return, Literal: "return", Line: 20, Col: 5},
								Value: ast.BinaryExpression{
									Token: token.Token{Type: token.Multiply, Literal: "*", Line: 20, Col: 13},
									First: ast.Atom{
										Token: token.Token{Type: token.Int, Literal: "2", Line: 20, Col: 15},
										Value: "2",
									},
									Second: ast.Atom{
										Token: token.Token{Type: token.Ident, Literal: "2", Line: 20, Col: 17},
										Value: "z",
									},
								},
							},
						},
					},
					ast.ForLoopStatement{
						Token: token.Token{Type: token.For, Literal: "for", Line: 22, Col: 1},
						Initalizer: ast.AssignStatement{
							Token: token.Token{Type: token.Assign, Literal: ":=", Line: 22, Col: 6},
							Name: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "i", Line: 22, Col: 9},
								Value: "i",
							},
							Value: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "0", Line: 22, Col: 11},
								Value: "0",
							},
						},
						Condition: ast.BinaryExpression{
							Token: token.Token{Type: token.LessThan, Literal: "<", Line: 22, Col: 15},
							First: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "i", Line: 22, Col: 17},
								Value: "i",
							},
							Second: ast.Atom{
								Token: token.Token{Type: token.Int, Literal: "5", Line: 22, Col: 19},
								Value: "5",
							},
						},
						Update: ast.ReassignStatement{
							Token: token.Token{Type: token.Reassign, Literal: "=", Line: 22, Col: 23},
							Name: ast.Atom{
								Token: token.Token{Type: token.Ident, Literal: "i", Line: 22, Col: 25},
								Value: "i",
							},
							Value: ast.BinaryExpression{
								Token: token.Token{Type: token.Add, Literal: "+", Line: 22, Col: 28},
								First: ast.Atom{
									Token: token.Token{Type: token.Ident, Literal: "i", Line: 22, Col: 30},
									Value: "i",
								},
								Second: ast.Atom{
									Token: token.Token{Type: token.Int, Literal: "1", Line: 22, Col: 32},
									Value: "1",
								},
							},
						},
						Statements: []ast.Statement{
							ast.FnCall{
								Token: token.Token{Type: token.Ident, Literal: "addThenDouble", Line: 23, Col: 5},
								Arguments: []ast.Atom{
									{
										Token: token.Token{Type: token.Ident, Literal: "i", Line: 23, Col: 19},
										Value: "i",
									},
									{
										Token: token.Token{Type: token.Int, Literal: "2", Line: 23, Col: 21},
										Value: "2",
									},
								},
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
				t.Fatalf("Parsing failed with error: %s", err.Error())
			}
			if !isEqualPrograms(program, test.want) {
				if len(program.Statements) != len(test.want.Statements) {
					t.Fatalf("len(want)=%d, len(got)=%d", len(test.want.Statements), len(program.Statements))
				}

				bErr := newBatchError()
				for i := range program.Statements {
					bErr.add(fmt.Sprintf("got=%v\nwant=%v\n", program.Statements[i], test.want.Statements[i]))
				}
				t.Fatal(bErr)
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
	case ast.ReassignStatement:
		stmtTwo, ok := second.(ast.ReassignStatement)
		if !ok {
			return false
		}
		return isEqualReassignStatement(stmtOne, stmtTwo)
	case ast.ConditionalStatement:
		stmtTwo, ok := second.(ast.ConditionalStatement)
		if !ok {
			return false
		}
		return isEqualConditionalStatement(stmtOne, stmtTwo)
	case ast.FunctionAssignStatement:
		stmtTwo, ok := second.(ast.FunctionAssignStatement)
		if !ok {
			return false
		}
		return isEqualFunctionAssignStatements(stmtOne, stmtTwo)
	case ast.ReturnStatement:
		stmtTwo, ok := second.(ast.ReturnStatement)
		if !ok {
			return false
		}
		return isEqualReturnStatements(stmtOne, stmtTwo)
	case ast.FnCall:
		stmtTwo, ok := second.(ast.FnCall)
		if !ok {
			return false
		}
		return isEqualFnCalls(stmtOne, stmtTwo)
	case ast.ForLoopStatement:
		stmtTwo, ok := second.(ast.ForLoopStatement)
		if !ok {
			return false
		}
		return isEqualForLoopStatements(stmtOne, stmtTwo)
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

func isEqualConditionalStatement(first, second ast.ConditionalStatement) bool {
	if !isEqualTokens(first.Token, second.Token) {
		return false
	}
	if !isEqualExpressions(first.IfCondition, second.IfCondition) {
		return false
	}
	if len(first.IfStatements) != len(second.IfStatements) {
		return false
	}
	for i := range first.IfStatements {
		if !isEqualStatements(first.IfStatements[i], second.IfStatements[i]) {
			return false
		}
	}
	if len(first.ElifBlocks) != len(second.ElifBlocks) {
		return false
	}
	for i := range first.ElifBlocks {
		if !isEqualElifBlocks(first.ElifBlocks[i], second.ElifBlocks[i]) {
			return false
		}
	}
	return isEqualElseBlocks(first.ElseBlock, second.ElseBlock)
}

func isEqualFunctionAssignStatements(first, second ast.FunctionAssignStatement) bool {
	if !isEqualTokens(first.Token, second.Token) {
		return false
	}

	if !isEqualAtoms(first.Name, second.Name) {
		return false
	}

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

func isEqualReturnStatements(first, second ast.ReturnStatement) bool {
	if !isEqualTokens(first.Token, second.Token) {
		return false
	}

	if !isEqualExpressions(first.Value, second.Value) {
		return false
	}

	return true
}

func isEqualFnCalls(first, second ast.FnCall) bool {
	if !isEqualTokens(first.Token, second.Token) {
		return false
	}

	if len(first.Arguments) != len(second.Arguments) {
		return false
	}

	for i := range first.Arguments {
		if !isEqualAtoms(first.Arguments[i], second.Arguments[i]) {
			return false
		}
	}

	return true
}

func isEqualForLoopStatements(first, second ast.ForLoopStatement) bool {
	if !isEqualTokens(first.Token, second.Token) {
		return false
	}

	if !isEqualAssignStatement(first.Initalizer, second.Initalizer) {
		return false
	}

	if !isEqualBinaryExpressions(first.Condition, second.Condition) {
		return false
	}

	if !isEqualReassignStatement(first.Update, second.Update) {
		return false
	}

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
	if first.Line != second.Line {
		return false
	}
	if first.Col != second.Col {
		return false
	}
	return true
}

func isEqualElifBlocks(first, second ast.ElifBlock) bool {
	if !isEqualTokens(first.Token, second.Token) {
		return false
	}
	if !isEqualExpressions(first.Condition, second.Condition) {
		return false
	}
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

func isEqualElseBlocks(first, second ast.ElseBlock) bool {
	if !isEqualTokens(first.Token, second.Token) {
		return false
	}
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
