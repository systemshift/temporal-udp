package server

import (
	"fmt"
	"net"
	"os"
)

const (
	HANDSHAKE_PORT  = 8080
	FILE_PORT       = 8081
	MAX_PACKET_SIZE = 1400
)

func HandshakeListen(ip string, file string) {
	// listen for income connection
	handshake_conn, err := net.ListenPacket("udp", ip+":8080")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer handshake_conn.Close()

	// listen for incoming message
	buf := make([]byte, MAX_PACKET_SIZE)
	_, addr, err := handshake_conn.ReadFrom(buf)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// print received seed from addr
	fmt.Printf("received from %s: %s\n", addr, string(buf))

	// prepare file metadate to send back to client
	file_metadata, err := os.Stat(file)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	file_size := file_metadata.Size()
	file_size_str := fmt.Sprintf("%d", file_size)
	// calculate the number of pieces to send
	num_pieces := file_size / 1200 // 1200 for data, 200 for header

	// message to send back to client with file name, size, and number of pieces
	message := "filename:" + file + "\\ filesize:" + file_size_str + "\\ numpieces:" + fmt.Sprintf("%d", num_pieces)

	// send file name and sizz back to client
	_, err = handshake_conn.WriteTo([]byte(message), addr)
	if err != nil {
		fmt.Println("Error: ", err)
	}

}
