package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
)

const (
	CONN_HOST       = "localhost"
	HANDSHAKE_PORT  = "8080"
	CONN_PORT       = "8081"
	HANDSHAKE_ADDR  = CONN_HOST + ":" + HANDSHAKE_PORT
	CONN_ADDR       = CONN_HOST + ":" + CONN_PORT
	MAX_PACKET_SIZE = 1400
	HEADER_SIZE     = 8
)

func main() {
	prng := rand.Int63()

	serverAddr, err := net.ResolveUDPAddr("udp", HANDSHAKE_ADDR)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	localAddr, err := net.ResolveUDPAddr("udp", "localhost:0")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	// send prng to server
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, prng)
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// wait for ack from server
	reply_buffer := make([]byte, MAX_PACKET_SIZE)
	n, err := conn.Read(reply_buffer)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fmt.Println(string(reply_buffer[:n]))
}
