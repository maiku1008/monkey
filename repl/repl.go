// Package repl provides facilities to read, evaluat, print and loop
// monkey code.
// For now it limits itself to print out the tokens of input code.
package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/parser"
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
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
