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

	// generate random sequence from shared seed
	rand.Seed(seed)

	// starting at the agreed start time, use the seeded RNG to generate the sequence of printing times. Server should produce same sequence
	interval := rand.Int63n(mod) * time.Hour.Milliseconds()

	ticker := time.NewTicker(time.Duration(interval))
	for range ticker.C {
		// wait until next_time: sleep for random number generated from seeded RNG
		// show the current time
		// update next_time using the seeded RNG to generate the next time
		// repeat until the number of packets has been sent

		// print current time
		fmt.Println(time.Now())

		// set new interval
		interval = rand.Int63n(mod) * time.Hour.Milliseconds()
		ticker.Reset(time.Duration(interval))

	}
	// wait until start time
	//time.Sleep(time.Duration(start_time.UnixMilli()-time.Now().UnixMilli()) * time.Millisecond)

}
