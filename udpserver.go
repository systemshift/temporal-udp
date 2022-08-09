package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	open_buffer  = 1024
	close_buffer = 0
)

func main() {
	source, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	conn, err := net.ListenUDP("udp", source)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("udp server is %s\n", conn.LocalAddr().String())

	defer conn.Close()

	buffer := make([]byte, 1024)
	close_buffer := make([]byte, 0)

	open_ticker := time.NewTicker(time.Millisecond * 500)
	close_ticker := time.NewTicker(time.Millisecond * 1000)

	// for {
	// 	_, err := conn.Read(buffer)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("message from client is %s\n", string(buffer))
	// 	time.Sleep(time.Second)

	// }

	for {
		select {
		case <-open_ticker.C:
			_, err := conn.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("message from client is %s\n", string(buffer))
			time.Sleep(time.Second)
		case <-close_ticker.C:
			_, err := conn.Read(close_buffer)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("message from client is %s\n", string(close_buffer))
			time.Sleep(time.Second)
		}

	}
}
