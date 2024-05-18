package main

import (
	"event-streaming-system/pkg/client"
	"log"
)

func main() {
	c, err := client.Connect("localhost:8080")
	if err != nil {
		log.Println("Error connecting to server:", err)
		return
	}
	defer c.Close()

	log.Println("Client connected successfully")

	err = c.Publish("stream1", "Hello World!")
	if err != nil {
		log.Println("Error Publish:", err)
		return
	}

	err = c.Publish("stream1", "Another message")
	if err != nil {
		log.Println("Error Publish:", err)
		return
	}

	err = c.Publish("stream2", "This one goes to stream2")
	if err != nil {
		log.Println("Error Publish:", err)
		return
	}
}
