package server

import (
	"event-streaming-system/pkg/server/eventstream"
	"log"
	"net"
)

// Start initializes and runs the server
func Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Printf("Server listening on %s", address)

	eventStream := eventstream.GetInstance()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}
		go eventStream.HandleConnection(conn)
	}
}
