package cmd

import "fmt"

func loxError(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	// FUTURE: Add more useful error context
	//
	// Error: Unexpected "," in argument list.

	// 15 | function(first, second,);
	//                            ^-- Here.
	fmt.Printf("[line %d] Error %s: %s", line, where, message)
}
