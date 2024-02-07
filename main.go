package main

import (
	"bufio"
	"dot/eval"
	"dot/lexer"
	"dot/object"
	"dot/parser"
	"fmt"
	"os"
)

func main() {
	if os.Args[1] == "repl" {
		startRepl()
		return
	}
	filename := os.Args[1]
	if filename == "" {
		fmt.Printf("Usage: %s <filename>\n", filename)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	contentStr := string(content)
	lexer := lexer.NewLexer(contentStr)
	parser := parser.NewParser(lexer)
	program := parser.ParseProgram()
	parser.PrintErrors()
	// fmt.Print(program.String())
	env := object.NewEnvironment()
	evaluated := eval.Eval(program, env)
	fmt.Println(evaluated)
}

func startRepl() {
	in := os.Stdin
	out := os.Stdout
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if line == "exit()" {
			return
		}
		lexer := lexer.NewLexer(line)
		parser := parser.NewParser(lexer)
		program := parser.ParseProgram()
		parser.PrintErrors()
		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			fmt.Fprintln(out, evaluated.String())
		}
	}

}
