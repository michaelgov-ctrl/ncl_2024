package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	var hashesFilename, wordlistFilename string
	flag.StringVar(&hashesFilename, "hashes", "", "REQUIRED: filepath to hashes file")
	flag.StringVar(&wordlistFilename, "wordlist", "", "REQUIRED: filepath to wordlist file")

	flag.Parse()

	if hashesFilename == "" || wordlistFilename == "" {
		flag.Usage()
		os.Exit(1)
	}

	hashlist, err := os.Open(hashesFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer hashlist.Close()

	hashes := loadNewAppMap(hashlist)

	wordlist, err := os.Open(wordlistFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer wordlist.Close()

	scanner := bufio.NewScanner(wordlist)

	var (
		workerSliceLen = 500
		wg             sync.WaitGroup
		buffer         = make([]string, workerSliceLen)
	)

	for i := 0; scanner.Scan(); i++ {
		if i < workerSliceLen {
			buffer[i] = scanner.Text()
			continue
		}

		wg.Add(1)
		go worker(&wg, hashes, buffer)

		buffer = make([]string, workerSliceLen)
		buffer[0] = scanner.Text()
		i = 0
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if 0 < len(buffer) {
		wg.Add(1)
		go worker(&wg, hashes, buffer)
	}

	wg.Wait()

	fmt.Println(hashes.m)
}
