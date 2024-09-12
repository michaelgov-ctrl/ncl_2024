package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

var hashes appMap

func main() {
	hashes.new()
	hashes.store("68a96446a5afb4ab69a2d15091771e39", "")
	hashes.store("ec5f0b1826389df8622133014e88afde", "")
	hashes.store("32e5f63b189b78dccf0b97ac41f0d228", "")
	hashes.store("2233287f476ba63323e60addca1f6b64", "")
	hashes.store("6539bbb84fe2de2628fc5e4f2a31f23a", "")

	file, err := os.Open("rockyou.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

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
		go worker(&wg, buffer)

		buffer = make([]string, workerSliceLen)
		buffer[0] = scanner.Text()
		i = 0
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	fmt.Println(hashes.m)
}
