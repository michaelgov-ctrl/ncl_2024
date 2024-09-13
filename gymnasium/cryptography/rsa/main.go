package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// https://en.wikipedia.org/wiki/RSA_(cryptosystem)#Operation
// https://trove.cyberskyline.com/fafad0e3926b4f808e7616039ca6a6af

// go build -ldflags '-s'
func main() {
	var n, e int
	flag.IntVar(&n, "n", -1, "REQUIRED: n value")
	flag.IntVar(&e, "e", -1, "REQUIRED: e value")

	flag.Parse()
	if n == -1 || e == -1 {
		flag.Usage()
		os.Exit(1)
	}

	ciphertext := flag.Arg(0)
	if ciphertext == "" {
		log.Fatal("missing positional argument for ciphertext")
	}

	plaintext, err := nclDecryptRSA(n, e, ciphertext)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(plaintext)
}
