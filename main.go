package main

import (
	"dot/lexer"
	"dot/parser"
	"fmt"
	"os"
)

func main() {
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
	fmt.Println(program.String())
}
