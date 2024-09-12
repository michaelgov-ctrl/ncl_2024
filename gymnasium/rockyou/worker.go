package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
)

func worker(wg *sync.WaitGroup, strings []string) {
	defer wg.Done()

	for _, str := range strings {
		hash := md5.Sum([]byte(str))

		ok, err := hashes.updateOnceIfExists(hex.EncodeToString(hash[:]), str)
		if err != nil {
			panic(fmt.Sprintf("hash colision on %s", hash))
		}

		if ok {
			hashes.incFound()
			if hashes.allFound() {
				//TODO stop all processing
			}
			fmt.Printf("%s: %s\n", hex.EncodeToString(hash[:]), str)
		}
	}
}
