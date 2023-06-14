package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
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
	var seed int64
	mod := int64(10000)
	client_addr := &net.UDPAddr{}

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
			seed_string := strs[0]
			start_time := strs[1]

			// convert seed to int64
			seed, err = strconv.ParseInt(seed_string, 10, 64)
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}

			// print seed and start time
			fmt.Println("Seed: ", seed_string)
			fmt.Println("Start time: ", start_time)

			// send ack
			handshake_conn.WriteToUDP([]byte("ack"), client_addr)

			time.Sleep(3 * time.Second)

			// break out of loop, make sure to fix this later, we are trying to get the random seed working
			break
		} else {
			// print raw data
			fmt.Println("Packet does not = 2. Received: ", data)
		}
	}

	// seed random number generator
	rand.Seed(seed)

	// starting the agreed upon time, generate random numbers and send them to the client
	interval := rand.Int63n(mod) * time.Hour.Milliseconds()

	ticker := time.NewTicker(time.Duration(interval))

	// setup file in memory and prepare to send, will use later to send packets over the time loop
	file_info, err := os.Stat("server_storage/01 - Angel Attack.mkv")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	file_array := make([]byte, file_info.Size())
	// read file into memory
	file, err := os.Open("server_storage/01 - Angel Attack.mkv")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	n, err := file.Read(file_array)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Read ", n, " bytes from file")

	// setup connection back to client
	client_file_addr := &net.UDPAddr{IP: client_addr.IP, Port: 8080}
	conn, err := net.DialUDP("udp", nil, client_file_addr)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	i := 0

	for range ticker.C {
		// print current time
		fmt.Println(time.Now())

		// update next_time using the seeded RNG to generate the next time
		interval = rand.Int63n(mod) * time.Hour.Milliseconds()
		ticker.Reset(time.Duration(interval))

		// send packet to client
		_, err = conn.Write(file_array[i : i+MAX_PACKET_SIZE])
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		// todo: check for ack from client
		// increment i
		i += MAX_PACKET_SIZE
		// check if EOF
		if i >= len(file_array) {
			fmt.Println("EOF")
			break
		}

	}

}
