// Package repl provides facilities to read, evaluat, print and loop
// monkey code.
// For now it limits itself to print out the tokens of input code.
package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> "

// Start takes an input and output, and initiates the main REPL loop.
func Start(in io.Reader, out io.Writer) {
	// Start a new scanner
	scanner := bufio.NewScanner(in)
	for {
		// Print the prompt
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		// Takes a string of bytes
		line := scanner.Text()
		// Start a new lexer with said string
		l := lexer.New(line)
		// Read each token until the end of the file
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			// Print each token
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
