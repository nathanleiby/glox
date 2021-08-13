package cmd

import (
	"fmt"
	"os"
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
	s := NewScanner(source)
	err := s.ScanTokens()
	if err != nil {
		return err
	}

	for _, token := range s.Tokens() {
		fmt.Println(token)
	}

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
