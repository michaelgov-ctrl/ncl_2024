package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	f, err := os.Open("green_file")
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	var files []string
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == '\x00' {
			fmt.Println(i)
		}

		filetype := http.DetectContentType(bytes[i:])
		if filetype != "application/octet-stream" && filetype != "text/plain; charset=utf-8" {
			fmt.Println(i, filetype)
			files = append(files, filetype)
		}
	}

	fmt.Println(files, len(files))
}
