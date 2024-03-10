package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("Usage: %s [script]\n", os.Args[0])
		return
	}

	lox := NewLox()

	if len(os.Args) == 2 {
		if err := lox.RunFile(os.Args[1]); err != nil {
			fmt.Printf("could not execute file %s: %+v", os.Args[1], err)
			return
		}
		return
	}

	if err := lox.RunPrompt(); err != nil {
		fmt.Printf("could not execute input: %+v", err)
		return
	}

	return
}
