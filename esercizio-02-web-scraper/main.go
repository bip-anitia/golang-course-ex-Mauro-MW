package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type PageInfo struct {
	URL         string
	Title       string
	StatusCode  int
	ContentSize int
	LinkCount   int
	Error       error
}

func main() {
	workers := flag.Int("workers", 5, "numero massimo di workers")
	timeout := flag.Duration("timeout", 10*time.Second, "timeout per richiesta HTTP")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Uso: go run main.go [-workers=N] [-timeout=10s] <urls.txt | url1 url2 ...>")
		return
	}

	urls, err := readURLs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore lettura URL: %v\n", err)
		return
	}
	if len(urls) == 0 {
		fmt.Println("Nessun URL valido")
		return
	}
	if *workers < 1 {
		*workers = 1
	}

	start := time.Now()
	fmt.Printf("Scraping %d URLs con %d workers...\n\n", len(urls), *workers)

	client := &http.Client{Timeout: *timeout}
	jobs := make(chan string)
	results := make(chan PageInfo)

	var wg sync.WaitGroup
	wg.Add(*workers)
	for i := 0; i < *workers; i++ {
		go func() {
			defer wg.Done()
			for url := range jobs {
				results <- fetch(url, client)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for _, u := range urls {
			jobs <- u
		}
		close(jobs)
	}()

	successes := 0
	for res := range results {
		if res.Error != nil {
			fmt.Printf("[ERROR] %s\n     Error: %v\n\n", res.URL, res.Error)
			continue
		}
		successes++
		fmt.Printf("[OK] %s\n     Status: %d | Size: %d bytes | Links: %d | Title: %q\n\n",
			res.URL, res.StatusCode, res.ContentSize, res.LinkCount, res.Title)
	}

	fmt.Printf("Completato in %s\nSuccessi: %d/%d\n", time.Since(start), successes, len(urls))
}

func readURLs(args []string) ([]string, error) {
	if len(args) == 1 {
		if info, err := os.Stat(args[0]); err == nil && !info.IsDir() {
			return readURLsFromFile(args[0])
		}
	}

	urls := make([]string, 0, len(args))
	for _, a := range args {
		u := strings.TrimSpace(a)
		if u != "" {
			urls = append(urls, u)
		}
	}
	return urls, nil
}

func readURLsFromFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	urls := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func extractTitleAndLinks(r io.Reader) (string, int) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", 0
	}

	var title string
	linkCount := 0

	var visit func(*html.Node)
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
		}
		if n.Type == html.ElementNode && n.Data == "a" {
			linkCount++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}
	visit(doc)

	return title, linkCount
}

func fetch(url string, client *http.Client) PageInfo {
	page := PageInfo{
		URL: url,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		page.Error = err
		return page
	}

	req.Header.Set("User-Agent", "go-scraper/1.0")

	resp, err := client.Do(req)

	if err != nil {
		page.Error = err
		return page
	}
	defer resp.Body.Close()
	page.StatusCode = resp.StatusCode
	if resp.StatusCode >= 400 {
		page.Error = fmt.Errorf("bad status: %d", resp.StatusCode)
		return page
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		page.Error = err
		return page
	}
	page.ContentSize = len(data)
	page.Title, page.LinkCount = extractTitleAndLinks(bytes.NewReader(data))
	return page

}
