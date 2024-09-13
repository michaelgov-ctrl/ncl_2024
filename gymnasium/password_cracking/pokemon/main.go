package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("pokemon.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	wl, err := os.Create("./wordlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer wl.Close()

	w := bufio.NewWriter(wl)
	for _, r := range records {
		if r[1] != "" {
			_, err := w.WriteString(fmt.Sprintf("%s\n", strings.ToLower(r[0])))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	w.Flush()
}
