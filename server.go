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

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer conn.Close()

	file, err := os.Create("server_storage/01 - Angel Attack.mkv")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer file.Close()

	buffer := make([]byte, MAX_PACKET_SIZE)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		_, err = file.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}
}
