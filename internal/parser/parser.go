package parser

import (
	"fmt"

	"github.com/avearmin/simple/internal/ast"
	"github.com/avearmin/simple/internal/lexer"
	"github.com/avearmin/simple/internal/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		program.Statements = append(program.Statements, stmt)

		p.nextToken()
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	if p.expectCur(token.Delimiter) {
		p.nextToken()
	}
	if !p.expectCur(token.LParen) {
		return nil, fmt.Errorf("expected token '(' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
	p.nextToken()

	switch p.curToken.Type {
	case token.Assign:
		stmt, err := p.parseAssignStatement()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	case token.Reassign:
		stmt, err := p.parseReassignStatement()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	default:
		return nil, fmt.Errorf("cannot begin statement with token '%s' on line %d col %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Col)
	}
}

func (p *Parser) parseAssignStatement() (ast.AssignStatement, error) {
	stmt := ast.AssignStatement{Token: p.curToken}

	if !p.expectPeek(token.Delimiter) {
		return ast.AssignStatement{}, fmt.Errorf("expected token 'DELIMITER' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.peekToken.Type)
	}
	p.nextToken()

	if !p.expectPeek(token.Ident) {
		return ast.AssignStatement{}, fmt.Errorf("expected token 'IDENT' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.peekToken.Type)
	}
	p.nextToken()

	stmt.Name = p.parseAtomExpression()

	if !p.expectPeek(token.Delimiter) {
		return ast.AssignStatement{}, fmt.Errorf("expected token 'DELIMITER' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.peekToken.Type)
	}
	p.nextToken()
	p.nextToken()

	exp, err := p.parseExpression()
	if err != nil {
		return ast.AssignStatement{}, err
	}
	stmt.Value = exp

	if !p.expectCur(token.RParen) {
		return ast.AssignStatement{}, fmt.Errorf("expected token ')' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.peekToken.Type)
	}

	return stmt, nil
}

func (p *Parser) parseReassignStatement() (ast.ReassignStatement, error) {
	stmt := ast.ReassignStatement{Token: p.curToken}

	if !p.expectPeek(token.Delimiter) {
		return ast.ReassignStatement{}, fmt.Errorf("expected token 'DELIMITER' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.peekToken.Type)
	}
	p.nextToken()

	if !p.expectPeek(token.Ident) {
		return ast.ReassignStatement{}, fmt.Errorf("expected token 'IDENT' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.peekToken.Type)
	}
	p.nextToken()

	stmt.Name = p.parseAtomExpression()

	if !p.expectPeek(token.Delimiter) {
		return ast.ReassignStatement{}, fmt.Errorf("expected token 'DELIMITER' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.peekToken.Type)
	}
	p.nextToken()
	p.nextToken()

	exp, err := p.parseExpression()
	if err != nil {
		return ast.ReassignStatement{}, err
	}
	stmt.Value = exp

	if !p.expectCur(token.RParen) {
		return ast.ReassignStatement{}, fmt.Errorf("expected token ')' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.peekToken.Type)
	}

	return stmt, nil
}

// expressions are either lists, or atoms
func (p *Parser) parseExpression() (ast.Expression, error) {
	switch p.curToken.Type {
	case token.LParen:
		p.nextToken() // advance off the RParen token
		exp, err := p.parseListExpression()
		if err != nil {
			return nil, err
		}
		return exp, nil
	case token.Ident, token.Int:
		return p.parseAtomExpression(), nil
	default:
		err := fmt.Errorf("Unexpected token '%s' on line %d col %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Col)
		return nil, err
	}
}

func (p *Parser) parseListExpression() (ast.Expression, error) {
	switch p.curToken.Type {
	case "+", "-", "*", "/", "%":
		exp, err := p.parseBinaryExpression()
		if err != nil {
			return nil, err
		}
		return exp, nil
	default:
		return nil, fmt.Errorf("Unexpected token on line %d col %d in list expression '%s'",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
}

func (p *Parser) parseBinaryExpression() (ast.BinaryExpression, error) {
	var binaryExp ast.BinaryExpression
	switch p.curToken.Type {
	case "+", "-", "*", "/", "%":
		binaryExp.Token = p.curToken
	default:
		return ast.BinaryExpression{}, fmt.Errorf("Expected tokens '+' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}

	if !p.expectPeek(token.Delimiter) {
		return ast.BinaryExpression{}, fmt.Errorf("Expected token 'DELIMITER on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.peekToken.Type)
	}
	p.nextToken()
	p.nextToken()

	expOne, err := p.parseExpression()
	if err != nil {
		return ast.BinaryExpression{}, err
	}

	if !p.expectPeek(token.Delimiter) {
		return ast.BinaryExpression{}, fmt.Errorf("Expected token 'DELIMITER on line %d col %d, but got '%s'",
			p.peekToken.Line, p.peekToken.Col, p.peekToken.Type)
	}
	p.nextToken()
	p.nextToken()

	expTwo, err := p.parseExpression()
	if err != nil {
		return ast.BinaryExpression{}, err
	}

	if !p.expectPeek(token.RParen) {
		return ast.BinaryExpression{}, fmt.Errorf("expected token ')' on line %d col %d, but got '%s'",
			p.peekToken.Line, p.peekToken.Col, p.peekToken.Type)
	}
	p.nextToken()
	p.nextToken()

	binaryExp.First = expOne
	binaryExp.Second = expTwo

	return binaryExp, nil
}

func (p *Parser) parseAtomExpression() ast.Atom {
	return ast.Atom{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) expectCur(tokType token.Type) bool {
	return tokType == p.curToken.Type
}

func (p *Parser) expectPeek(tokType token.Type) bool {
	return tokType == p.peekToken.Type
}
