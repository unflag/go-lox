package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("Usage: %s [script]\n", os.Args[0])
		return
	}

	if len(os.Args) == 2 {
		if err := runFile(os.Args[1]); err != nil {
			fmt.Printf("could not execute file %s: %+v", os.Args[1], err)
			return
		}
		return
	}

	if err := runPrompt(); err != nil {
		fmt.Printf("could not execute input: %+v", err)
		return
	}

	return
}

func runFile(path string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	return run(string(src))
}

func runPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("could not read input: %+v", err)
			}
			break
		}

		if err := run(scanner.Text()); err != nil {
			return fmt.Errorf("could not run input: %+v", err)
		}
	}

	return nil
}

func run(src string) error {
	if src != "" {
		s := newScanner(src)
		tokens, err := s.Scan()
		if err != nil {
			return fmt.Errorf("could not parse input: %w", err)
		}

		p := newParser(tokens)
		expr := p.Parse()
		fmt.Println("Parser: ", newPrinter().Print(expr))

		result := newInterpreter().Evaluate(expr)
		fmt.Println("Interpreter: ", result)
	}

	return nil
}
