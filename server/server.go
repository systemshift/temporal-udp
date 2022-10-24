package server

import (
	"fmt"
	"log"
	"net"
)

var opentimewindow = 100
var closetimewindow = 900
var listenAddr string

func Listen(listenAddr string) {
	// listen for incoming udp packets
	source, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", source)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("udp server is", conn.LocalAddr().String())
	buffer := make([]byte, 1024)

	message, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("message from client is", string(buffer[:message]))
	conn.Close()

}
