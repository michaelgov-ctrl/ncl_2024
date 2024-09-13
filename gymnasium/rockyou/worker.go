package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
)

func worker(wg *sync.WaitGroup, hashes *appMap, strings []string) {
	defer wg.Done()

	for _, str := range strings {
		hash := md5.Sum([]byte(str))

		ok, err := hashes.updateOnceIfExists(hex.EncodeToString(hash[:]), str)
		if err != nil {
			panic(fmt.Sprintf("hash colision on %s", hash))
		}

		if ok {
			fmt.Printf("%s: %s\n", hex.EncodeToString(hash[:]), str)

			hashes.incFound()
			if hashes.allFound() {
				fmt.Println(hashes.m)
				os.Exit(0)
			}
		}
	}
}
