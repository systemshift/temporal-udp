package main

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST       = "localhost"
	CONN_PORT       = "8080"
	CONN_ADDR       = CONN_HOST + ":" + CONN_PORT
	MAX_PACKET_SIZE = 1400
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", CONN_ADDR)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer conn.Close()

	file, err := os.Open("client_storage/01 - Angel Attack.mkv")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer file.Close()

	buffer := make([]byte, MAX_PACKET_SIZE)

	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
		}

		_, err = conn.Write(buffer[:n])

	}
}
