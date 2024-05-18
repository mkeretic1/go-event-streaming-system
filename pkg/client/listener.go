package client

import (
	"bufio"
	"io"
	"log"
	"strings"
)

// EventListener manages a channel to receive messages from a subscription.
type EventListener struct {
	ch chan string
}

// Chan returns the channel that receives event messages.
func (s *EventListener) Chan() chan string {
	return s.ch
}

// createListener initializes an event listener for incoming messages.
func (c *Client) createListener() *EventListener {
	eventListener := &EventListener{
		ch: make(chan string),
	}

	go c.listen(eventListener)
	return eventListener
}

// listen continuously reads messages from the connection and sends them to the listener channel.
func (c *Client) listen(eventListener *EventListener) {
	reader := bufio.NewReader(c.conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Fatalf("Connection closed by server")
			}
			log.Println("Error reading message:", err)
			continue
		}

		msg = strings.TrimSuffix(msg, "")
		eventListener.ch <- msg
	}
}
