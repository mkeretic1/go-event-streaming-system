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

	// receive new events on stream
	stream1, err := c.Subscribe("stream1", eventstream.SubscriptionModeNew)
	if err != nil {
		log.Fatal("Error subscribing to stream:", err)
	}

	// receive all previous events and new events on stream
	stream2, err := c.Subscribe("stream2", eventstream.SubscriptionModeAll)
	if err != nil {
		log.Fatal("Error subscribing to stream:", err)
	}

	// receive all previous events on stream with id > afterId and new events on stream
	stream3, err := c.Subscribe("stream3", eventstream.SubscriptionModeAfter, 5)
	if err != nil {
		log.Fatal("Error subscribing to stream:", err)
	}

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
