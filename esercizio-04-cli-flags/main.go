package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "filetools",
	Short: "A versatile file processing tool",
}

var countCmd = &cobra.Command{
	Use:   "count [files...]",
	Short: "Count lines, words, and characters",
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagVerbose && flagQuiet {
			return fmt.Errorf("cannot use --verbose and --quiet together")
		}
		if len(args) == 0 {
			return fmt.Errorf("no files provided")
		}
		f := strings.ToLower(flagFormat)
		if f != "text" && f != "json" && f != "csv" {
			return fmt.Errorf("invalid format: %s", flagFormat)
		}
		flagFormat = f
		type FileStats struct {
			File  string
			Stats Stats
		}
		results := []FileStats{}
		for _, path := range args {
			stats, err := countFile(path, flagLines)
			if err != nil {
				return err
			}
			results = append(results, FileStats{File: path, Stats: stats})
		}

		switch flagFormat {
		case "text":
			for _, r := range results {
				fmt.Printf("%s: lines=%d words=%d chars=%d\n", r.File, r.Stats.Lines, r.Stats.Words, r.Stats.Chars)
			}
		case "json":
			json.NewEncoder(os.Stdout).Encode(results)
		case "csv":
			fmt.Println("file,lines,words,chars")
			for _, r := range results {
				fmt.Printf("%s,%d,%d,%d\n", r.File, r.Stats.Lines, r.Stats.Words, r.Stats.Chars)
			}

		}

		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [files...]",
	Short: "Search for a pattern in files",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO
		return nil
	},
}

var (
	flagLines   int
	flagFormat  string
	flagVerbose bool
	flagQuiet   bool
	flagPattern string
)

type Stats struct{ Lines, Words, Chars int }

func init() {
	rootCmd.AddCommand(countCmd)
	countCmd.Flags().IntVar(&flagLines, "lines", 0, "number of lines to process")
	countCmd.Flags().StringVar(&flagFormat, "format", "text", "output format")
	countCmd.Flags().BoolVar(&flagVerbose, "verbose", false, "verbose output")
	countCmd.Flags().BoolVar(&flagQuiet, "quiet", false, "quiet output")
	searchCmd.Flags().StringVar(&flagPattern, "pattern", "", "pattern to search")
	searchCmd.MarkFlagRequired("pattern")
	rootCmd.AddCommand(searchCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func countFile(path string, maxLines int) (Stats, error) {
	f, err := os.Open(path)
	if err != nil {
		return Stats{}, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	stats := Stats{}
	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++
		stats.Words += len(strings.Fields(line))
		stats.Chars += len(line)
		if maxLines > 0 && stats.Lines >= maxLines {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return Stats{}, err
	}

	return stats, nil

}
