package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type WordCount struct {
	Word  string
	Count int
}

func main() {
	counts := make(map[string]int)

	if len(os.Args) > 1 {
		for _, filename := range os.Args[1:] {
			f, err := os.Open(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "errore apertura file:", err)
				continue
			}
			if err := countLines(f, counts); err != nil {
				fmt.Fprintln(os.Stderr, "errore lettura file:", err)
			}
			f.Close()
		}
	} else {
		if err := countLines(os.Stdin, counts); err != nil {
			fmt.Fprintln(os.Stderr, "errore lettura stdin:", err)
		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		words := strings.FieldsFunc(line, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
		for _, w := range words {
			if w != "" {
				counts[w]++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
