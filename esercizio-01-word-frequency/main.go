package main

import (
	"bufio"
	"flag"
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
	top := flag.Int("top", 0, "numero di parole da mostrare (0 = tutte)")
	ignoreCase := flag.Bool("ignore-case", true, "ignora maiuscole/minuscole")
	flag.Parse()
	files := flag.Args()

	// Leggi da file se forniti, altrimenti da stdin.
	if len(files) > 0 {
		for _, filename := range files {
			f, err := os.Open(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "errore apertura file:", err)
				continue
			}
			if err := countLines(f, counts, *ignoreCase); err != nil {
				fmt.Fprintln(os.Stderr, "errore lettura file:", err)
			}
			f.Close()
		}
	} else {
		if err := countLines(os.Stdin, counts, *ignoreCase); err != nil {
			fmt.Fprintln(os.Stderr, "errore lettura stdin:", err)
		}
	}

	// Converti la mappa in slice per poter ordinare per frequenza.
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
	// Calcola statistiche globali.
	totalWords := 0
	for _, item := range items {
		totalWords += item.Count
	}

	uniqueWords := len(counts)
	fmt.Printf("Parole totali: %d\n", totalWords)
	fmt.Printf("Parole uniche: %d\n\n", uniqueWords)

	if *top > 0 {
		fmt.Printf("Top %d parole più frequenti:\n", *top)
	} else {
		fmt.Println("Tutte le parole (ordinate per frequenza):")
	}

	// Limita la stampa se è stato richiesto un top N.
	limit := len(items)
	if *top > 0 && *top < limit {
		limit = *top
	}

	for i := 0; i < limit; i++ {
		if items[i].Count <= 1 {
			continue
		}
		fmt.Printf("%d. %q - %d occorrenze\n", i+1, items[i].Word, items[i].Count)
	}

}

func countLines(f *os.File, counts map[string]int, ignoreCase bool) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if ignoreCase {
			line = strings.ToLower(line)
		}
		// Spezza la riga in parole ignorando punteggiatura e spazi.
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
