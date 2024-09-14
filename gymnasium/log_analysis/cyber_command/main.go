package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ghostiam/binstruct"
)

type Header struct {
	Magic        []byte `bin:"len:8"`
	Version      uint8  `bin:"len:1"`
	Timestamp    uint32 `bin:"len:4"`
	HostnameLen  uint32 `bin:"len:4"`
	Hostname     []byte `bin:"len:HostnameLen"`
	FlagLen      uint32 `bin:"len:4"`
	Flag         []byte `bin:"len:FlagLen"`
	NumOfEntries uint32 `bin:"len:4"`
}

type BodyItem struct {
	SourceIp         uint32 `bin:"len:4"`
	DestIp           uint32 `bin:"len:4"`
	Timestamp        uint32 `bin:"len:4"`
	BytesTransferred uint32 `bin:"len:4"`
}

type SkyPacket struct {
	Header
	Body []BodyItem `bin:"len:NumOfEntries"`
}

func main() {
	f, err := os.Open("custom_file_format.sky")
	if err != nil {
		//log.Fatalf("%s\n", err)
	}

	var sp SkyPacket
	decoder := binstruct.NewDecoder(f, binary.BigEndian)

	err = decoder.Decode(&sp)
	if err != nil {
		//log.Fatal(err)
	}

	logs := parseLogs(sp.Body)

	var totalTransferredBytes int
	uniqueIps := make(map[string]bool)
	ipSentData := make(map[string]int)
	dayCounter := make(map[string]int)
	for _, l := range logs {
		totalTransferredBytes += l.BytesTransferred
		uniqueIps[l.DestIp.String()] = true
		uniqueIps[l.SourceIp.String()] = true

		ipSentData[l.SourceIp.String()] += l.BytesTransferred

		y, m, d := l.Timestamp.Date()
		dayCounter[fmt.Sprintf("%d-%d-%d", y, m, d)] += l.BytesTransferred
	}

	mostSent, bigSender := 0, ""
	for k, v := range ipSentData {
		if mostSent < v {
			bigSender, mostSent = k, v
		}
	}

	fmt.Println(dayCounter)
	busiestDay, bigDay := 0, ""
	for k, v := range dayCounter {
		if busiestDay < v {
			bigDay, busiestDay = k, v
		}
	}

	flag, _ := base64.StdEncoding.DecodeString(string(sp.Flag))

	//fmt.Printf("%+v\n", sp)
	fmt.Println("file created on: ", time.Unix(int64(sp.Timestamp), 0))
	fmt.Println("hostname: ", string(sp.Hostname))
	fmt.Println("flag: ", string(flag))
	fmt.Println("totalTransferredBytes: ", totalTransferredBytes)
	fmt.Println("entries in the log file: ", len(logs))
	fmt.Println("num of unique ips: ", len(uniqueIps))
	fmt.Println("biggest sender: ", bigSender, "sent: ", mostSent)
	fmt.Println("busiest day: ", bigDay)
}

type log struct {
	SourceIp         net.IP
	DestIp           net.IP
	Timestamp        time.Time
	BytesTransferred int
}

func parseLogs(items []BodyItem) []log {
	logs := make([]log, len(items))
	for idx, item := range items {
		sourceip := make(net.IP, 4)
		binary.BigEndian.PutUint32(sourceip, item.SourceIp)

		destip := make(net.IP, 4)
		binary.BigEndian.PutUint32(destip, item.SourceIp)

		fmt.Println()
		logs[idx] = log{
			SourceIp:         sourceip,
			DestIp:           destip,
			BytesTransferred: int(item.BytesTransferred),
			Timestamp:        time.Unix(int64(item.Timestamp), 0),
		}
	}

	return logs
}
