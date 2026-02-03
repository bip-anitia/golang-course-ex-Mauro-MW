package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

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
	// TODO: Implementare il concurrent web scraper
	fmt.Println("Concurrent Web Scraper")
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

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "go-scraper/1.0")

	resp, err := client.Do(req)
	if err != nil {
		page.Error = err
		return page
	}
	defer resp.Body.Close()

	page.StatusCode = resp.StatusCode

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		page.Error = err
		return page
	}
	page.ContentSize = len(data)
	page.Title, page.LinkCount = extractTitleAndLinks(bytes.NewReader(data))
	return page

}
