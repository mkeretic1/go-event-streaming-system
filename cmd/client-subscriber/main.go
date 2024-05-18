package main

import (
	"event-streaming-system/pkg/client"
	"event-streaming-system/pkg/server/eventstream"
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

	stream1 := c.Subscribe("stream1", eventstream.SubscriptionModeNew)
	stream2 := c.Subscribe("stream2", eventstream.SubscriptionModeAll)
	stream3 := c.Subscribe("stream3", eventstream.SubscriptionModeAfter, 5)

	for {
		select {
		case msg := <-stream1.Chan():
			log.Printf("Message received: %s", msg)
		case msg := <-stream2.Chan():
			log.Printf("Message received: %s", msg)
		case msg := <-stream3.Chan():
			log.Printf("Message received: %s", msg)
		}
	}
}