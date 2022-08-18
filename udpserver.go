package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {

	ticker := time.NewTicker(time.Millisecond * 1000)
	source, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	for range ticker.C {

		conn, err := net.ListenUDP("udp", source)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("udp server is %s\n", conn.LocalAddr().String())
		buffer := make([]byte, 1024)

		conn.SetDeadline(time.Now().Add(time.Millisecond * 300))
		message, err := conn.Read(buffer)
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			log.Fatal(err)
		}

		fmt.Printf("message from client is %s\n", string(buffer[:message]))
		conn.Close()
		//time.Sleep(time.Millisecond * 350)

	}

}
