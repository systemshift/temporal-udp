package client

import (
	"fmt"
	"log"
	"net"
)

var opentimewindow = 100
var closetimewindow = 900
var sendAddr string
var userInput string

func Connect(sendAddr, userInput string) {
	source, err := net.ResolveUDPAddr("udp", sendAddr)
	conn, err := net.DialUDP("udp", nil, source)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("udp server is %s\n", conn.RemoteAddr().String())

	message := userInput
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Default().Println("error writing to udp server")
	}
	conn.Close()
}
