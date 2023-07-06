package client

import (
	"fmt"
	"math/rand"
	"net"
)

const (
	HANDSHAKE_PORT  = 8080
	FILE_PORT       = 8081
	MAX_PACKET_SIZE = 1400
)

func Handshake(ip string) {
	// establish connection
	handshake_conn, err := net.Dial("udp", ip+":8080")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer handshake_conn.Close()

	// create seed
	seed := rand.Int63()
	seed_str := fmt.Sprintf("%d", seed)

	// print seed
	fmt.Println("seed: ", seed)

	// send seed
	handshake_conn.Write([]byte(seed_str))

	// listen for income reply
	buf := make([]byte, MAX_PACKET_SIZE)

	// read from connection into buffer
	_, err = handshake_conn.Read(buf)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// print received message
	fmt.Printf("received: %s\n", string(buf))
}
