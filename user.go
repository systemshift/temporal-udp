package main

import (
	"temporal-udp/client"
)

func main() {
	client.Connect("127.0.0.1:8080", "hello")
}
