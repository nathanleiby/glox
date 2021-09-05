package cmd

import (
	"fmt"
	"os"

	"github.com/nathanleiby/glox/cmd/parser"
	"github.com/nathanleiby/glox/cmd/scanner"
)

func lox(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		err := runFile(args[0])
		if err != nil {
			os.Exit(65)
		}
	} else {
		err := runPrompt()
		if err != nil {
			os.Exit(65)
		}
	}
}

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return run(string(bytes))
}

func run(source string) error {
	s := scanner.NewScanner(source)
	err := s.ScanTokens()
	if err != nil {
		fmt.Println("scan error:", err)
		return err
	}

	p := parser.NewParser(s.Tokens())
	expr, err := p.Parse()
	if err != nil {
		fmt.Println("parse error:", err)
		return err
	}

	fmt.Println(parser.Parenthesize([]parser.Expr{expr}))

	return nil
}

func runPrompt() error {
	for {
		fmt.Printf("> ")
		var input string
		fmt.Scanln(&input)
		if input == "" {
			break
		}
		_ = run(input)
	}

	return nil
}
