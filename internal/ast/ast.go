package ast

import "github.com/avearmin/simple/internal/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type AssignStatement struct {
	Token token.Token
	Name  Atom
	Value Expression
}

func (as AssignStatement) statementNode()       {}
func (as AssignStatement) TokenLiteral() string { return as.Token.Literal }

type ReassignStatement struct {
	Token token.Token
	Name  Atom
	Value Expression
}

func (rs ReassignStatement) statementNode()       {}
func (rs ReassignStatement) TokenLiteral() string { return rs.Token.Literal }

type Atom struct {
	Token token.Token
	Value string
}

func (a Atom) expressionNode()      {}
func (a Atom) TokenLiteral() string { return a.Token.Literal }

type BinaryExpression struct {
	Token  token.Token
	First  Expression
	Second Expression
}

func (be BinaryExpression) expressionNode()      {}
func (be BinaryExpression) TokenLiteral() string { return be.Token.Literal }
