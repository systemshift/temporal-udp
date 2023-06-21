package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	CONN_HOST       = "localhost"
	FILE_PORT       = "8080"
	LATENCY_PORT    = "8081"
	FILE_ADDR       = CONN_HOST + ":" + FILE_PORT
	LATENCY_ADDR    = CONN_HOST + ":" + LATENCY_PORT
	MAX_PACKET_SIZE = 1400
	HEADER_SIZE     = 8
	K               = 5
)

func main() {
	fileTransferConn, err := net.Dial("udp", FILE_ADDR)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer fileTransferConn.Close()

	latencyConn, err := net.Dial("udp", LATENCY_ADDR)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer latencyConn.Close()

	go measureLatency(latencyConn)

	// send file transfer request
	file, err := os.Open("client_storage/01 - Angel Attack.mkv")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer file.Close()

	buffer := make([]byte, MAX_PACKET_SIZE)
	sequenceNumber := uint64(0)
	for {
		n, err := file.Read(buffer[HEADER_SIZE:])
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error: ", err)
			}
			break
		}

		data := append(intToBytes(sequenceNumber), buffer[HEADER_SIZE:n+HEADER_SIZE]...)
		fileTransferConn.Write(data)

		sequenceNumber++
	}

}

func measureLatency(conn net.Conn) {
	latencies := make([]time.Duration, 0, K)
	buffer := make([]byte, 1024)

	for i := 0; i < K; i++ {
		//send lmb message
		_, err := fmt.Fprintf(conn, "lmb:%d", i)
		if err != nil {
			fmt.Println("Error: ", err)
			i-- // retry
			continue
		}

		startTime := time.Now()

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error: ", err)
			i-- // retry
			continue
		}

		message := string(buffer[:n])
		// check if message is a reply for latency
		if message == fmt.Sprintf("lme:%d", i) {
			latency := time.Since(startTime)
			latencies = append(latencies, latency)
			fmt.Printf("Latency %d: %v\n", i, latency)
		}
	}

	// compute average latency
	total := time.Duration(0)
	for _, latency := range latencies {
		total += latency
	}
	avgLatency := total / time.Duration(len(latencies))

	fmt.Printf("Average latency: %v\n", avgLatency)
}

func intToBytes(n uint64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes()
}
