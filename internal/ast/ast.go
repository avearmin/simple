package ast

import "github.com/avearmin/simple/internal/token"

type Node interface {
	TokenLiteral() string
	TokenType() token.Type
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
func (as AssignStatement) TokenType() token.Type { return as.Token.Type }

type ReassignStatement struct {
	Token token.Token
	Name  Atom
	Value Expression
}

func (rs ReassignStatement) statementNode()       {}
func (rs ReassignStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs ReassignStatement) TokenType() token.Type { return rs.Token.Type }


type ElifBlock struct {
	Token      token.Token
	Condition  Expression
	Statements []Statement
}

type ElseBlock struct {
	Token      token.Token
	Statements []Statement
}

type ConditionalStatement struct {
	Token        token.Token
	IfCondition  Expression
	IfStatements []Statement
	ElifBlocks   []ElifBlock
	ElseBlock    ElseBlock
}

func (cs ConditionalStatement) statementNode()       {}
func (cs ConditionalStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs ConditionalStatement) TokenType() token.Type { return cs.Token.Type }


type Atom struct {
	Token token.Token
	Value string
}

func (a Atom) expressionNode()      {}
func (a Atom) TokenLiteral() string { return a.Token.Literal }
func (a Atom) TokenType() token.Type { return a.Token.Type }


type BinaryExpression struct {
	Token  token.Token
	First  Expression
	Second Expression
}

func (be BinaryExpression) expressionNode()      {}
func (be BinaryExpression) TokenLiteral() string { return be.Token.Literal }
func (be BinaryExpression) TokenType() token.Type { return be.Token.Type }

