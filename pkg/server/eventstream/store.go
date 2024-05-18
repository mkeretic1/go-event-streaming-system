package eventstream

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

var storeLock = &sync.RWMutex{}

const storeFile = "go-event-streaming-system.json"

// persistEvents writes all current stream events to a JSON file.
func (e *EventSystem) persistEvents() {
	storeLock.Lock()
	defer storeLock.Unlock()

	eventsData := make(map[string][]EsEvent)
	for streamName, stream := range e.Streams {
		eventsData[streamName] = stream.Events
	}

	data, err := json.MarshalIndent(eventsData, "", "    ")
	if err != nil {
		log.Println("Failed to marshal event stream: ", err)
		return
	}

	if err := os.WriteFile(storeFile, data, 0644); err != nil {
		log.Println("Failed to write event stream to file: ", err)
	}
}

// loadEvents loads streams events from a JSON file into a map.
func loadEvents() (map[string]EsStream, error) {
	storeLock.Lock()
	defer storeLock.Unlock()

	log.Println("Fetching existing events")

	data, err := os.ReadFile(storeFile)
	if err != nil {
		if os.IsNotExist(err) { // return empty stream map
			return make(map[string]EsStream), nil
		}
		return nil, fmt.Errorf("failed to read event stream file: %w", err)
	}

	if len(data) == 0 {
		return make(map[string]EsStream), nil
	}

	var eventData map[string][]EsEvent
	if err := json.Unmarshal(data, &eventData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event stream: %w", err)
	}

	eventStreams := make(map[string]EsStream)
	for streamName, events := range eventData {
		eventStreams[streamName] = EsStream{
			Events:      events,
			Subscribers: []string{},
		}
	}

	return eventStreams, nil
}
