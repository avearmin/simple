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

		p.ignoreDelimiters()
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
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
	case token.If:
		stmt, err := p.parseConditionalStatement()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	case token.Fn:
		stmt, err := p.parseFunctionAssignStatement()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	case token.Return:
		stmt, err := p.parseReturnStatement()
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
	if !p.expectCur(token.Assign) {
		return ast.AssignStatement{}, fmt.Errorf("expected token ':=' on line %d col %d but got '%s'",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
	stmt := ast.AssignStatement{Token: p.curToken}
	p.nextToken()

	if err := p.eatDelimiter(); err != nil {
		return ast.AssignStatement{}, err
	}

	atom, err := p.parseAtomExpression()
	if err != nil {
		return ast.AssignStatement{}, err
	}
	if atom.TokenType() != token.Ident {
		return ast.AssignStatement{}, fmt.Errorf("expected token 'IDENT' on line %d col %d, but got %s",
			atom.Token.Line, atom.Token.Col, atom.TokenType())
	}
	stmt.Name = atom

	if err := p.eatDelimiter(); err != nil {
		return ast.AssignStatement{}, err
	}

	exp, err := p.parseExpression()
	if err != nil {
		return ast.AssignStatement{}, err
	}
	stmt.Value = exp

	if !p.expectCur(token.RParen) {
		return ast.AssignStatement{}, fmt.Errorf("expected token ')' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.peekToken.Type)
	}

	p.nextToken()

	return stmt, nil
}

func (p *Parser) parseReassignStatement() (ast.ReassignStatement, error) {
	if !p.expectCur(token.Reassign) {
		return ast.ReassignStatement{}, fmt.Errorf("expected token '=' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.peekToken.Type)
	}
	stmt := ast.ReassignStatement{Token: p.curToken}
	p.nextToken()

	if err := p.eatDelimiter(); err != nil {
		return ast.ReassignStatement{}, err
	}

	atom, err := p.parseAtomExpression()
	if err != nil {
		return ast.ReassignStatement{}, err
	}
	if atom.TokenType() != token.Ident {
		return ast.ReassignStatement{}, fmt.Errorf("expected token 'IDENT' on line %d col %d, but got %s",
			atom.Token.Line, atom.Token.Col, atom.TokenType())
	}
	stmt.Name = atom

	if err := p.eatDelimiter(); err != nil {
		return ast.ReassignStatement{}, err
	}

	exp, err := p.parseExpression()
	if err != nil {
		return ast.ReassignStatement{}, err
	}
	stmt.Value = exp

	if !p.expectCur(token.RParen) {
		return ast.ReassignStatement{}, fmt.Errorf("expected token ')' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
	p.nextToken()

	return stmt, nil
}

func (p *Parser) parseConditionalStatement() (ast.ConditionalStatement, error) {
	if !p.expectCur(token.If) {
		return ast.ConditionalStatement{}, fmt.Errorf("expected token 'IF' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
	stmt := ast.ConditionalStatement{
		Token:        p.curToken,
		IfStatements: []ast.Statement{},
		ElifBlocks:   []ast.ElifBlock{},
	}
	p.nextToken()

	if err := p.eatDelimiter(); err != nil {
		return ast.ConditionalStatement{}, err
	}

	ifCondExp, err := p.parseExpression()
	if err != nil {
		return ast.ConditionalStatement{}, err
	}
	stmt.IfCondition = ifCondExp

	if err := p.eatDelimiter(); err != nil {
		return ast.ConditionalStatement{}, err
	}

	for !p.expectCur(token.Elif) && !p.expectCur(token.Else) {
		ifStmt, err := p.parseStatement()
		if err != nil {
			return ast.ConditionalStatement{}, err
		}
		stmt.IfStatements = append(stmt.IfStatements, ifStmt)

		if p.expectCur(token.RParen) {
			return stmt, nil
		}

		if err := p.eatDelimiter(); err != nil {
			return ast.ConditionalStatement{}, err
		}
	}

	for !p.expectCur(token.Else) {
		elifBlock, err := p.parseElifBlock()
		if err != nil {
			return ast.ConditionalStatement{}, err
		}
		stmt.ElifBlocks = append(stmt.ElifBlocks, elifBlock)

		if p.expectCur(token.RParen) {
			return stmt, nil
		}
	}

	elseBlock, err := p.parseElseBlock()
	if err != nil {
		return ast.ConditionalStatement{}, err
	}
	stmt.ElseBlock = elseBlock

	if !p.expectCur(token.RParen) {
		return ast.ConditionalStatement{}, fmt.Errorf("expected token ')' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}

	p.nextToken()

	return stmt, nil
}

func (p *Parser) parseElifBlock() (ast.ElifBlock, error) {
	if !p.expectCur(token.Elif) {
		return ast.ElifBlock{}, fmt.Errorf("expected token 'ELIF' on line %d col %d but got %s",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}

	block := ast.ElifBlock{Token: p.curToken, Statements: []ast.Statement{}}
	p.nextToken()

	if err := p.eatDelimiter(); err != nil {
		return ast.ElifBlock{}, err
	}

	exp, err := p.parseExpression()
	if err != nil {
		return ast.ElifBlock{}, err
	}
	block.Condition = exp

	if err := p.eatDelimiter(); err != nil {
		return ast.ElifBlock{}, err
	}

	for p.curToken.Type != token.Elif && p.curToken.Type != token.Else && p.peekToken.Type != token.RParen {
		stmt, err := p.parseStatement()
		if err != nil {
			return ast.ElifBlock{}, err
		}
		block.Statements = append(block.Statements, stmt)
		p.eatDelimiter()
	}

	return block, nil
}

func (p *Parser) parseElseBlock() (ast.ElseBlock, error) {
	if !p.expectCur(token.Else) {
		return ast.ElseBlock{}, fmt.Errorf("expected token 'ELSE' on line %d col %d but got %s",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}

	block := ast.ElseBlock{Token: p.curToken, Statements: []ast.Statement{}}
	p.nextToken()

	if err := p.eatDelimiter(); err != nil {
		return ast.ElseBlock{}, err
	}

	for !p.expectCur(token.Elif) && !p.expectCur(token.Else) && !p.expectCur(token.RParen) {
		stmt, err := p.parseStatement()
		if err != nil {
			return ast.ElseBlock{}, err
		}
		block.Statements = append(block.Statements, stmt)
	}

	return block, nil
}

func (p *Parser) parseFunctionAssignStatement() (ast.FunctionAssignStatement, error) {
	if !p.expectCur(token.Fn) {
		return ast.FunctionAssignStatement{}, fmt.Errorf("%d:%d expected 'FN', got '%s'", p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
	fnStmt := ast.FunctionAssignStatement{Token: p.curToken, Params: []ast.Atom{}, Statements: []ast.Statement{}}
	p.nextToken()

	if err := p.eatDelimiter(); err != nil {
		return ast.FunctionAssignStatement{}, err
	}

	fnName, err := p.parseAtomExpression()
	if err != nil {
		return ast.FunctionAssignStatement{}, err
	}
	if fnName.TokenType() != token.Ident {
		return ast.FunctionAssignStatement{}, fmt.Errorf("%d:%d cannot use '%s' as name of a function", p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
	fnStmt.Name = fnName

	if err := p.eatDelimiter(); err != nil {
		return ast.FunctionAssignStatement{}, err
	}

	for !p.expectCur(token.LParen) {
		param, err := p.parseAtomExpression()
		if err != nil {
			return ast.FunctionAssignStatement{}, err
		}
		if param.TokenType() != token.Ident {
			return ast.FunctionAssignStatement{}, fmt.Errorf("%d:%d cannot use '%s' as parameter to a function", p.curToken.Line, p.curToken.Col, p.curToken.Type)

		}

		fnStmt.Params = append(fnStmt.Params, param)

		if p.expectCur(token.RParen) {
			return fnStmt, nil
		}

		if err := p.eatDelimiter(); err != nil {
			return ast.FunctionAssignStatement{}, err
		}

	}

	for {
		innerStmt, err := p.parseStatement()
		if err != nil {
			return ast.FunctionAssignStatement{}, err
		}
		fnStmt.Statements = append(fnStmt.Statements, innerStmt)

		if p.expectCur(token.RParen) {
			break
		}

		if err := p.eatDelimiter(); err != nil {
			return ast.FunctionAssignStatement{}, err
		}
	}

	p.nextToken()
	return fnStmt, nil
}

func (p *Parser) parseReturnStatement() (ast.ReturnStatement, error) {
	if !p.expectCur(token.Return) {
		return ast.ReturnStatement{}, fmt.Errorf("%d:%d expected 'RETURN' but got '%s'", p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
	returnStmt := ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	if err := p.eatDelimiter(); err != nil {
		return ast.ReturnStatement{}, err
	}

	exp, err := p.parseExpression()
	if err != nil {
		return ast.ReturnStatement{}, err
	}

	returnStmt.Value = exp

	if !p.expectCur(token.RParen) {
		return ast.ReturnStatement{}, fmt.Errorf("%d:%d expected ')' but got %s", p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
	p.nextToken()

	return returnStmt, nil
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
	case token.Ident, token.Int, token.Bool:
		atom, err := p.parseAtomExpression()
		if err != nil {
			return nil, err
		}
		return atom, nil
	default:
		err := fmt.Errorf("Unexpected token '%s' on line %d col %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Col)
		return nil, err
	}
}

func (p *Parser) parseListExpression() (ast.Expression, error) {
	switch p.curToken.Type {
	case "+", "-", "*", "/", "%", "==", "!=", "<=", ">=", "<", ">":
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
	case "+", "-", "*", "/", "%", "==", "!=", "<=", ">=", "<", ">":
		binaryExp.Token = p.curToken
	default:
		return ast.BinaryExpression{}, fmt.Errorf("Cannot begin binary expression with '%s' on line %d col %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Col)
	}
	p.nextToken()

	if err := p.eatDelimiter(); err != nil {
		return ast.BinaryExpression{}, err
	}

	expOne, err := p.parseExpression()
	if err != nil {
		return ast.BinaryExpression{}, err
	}
	binaryExp.First = expOne

	if err := p.eatDelimiter(); err != nil {
		return ast.BinaryExpression{}, nil
	}

	expTwo, err := p.parseExpression()
	if err != nil {
		return ast.BinaryExpression{}, err
	}
	binaryExp.Second = expTwo

	if !p.expectCur(token.RParen) {
		return ast.BinaryExpression{}, fmt.Errorf("expected token ')' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
	p.nextToken()

	return binaryExp, nil
}

func (p *Parser) parseAtomExpression() (ast.Atom, error) {
	if !p.expectCur(token.Ident) && !p.expectCur(token.Int) && !p.expectCur(token.Bool) {
		return ast.Atom{}, fmt.Errorf("expected token 'IDENT' on line %d col %d, but got %s",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
	atom := ast.Atom{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	return atom, nil
}

func (p *Parser) expectCur(tokType token.Type) bool {
	return tokType == p.curToken.Type
}

func (p *Parser) expectPeek(tokType token.Type) bool {
	return tokType == p.peekToken.Type
}

func (p *Parser) eatDelimiter() error {
	if !p.expectCur(token.Delimiter) {
		return fmt.Errorf("expected token 'DELIMITER' on line %d col %d, but got '%s'",
			p.curToken.Line, p.curToken.Col, p.curToken.Type)
	}
	p.nextToken()
	return nil
}

func (p *Parser) ignoreDelimiters() {
	for p.expectCur(token.Delimiter) {
		p.nextToken()
	}
}
