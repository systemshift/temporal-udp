package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
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
	mod := int64(5)
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
			handshake_conn.Close()
			break
		}

	}
	// server now has seed and start time

	// generate random sequence from shared seed
	rand.Seed(seed)

	// starting at the agreed start time, use the seeded RNG to generate the sequence of printing times. Server should produce same sequence
	interval := rand.Int63n(mod) * time.Hour.Milliseconds()

	ticker := time.NewTicker(time.Duration(interval))

	// listen back for incoming packets
	file_conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8080})
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	file_conn_file, err := file_conn.File()
	// Set the receive buffer size
	err = syscall.SetsockoptInt(int(file_conn_file.Fd()), syscall.SOL_SOCKET, syscall.SO_RCVBUF, MAX_PACKET_SIZE)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, MAX_PACKET_SIZE)

	// receive metadata packet that includes file name and size
	file_conn.ReadFromUDP(buffer)

	// print metadata
	fmt.Printf("received: %s\n", string(buffer))

	// store file name and size split by :
	full_string := strings.Split(string(buffer), ":")
	file_name := strings.TrimRight(full_string[1], "\x00")

	file_size, err := strconv.Atoi(full_string[0])
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("file name: ", file_name)
	// create file to write to
	file, err := os.Create("client_storage/" + file_name)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer file.Close()

	file_counter := 0
	for range ticker.C {
		// wait until next_time: sleep for random number generated from seeded RNG
		// show the current time
		// update next_time using the seeded RNG to generate the next time
		// repeat until the number of packets has been sent

		// print current time
		fmt.Println(time.Now())

		// set new interval
		minDuration := int64(1) // add 1 milisecond to min to avoid 0
		interval = rand.Int63n(mod)*time.Hour.Milliseconds() + minDuration
		ticker.Reset(time.Duration(interval))

		// read from connection
		fmt.Println("reading from connection into buffer to print file")
		n, addr, err := file_conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		// print size of n
		fmt.Printf("received: %d\n", n)

		// writing buffer to file
		fmt.Println("writing buffer to file")

		// write to file
		file.Write(buffer[:n])
		file.Sync()

		// reply with ack
		_, err = file_conn.WriteToUDP([]byte("ack"), addr)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		// increment file counter
		file_counter++ // change later to increment by n

		// check if all packets have been received
		if file_counter == file_size {
			break
		}

	}

}
