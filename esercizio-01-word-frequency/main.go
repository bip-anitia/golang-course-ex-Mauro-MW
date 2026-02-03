package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
	items := make([]WordCount, 0, len(counts))
	for w, c := range counts {
		items = append(items, WordCount{Word: w, Count: c})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count
		}
		return items[i].Word < items[j].Word
	})
	totalWords := 0
	for _, item := range items {
		totalWords += item.Count
	}
	uniqueWords := len(counts)

	fmt.Printf("Parole totali: %d\n", totalWords)
	fmt.Printf("Parole uniche: %d\n\n", uniqueWords)

	for _, item := range items {
		if item.Count > 1 {
			fmt.Printf("%d\t%s\n", item.Count, item.Word)
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
