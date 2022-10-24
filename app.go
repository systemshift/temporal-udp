package main

import (
	"temporal-udp/server"
)

func main() {
	//server.listenAddr("127.0.0.1:8080")
	server.Listen("127.0.0.1:8080")
}
