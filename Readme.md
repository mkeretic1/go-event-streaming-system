# Go Event Streaming System

## Overview

The Go Event Streaming System is a simple system for event streaming built entirely in Go. It features separate client and server components that communicate over TCP, allowing clients to publish messages to streams and subscribe to receive messages from streams.

## Architecture

The system is divided into several parts:
- **Server**: Handles incoming TCP connections, manages streams, and distributes messages.
- **Client**: Can either publish messages to a stream or subscribe to receive updates from streams.

## Getting Started

### Prerequisites
- Go 1.21+ installed on your machine

### Installation

Clone the repository to your local machine:

```bash
git clone git@github.com:mkeretic1/go-event-streaming-system.git
cd go-event-streaming-system
```

## Running

### Running the Server
By default, server starts on `localhost:8080`
```bash
cd cmd/server
go run main.go
```

### Using the Client

#### Subscribing to a Stream
Here's how to subscribe to a stream called "stream1". See the full example [here](/cmd/client-subscriber/main.go).

```go
func main() {
    c, err := client.Connect("localhost:8080")
    if err != nil {
        log.Println("Error connecting to server:", err)
        return
    }
    defer c.Close()

    log.Println("Client connected successfully")

    stream := c.Subscribe("stream1", eventstream.SubscriptionModeAll)

    for msg := <-stream.Chan() {
        log.Printf("Message received: %s", msg)
    }
```

#### Publishing a Message
Here's how to publish a message to a stream called "stream1".  See the full example [here](/cmd/client-publisher/main.go).

```go
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
}
```
