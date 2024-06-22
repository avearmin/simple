package parser

import (
	"fmt"

	"github.com/avearmin/simple/internal/ast"
	"github.com/avearmin/simple/internal/lexer"
	"github.com/avearmin/simple/internal/token"
)

type Parser struct {
	l         *lexer.Lexer
	errors    []string
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

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	if !p.expectCur(token.LParen) {
		err := fmt.Sprintf("expected token '(', but got '%s'", p.curToken.Type)
		p.errors = append(p.errors, err)
		return nil
	}
	p.nextToken()

	switch p.curToken.Type {
	case token.Assign:
		stmt, err := p.parseAssignStatement()
		if err != nil {
			p.errors = append(p.errors, err.Error())
			return nil
		}
		return stmt
	default:
		return nil
	}
}

func (p *Parser) parseAssignStatement() (ast.AssignStatement, error) {
	stmt := ast.AssignStatement{Token: p.curToken}

	if !p.expectPeek(token.Delimiter) {
		return ast.AssignStatement{}, fmt.Errorf("expected token 'DELIMITER', but got '%s'", p.peekToken.Type)
	}
	p.nextToken()

	if !p.expectPeek(token.Ident) {
		return ast.AssignStatement{}, fmt.Errorf("expected token 'IDENT', but got '%s'", p.peekToken.Type)
	}
	p.nextToken()

	stmt.Name = ast.Atom{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.Delimiter) {
		return ast.AssignStatement{}, fmt.Errorf("expected token 'DELIMITER', but got '%s'", p.peekToken.Type)
	}
	p.nextToken()

	stmt.Value = p.parseExpression() // TODO: ADD THIS FUNCTION

	if !p.expectCur(token.RParen) {
		return ast.AssignStatement{}, fmt.Errorf("expected token ')', but got '%s'", p.peekToken.Type)
	}

	return stmt, nil
}

func (p *Parser) expectCur(tokType token.Type) bool {
	return tokType == p.curToken.Type
}

func (p *Parser) expectPeek(tokType token.Type) bool {
	return tokType == p.peekToken.Type
}
