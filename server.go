package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST       = "localhost"
	FILE_PORT       = "8080"
	LATENCY_PORT    = "8081"
	FILE_ADDR       = CONN_HOST + ":" + FILE_PORT
	LATENCY_ADDR    = CONN_HOST + ":" + LATENCY_PORT
	MAX_PACKET_SIZE = 1400
	HEADER_SIZE     = 8
)

func main() {
	go listenForFileTransfer()

	latency_conn, err := net.ListenPacket("udp", LATENCY_ADDR)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer latency_conn.Close()

	buffer := make([]byte, MAX_PACKET_SIZE)
	for {
		n, addr, err := latency_conn.ReadFrom(buffer)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Error: ", err)
			}
			break
		}

		message := string(buffer[:n])

		// check if message is a request for latency
		if len(message) >= 4 && message[:4] == "lmb:" {
			// reply back
			_, err = latency_conn.WriteTo([]byte("lme"), addr)
		}

	}

}

func listenForFileTransfer() {
	latency_conn, err := net.ListenPacket("udp", FILE_ADDR)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer latency_conn.Close()

	receivedPackets := make(map[uint64][]byte)
	buffer := make([]byte, MAX_PACKET_SIZE)

	for {
		n, _, err := latency_conn.ReadFrom(buffer)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Error: ", err)
			}
			break
		}

		// process file packet
		sequenceNumber := bytesToInt(buffer[:HEADER_SIZE])
		data := buffer[HEADER_SIZE:n]
		receivedPackets[sequenceNumber] = data

		if len(receivedPackets) == 1 {
			go writeFile(receivedPackets)

		}
	}
}

func writeFile(receivedPackets map[uint64][]byte) {
	file, err := os.Create("server_storage/01 - Angel Attack.mkv")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer file.Close()

	sequenceNumber := uint64(0)
	for {
		if data, ok := receivedPackets[sequenceNumber]; ok {
			file.Write(data)
			delete(receivedPackets, sequenceNumber)
			sequenceNumber++
		} else {
			break
		}

	}
}

func bytesToInt(b []byte) uint64 {
	buf := bytes.NewBuffer(b)
	var n uint64
	binary.Read(buf, binary.BigEndian, &n)
	return n
}
