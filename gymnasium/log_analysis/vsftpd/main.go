package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	EventTime time.Time
	PID       int
	User      string
	Action    string
	Status    string
	Client    string
	File      string
	Bytes     int
	Rate      string
}

var re = regexp.MustCompile(`\[[^\]]*\]`)

func main() {
	f, err := os.Open("vsftpd.log")
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var entries []Entry
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) > 36 && t[36] == '[' {
			entries = append(entries, NewEntryFromString(t))
		}
	}

	// What file extension was the most used by ftpuser?
	var extensions, max, common = make(map[string]int), 0, ""

	// How many total bytes did this other user upload?
	// How many total bytes did ftpuser upload?
	var ftpuserUpBytes, jimmyUpBytes int

	// How many total bytes did ftpuser download?
	var ftpuserDownBytes int
	for _, e := range entries {
		if e.User == "ftpuser" {
			if e.Action == "UPLOAD" {
				ftpuserUpBytes += e.Bytes
			} else if e.Action == "DOWNLOAD" {
				ftpuserDownBytes += e.Bytes
			}

		}

		if e.User == "jimmy" {
			if e.Action == "UPLOAD" {
				jimmyUpBytes += e.Bytes
			}
		}

		idx := strings.LastIndex(e.File, ".")
		if idx == -1 {
			continue
		}

		ext := e.File[idx:len(e.File)]
		extensions[ext]++
		if max < extensions[ext] {
			max = extensions[ext]
			common = ext
		}
	}

	fmt.Println("most frequent extension: ", common)
	fmt.Println("ftpuserUpBytes:  ", ftpuserUpBytes)
	fmt.Println("ftpuserDownBytes: ", ftpuserDownBytes)
	fmt.Println("jimmyUpBytes: ", jimmyUpBytes)
}

func NewEntryFromString(s string) Entry {
	e := Entry{}
	dateTime, _ := time.Parse("Mon Jan 02 15:04:05 2006", string(s[:24]))
	e.EventTime = dateTime

	bracketedTokens := re.FindAllString(s, -1)
	if len(bracketedTokens) > 0 {
		e.PID, _ = strconv.Atoi(strings.Trim(bracketedTokens[0], "[pid ]"))
	}

	if len(bracketedTokens) > 1 {
		e.User = strings.Trim(bracketedTokens[1], "[]")
	}

	split := strings.Split(s, "]")
	backEnd := split[2]

	fields := strings.Fields(backEnd)
	if len(fields) > 0 {
		e.Status = strings.Trim(fields[0], " ")
	}

	if len(fields) > 1 {
		e.Action = strings.Trim(fields[1], ": ")
	}

	if len(fields) > 3 {
		e.Client = strings.Trim(fields[3], "\",")
	}

	commaDelim := strings.Split(backEnd, ",")

	if len(commaDelim) > 1 {
		e.File = strings.Trim(commaDelim[1], "\",")
	}

	if len(commaDelim) > 2 {
		e.Bytes, _ = strconv.Atoi(strings.Trim(commaDelim[2], " bytes"))
	}

	if len(commaDelim) > 3 {
		e.Rate = commaDelim[3]
	}

	return e
}
