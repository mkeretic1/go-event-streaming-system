package client

import (
	"encoding/json"
	eventstream "event-streaming-system/pkg/server/eventstream"
	"log"
	"net"
)

type Client struct {
	conn net.Conn
}

// Connect establishes a connection to the server at the specified address
func Connect(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

// Publish sends a message to a specified stream on the server.
func (c *Client) Publish(stream string, message string) error {
	action, err := eventstream.NewPublishAction(stream, message)
	if err != nil {
		log.Fatalf("Failed to instantiate publish action: %v", err)
	}

	data, err := json.Marshal(action)
	if err != nil {
		log.Fatalf("Error marshaling data to JSON: %v", err)
	}

	_, err = c.conn.Write(data)
	if err != nil {
		log.Fatalf("Error Publishing to stream: %v", err)
	}
	return err
}

// Subscribe sets up a subscription to a stream and listens for messages from server.
func (c *Client) Subscribe(stream string, subscriptionMode eventstream.SubscriptionMode, afterId ...int) *EventListener {
	if subscriptionMode == eventstream.SubscriptionModeAfter && len(afterId) == 0 {
		log.Fatalf("Used SubscriptionModeAfter, but 'afterId' was not specified")
	}

	action, err := eventstream.NewSubscribeAction(stream, subscriptionMode, afterId...)
	if err != nil {
		log.Fatalf("Failed to instantiate subscriber action: %v", err)
	}

	data, err := json.Marshal(action)
	if err != nil {
		log.Fatalf("Error marshaling data to JSON: %v", err)
	}

	_, err = c.conn.Write(data)
	if err != nil {
		log.Fatalf("Error Subscribing to stream: %v", err)
	}

	return c.createListener()
}

// Close closes the connection to the server
func (c *Client) Close() {
	c.conn.Close()
}
