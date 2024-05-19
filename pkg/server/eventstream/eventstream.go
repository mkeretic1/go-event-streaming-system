package eventstream

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

// EventSystem will handle managing streams and messages.
type EventSystem struct {
	lock    sync.RWMutex
	Streams map[string]EsStream `json:"streams"`
	Clients map[string]EsClient
}

var once = &sync.Once{}
var eventStreamSingleton *EventSystem

// GetInstance initializes and returns a singleton instance of EventSystem.
func GetInstance() *EventSystem {
	log.Println("Initializing event stream")
	once.Do(func() {
		streams, err := loadEvents()
		if err != nil {
			log.Fatalf("Failed to load event stream from file system: %v", err)
		}

		eventStreamSingleton = &EventSystem{
			lock:    sync.RWMutex{},
			Streams: streams,
			Clients: map[string]EsClient{},
		}
	})
	log.Println("Successfully initialized event stream")
	return eventStreamSingleton
}

// HandleConnection handles new client connections, processing incoming actions.
func (e *EventSystem) HandleConnection(conn net.Conn) {
	log.Printf("Received new connection from %s", conn.RemoteAddr().String())
	defer conn.Close()

	esClient := EsClient{
		uid:  strconv.FormatInt(time.Now().UnixNano(), 10),
		conn: conn,
	}
	e.Clients[esClient.uid] = esClient
	log.Printf("Assigned UID to new client connection: %s", esClient.uid)

	decoder := json.NewDecoder(conn)
	for {
		var action Action
		if err := decoder.Decode(&action); err != nil {
			if err == io.EOF { // Client closed the connection
				log.Printf("Client '%s' closed the connection, removing it from eventstream", esClient.uid)
				e.unregisterSubscriber(esClient)
				break
			}
			log.Printf("Error decoding Message from client '%s': %v", esClient.uid, err)
			continue
		}

		// Check the action and process accordingly
		switch action.ActionType {
		case ActionTypePublish:
			var publish Publish
			if err := json.Unmarshal(action.Data, &publish); err != nil {
				log.Printf("Error decoding Publish Action from client '%s': %v", esClient.uid, err)
				continue
			}

			log.Printf("Received PUBLISH request from client '%s' to stream '%s'", esClient.uid, publish.Stream)
			e.publishMessage(publish)
		case ActionTypeSubscribe:
			var subscribe Subscribe
			if err := json.Unmarshal(action.Data, &subscribe); err != nil {
				log.Printf("Error decoding Subscribe Action from client '%s': %v", esClient.uid, err)
				continue
			}

			log.Printf("Received SUBSRCIBE request from client '%s' to stream '%s'", esClient.uid, subscribe.Stream)
			e.registerSubscriber(esClient, subscribe)
		default:
			log.Println("Unknown action:", action.ActionType)
		}
	}
}

// registerSubscriber adds a client to a stream's subscriber list and sends previous events based on the subscription mode.
func (e *EventSystem) registerSubscriber(esClient EsClient, subscribe Subscribe) {
	e.lock.Lock()
	defer e.lock.Unlock()

	log.Printf("Registering subscriber '%s' to stream '%s' using subscription mode '%s'", esClient.uid, subscribe.Stream, subscribe.SubscriptionMode)

	// add subscriber to list
	stream := e.Streams[subscribe.Stream]
	stream.Subscribers = append(stream.Subscribers, esClient.uid)
	e.Streams[subscribe.Stream] = stream

	if subscribe.SubscriptionMode == SubscriptionModeAll || subscribe.SubscriptionMode == SubscriptionModeAfter {
		go e.sendPreviousEvents(esClient, subscribe, stream.Events)
	}
}

// publishMessage handles the publication of messages to a stream and notifies subscribers.
func (e *EventSystem) publishMessage(publish Publish) {
	e.lock.Lock()
	defer e.lock.Unlock()

	log.Printf("Publishing new event on stream '%s' to all subsribers", publish.Stream)

	stream := e.Streams[publish.Stream]
	event := EsEvent{
		Id:      len(stream.Events) + 1,
		Message: publish.Message,
	}

	// add event to stream
	stream.Events = append(stream.Events, event)
	e.Streams[publish.Stream] = stream

	if len(stream.Subscribers) > 0 {
		go e.sendEvent(stream, event)
	}

	go e.persistEvents()
}

// unregisterSubscriber removes a client from all subscriptions and the client-subscriber map.
func (e *EventSystem) unregisterSubscriber(esClient EsClient) {
	for s, stream := range e.Streams {
		stream.Subscribers = removeValueFromSlice(stream.Subscribers, esClient.uid)
		e.Streams[s] = stream
	}
	delete(e.Clients, esClient.uid)
}

// sendPreviousEvents sends previous events to a newly subscribed client based on the subscription mode.
func (e *EventSystem) sendPreviousEvents(esClient EsClient, subscribe Subscribe, events []EsEvent) {
	switch subscribe.SubscriptionMode {
	case SubscriptionModeAll:
		log.Printf("Sending all previous messages on stream '%s' to client '%s'", subscribe.Stream, esClient.uid)
		for _, event := range events {
			e.writeToClient(esClient, event.Message)
		}
	case SubscriptionModeAfter:
		log.Printf("Sending all messages after id '%d' on stream '%s' to client '%s'", subscribe.AfterId, subscribe.Stream, esClient.uid)
		if len(events) > subscribe.AfterId {
			for i := subscribe.AfterId; i < len(events); i++ {
				e.writeToClient(esClient, events[i].Message)
			}
		}
	}
}

// sendEvent broadcasts a new event to all subscribers of the stream.
func (e *EventSystem) sendEvent(stream EsStream, event EsEvent) {
	for _, subscriber := range stream.Subscribers {
		e.writeToClient(e.Clients[subscriber], event.Message)
	}
}

// writeToClient sends a message to a client's connection.
func (e *EventSystem) writeToClient(esClient EsClient, message string) {
	_, err := esClient.conn.Write([]byte(message + "\n"))
	if err != nil {
		log.Printf("Failed to write to client: %v", err)
	}
}
