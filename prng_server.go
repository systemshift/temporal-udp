package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// first packet includes the following: random seed, start time, number of packets to send

// last packet includes: control message saying EOF
const (
	CONN_HOST       = "0.0.0.0"
	HANDSHAKE_PORT  = 8080
	FILE_PORT       = 8081
	HANDSHAKE_ADDR  = CONN_HOST + ":8080"
	CONN_ADDR       = CONN_HOST + ":8081"
	MAX_PACKET_SIZE = 1400
)

func main() {
	handshake_conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: HANDSHAKE_PORT})
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer handshake_conn.Close()
	buf := make([]byte, MAX_PACKET_SIZE)

	// establish listening connection
	for {
		// check for packet from client is the right size

		n, client_addr, err := handshake_conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Printf("n = %d, buf[:n] = %s\n", n, string(buf[:n]))
		data := string(buf[:n]) // trim packet to size
		strs := strings.Split(data, ":")

		// check for len(strs)==2
		if len(strs) == 2 {
			fmt.Println("Error: expected '%s' to have two numbers delimited by :", data)
			seed := strs[0]
			start_time := strs[1]

			// print seed and start time
			fmt.Println("Seed: ", seed)
			fmt.Println("Start time: ", start_time)

			// send ack
			handshake_conn.WriteToUDP([]byte("ack"), client_addr)

			time.Sleep(3 * time.Second)
		} else {
			// print raw data
			fmt.Println("Packet does not = 2. Received: ", data)
		}
	}

}
