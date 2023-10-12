package client

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
)

const (
	HANDSHAKE_PORT  = 8080
	FILE_PORT       = 8081
	MAX_PACKET_SIZE = 1400
)

func Handshake(ip string) (string, int, int) {
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

	// listen for incoming reply
	buf := make([]byte, MAX_PACKET_SIZE)

	// read from connection into buffer
	_, err = handshake_conn.Read(buf)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// print received message
	fmt.Printf("received: %s\n", string(buf))

	// parse message into filename, size, and number of pieces
	str_arr := strings.Split(string(buf), "\\")
	filename := strings.Split(str_arr[0], ":")[1]
	filesize := strings.Split(str_arr[1], ":")[1]
	pieces := strings.Split(str_arr[2], ":")[1]

	// convert filesize and pieces to int
	filesize_int := 0
	pieces_int := 0
	fmt.Sscanf(filesize, "%d", &filesize_int)
	fmt.Sscanf(pieces, "%d", &pieces_int)

	return filename, filesize_int, pieces_int
}
