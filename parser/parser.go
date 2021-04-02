package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	// Define order of precedence of operators
	// first one is 0, then 1 to 7
	_           int = iota
	LOWEST          // precedence with the lowest priority
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // +
	PRODUCT         // *
	PREFIX          // -X or !X
	CALL            // myFunction(X)
)

// map operator tokens to their desired precedence
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// Parser is the parser object for the Monkey language
type Parser struct {
	l      *lexer.Lexer
	errors []string

	// Store current and next tokens
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New initialises a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Initialise the prefix parse functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	// Initialise the infix parse functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances and sets the relevant tokens in the Parser object
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken() // calls the lexer's NextToken and stores it
}

// ParseProgram parses the statements and returns a working AST object
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// Advance until the end of the file
	for !p.curTokenIs(token.EOF) {
		// Get a statement and add it to the program's list
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)

		// advance
		p.nextToken()
	}
	return program
}

// curTokenIs returns if the parser's current token is of the given TokenType
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs returns if the parser's next token is of the given TokenType
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek asserts the correctness of order of tokens by checking the next type
func (p *Parser) expectPeek(t token.TokenType) bool {
	if !p.peekTokenIs(t) {
		p.peekError(t)
		return false
	}
	p.nextToken()
	return true
}

// peekError appends an error to p.errors
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Errors returns the errors field of the Parser
func (p *Parser) Errors() []string {
	return p.errors
}

// parseStatement returns a ast.Statement depending on the current token's type
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement returns an instance of *ast.LetStatement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// TODO: We're skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseLetStatement returns an instance of *ast.LetStatement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: We're skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseExpressionStatement returns an instance of *ast.ExpressionStatement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	// skip the semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseExpression return an expression node.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Check if there is a prefixParse function associated to the current token's type.
	prefixFn := p.prefixParseFns[p.curToken.Type]
	if prefixFn == nil {
		// attaches an error message to the parser for later use
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	// If the function exists, run it and assign its result
	// to a local variable representing the expression
	// on the left of the operator
	leftExp := prefixFn()

	// Main loop implementing Vaughan Pratt's "Top down operator precedence".
	// https://tdop.github.io/

	// Unless the next token is a semicolon and the precedence (an int) is smaller
	// then the the one of the next token:
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// check if there is an infixParseFn for the type of the next token
		infixFn := p.infixParseFns[p.peekToken.Type]
		if infixFn == nil {
			// If there is not, return just the left expression
			return leftExp
		}
		// If there is, move to the next token
		p.nextToken()
		// run the infix parse function on the current left expression,
		// and assign its result to leftExp
		leftExp = infixFn(leftExp)
		// Continue doing this until we reach a semicolon or the precedence changes
	}

	return leftExp
}

// helper for clearer parse expression errors
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// return an identifier with the current token and value
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// return an ast node which represents an expression integer literal
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as an integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// return an ast node which represents a prefix expression
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// return an ast node which represents an infix expression
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// Initialise an expression with the token, operator and the expression
	// which is the operand on the left of the operator.
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	// Take the current precedence, and advance to the next token
	precedence := p.curPrecedence()
	p.nextToken()
	// recursively call parseExpression() to generate the expression on the
	// right of the operator, with the correct precedence
	expression.Right = p.parseExpression(precedence)

	return expression
}

// return an ast node which represents a boolean expression
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// registerPrefix is a helper method that adds a prefixParseFn to the map
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix is a helper method that adds a infixParseFn to the map
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Returns the precedence associated with the token type in peekToken
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Return the precedence associated with th current token type
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
