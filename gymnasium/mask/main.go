package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	prefix := "SKY-HQNT-"

	f, err := os.Create("./wordlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i := 0; i < 10000; i++ {
		_, err := w.WriteString(fmt.Sprintf("%s%04d\n", prefix, i))
		if err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()
}
