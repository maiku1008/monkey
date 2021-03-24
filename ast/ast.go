// Package ast has the necessary facilities to represent
// an Abstract Syntax Tree for the Monkey language
package ast

import "monkey/token"

// Node represents a single node in the AST.
type Node interface {
	TokenLiteral() string
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

// Identifier is a node that identifies a monkey identifier
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

// Assert implementations
var _ Node = (*Identifier)(nil)
var _ Expression = (*Identifier)(nil)

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
