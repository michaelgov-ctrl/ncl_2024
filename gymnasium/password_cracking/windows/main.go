package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	filename := "titles.txt"
	urls := []string{
		"https://en.wikipedia.org/wiki/List_of_Law_%26_Order:_Special_Victims_Unit_episodes_(season_20%E2%80%93present)",
		"https://en.wikipedia.org/wiki/List_of_Law_%26_Order:_Special_Victims_Unit_episodes_(seasons_1%E2%80%9319)",
	}

	// output base titles to file for later reuse
	titles, err := scrapeToFile(filename, urls)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	wordlist := createWordlist(titles)

	wl, err := os.Create("./wordlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer wl.Close()

	w := bufio.NewWriter(wl)
	for _, word := range wordlist {
		_, err := w.WriteString(fmt.Sprintf("%s\n", word))
		if err != nil {
			log.Fatal(err)
		}
	}

	w.Flush()
}

func createWordlist(titles []string) []string {
	var wordlist []string
	for _, title := range titles {
		replaced := strings.ReplaceAll(title, " ", "")
		temp := []string{
			title,
			replaced,
			strings.ToLower(title),
			strings.ToLower(replaced),
		}

		var expanded []string
		for _, t := range temp {
			for i := 0; i < 100; i++ {
				expanded = append(expanded, fmt.Sprintf("%s%02d", t, i))
			}
		}

		wordlist = append(wordlist, expanded...)
	}

	return wordlist
}
