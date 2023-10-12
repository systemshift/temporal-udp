package main

import (
	"fmt"
	"temporal-udp/client"
)

func main() {
	filename, filesize, pieces := client.Handshake("34.218.138.6")
	fmt.Println("filename: ", filename, " filesize: ", filesize, " pieces: ", pieces)

}
