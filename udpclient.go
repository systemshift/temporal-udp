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
	//defer conn.Close()

	ticker := time.NewTicker(time.Millisecond * 200)
	for range ticker.C {
		message := time.Now().String()
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)

	}
}
