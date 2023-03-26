package client

import (
	"fmt"
	"log"
	"net"
	"time"
)

func Connect(address string, message []string) {
	source, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, source)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("udp server is %s\n", conn.RemoteAddr().String())

	ticker := time.NewTicker(time.Millisecond * 10)
	i := 0
	for range ticker.C {
		_, err := conn.Write([]byte(message[i]))
		if err != nil {
			log.Default().Println("error writing to udp server", err)
		} else {
			log.Default().Println("message sent to udp server", err)
		}
		fmt.Println(i)
		i += 1
		time.Sleep(time.Millisecond * 80)

	}
}
