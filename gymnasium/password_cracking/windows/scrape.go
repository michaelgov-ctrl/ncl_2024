package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func scrapeToFile(filename string, urls []string) ([]string, error) {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Scraping: %s\n", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Status: %d\n", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s; Failed with:\n%s\n", r.Request.URL, err)
	})

	var episodeTitles []string
	c.OnHTML("table > tbody .summary", func(h *colly.HTMLElement) {
		episodeTitles = append(episodeTitles, strings.Trim(h.Text, "\""))
	})

	for _, u := range urls {
		c.Visit(u)
	}

	wl, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer wl.Close()

	w := bufio.NewWriter(wl)
	for _, e := range episodeTitles {
		_, err := w.WriteString(fmt.Sprintf("%s\n", e))
		if err != nil {
			return nil, err
		}
	}

	w.Flush()

	return episodeTitles, err
}
