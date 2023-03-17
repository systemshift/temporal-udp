package server

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func Listen(listenAddr string) {
	// listen for udp packets
	ticker := time.NewTicker(time.Millisecond * 500)
	source, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	for range ticker.C {

		conn, err := net.ListenUDP("udp", source)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("udp server is", conn.LocalAddr().String())
		buffer := make([]byte, 1024)

		conn.SetDeadline(time.Now().Add(time.Millisecond * 50))
		message, err := conn.Read(buffer)
		if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
			log.Fatal(err)
		}

		fmt.Println("message from client is", string(buffer[:message]))
		conn.Close()
	}
}
