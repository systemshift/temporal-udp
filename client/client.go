package client

import (
	"fmt"
	"log"
	"net"
	"time"
)

func Connect(address, message string) {
	source, err := net.ResolveUDPAddr("udp", address)
	conn, err := net.DialUDP("udp", nil, source)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("udp server is %s\n", conn.RemoteAddr().String())

	ticker := time.NewTicker(time.Millisecond * 10)
	for range ticker.C {
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Default().Println("error writing to udp server")
		}
		time.Sleep(time.Millisecond * 80)

	}
}
