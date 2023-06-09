package main

import (
	"fmt"
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
	serverAddr, err := net.ResolveUDPAddr("udp", HANDSHAKE_ADDR)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// listen for incoming packets
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	buffer := make([]byte, MAX_PACKET_SIZE)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// trim buffer to actual packet size
	data := buffer[:n]
	if len(data) != 8 {
		fmt.Println("packet is not 8 bytes, drop it", err)
		os.Exit(1)
	} else {
		// print prng and reply with ack to client
		fmt.Println("prng seed: ", data)
		_, err = conn.WriteToUDP([]byte("ack"), addr)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}
}
