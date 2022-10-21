package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	source, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	conn, err := net.DialUDP("udp", nil, source)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("udp server is %s\n", conn.RemoteAddr().String())

	ticker := time.NewTicker(time.Millisecond * 342)
	for range ticker.C {
		message := time.Now().String()
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Default().Println("error writing to udp server")
		}
		time.Sleep(time.Second)

	}
}
