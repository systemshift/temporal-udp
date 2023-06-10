package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	CONN_HOST       = "34.218.138.6"
	HANDSHAKE_PORT  = 8080
	FILE_PORT       = 8081
	HANDSHAKE_ADDR  = CONN_HOST + ":8080"
	FILE_ADDR       = CONN_HOST + ":8081"
	MAX_PACKET_SIZE = 1400
)

func main() {
	seed := rand.Int63()
	mod := int64(10000)
	offset := rand.Int63n(mod)
	start_time := time.Now()

	// send seed and start time
	handshake_conn, err := net.Dial("udp", "34.218.138.6:8080")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer handshake_conn.Close()

	buf := make([]byte, MAX_PACKET_SIZE)
	for {
		// send seed and start time
		fmt.Println("sending seed and start time")
		handshake_conn.Write([]byte(fmt.Sprintf("%d:%d", seed, start_time.UnixMilli()+offset)))
		//handshake_conn.SetReadDeadline(start_time)
		time.Sleep(10 * time.Millisecond)
		fmt.Println("reading from connection into buffer")
		_, err := handshake_conn.Read(buf)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		fmt.Printf("received: %s\n", string(buf))
		time.Sleep(3 * time.Second)
		// check if ack has arrived

		if string(buf[:3]) == "ack" {
			fmt.Println("ack received")
			break
		}

	}
	// server now has seed and start time

}
