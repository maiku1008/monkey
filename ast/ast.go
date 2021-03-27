// Package ast has the necessary facilities to represent
// an Abstract Syntax Tree for the Monkey language
package ast

import (
	"bytes"
	"monkey/token"
)

// Node represents a single node in the AST.
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is a type of Node representing a statement
type Statement interface {
	Node
	statementNode()
}

// Expression is a type of Node representing an expression
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every AST of the parser
// Every valid monkey program is a series of statements
type Program struct {
	Statements []Statement
}

// Assert implementations
var _ Node = (*Program)(nil)

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String creates a buffer and writes the return value
// of each statement's String() method to it.
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// LetStatement is a node that identifies a monkey let statement
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier // hold the name of the identifier
	Value Expression  // hold the expression
}

// Assert implementations
var _ Node = (*LetStatement)(nil)
var _ Statement = (*LetStatement)(nil)

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

// Identifier is a node that identifies a monkey identifier
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

var _ Node = (*Identifier)(nil)
var _ Expression = (*Identifier)(nil)

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// ReturnStatement is a node that represents a monkey return statement
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

var _ Node = (*ReturnStatement)(nil)
var _ Statement = (*ReturnStatement)(nil)

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

// ExpressionStatement is a node that represents a monkey expression statement
// It acts as a wrapper for an actual expression
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

var _ Node = (*ExpressionStatement)(nil)
var _ Statement = (*ExpressionStatement)(nil)

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral is a node representing an expression with integer literal
// example: 5;
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

var _ Node = (*IntegerLiteral)(nil)
var _ Expression = (*IntegerLiteral)(nil)

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
