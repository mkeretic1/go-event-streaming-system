package main

import (
	"event-streaming-system/pkg/server"
	"fmt"
	"log"
)

const (
	addr string = "localhost"
	port string = "8080"
)

func main() {
	if err := server.Start(fmt.Sprintf("%s:%s", addr, port)); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
