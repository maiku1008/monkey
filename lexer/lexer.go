// Package lexer provides facilities to convert monkey source code into tokens
package lexer

import (
	"monkey/token"
)

// Lexer translates source code into tokens
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// New initialises a Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads each character and updates the Lexer's fields.
// It does so by advancing the current position one at a time at each call
// until the end of the input.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// peekChar returns the next byte character to the current one
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhiteSpace calls readChar() on the lexer if the current character
// is a whitespace of some kind
func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// NextToken returns a new Token depending on the current character
func (l *Lexer) NextToken() token.Token {
	l.skipWhiteSpace()
	var tok token.Token
	switch l.ch {
	case '=':
		// Check if this is an EQ operator "=="
		if l.peekChar() == '=' {
			ch := l.ch                           // save the current character
			l.readChar()                         // move forward, updating l.ch
			literal := string(ch) + string(l.ch) // compose the literal with previous and current ch
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '!':
		// Check if this is an NOT_EQ operator "!="
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// In this case, the character is a letter.
		// We read the whole word and determine if it's a valid identifier or not.
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}
		// In this case, the character is a number.
		// We simply read the whole number and return it as the literal.
		// Monkey supports integers only, for simplicity.
		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		}
		// If we end up here, we don't know how to handle this character
		// and we mark it as illegal
		tok = newToken(token.ILLEGAL, l.ch)
	}
	// advance
	l.readChar()
	return tok
}

// newToken initialises a Token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// isLetter identifies whether a byte represents a letter or not
// An underscore is considered a valid letter,
// so we can enable identifier such as some_number
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit identifies whether a character represents a number or not
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// readIdentifier continues reading the string from the current position
// until the byte is not a letter, and returns the resulting string.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber continues reading the string from the current position
// until the byte is not a digit anymore, and returns the resulting string
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}
