package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type EventType string

const (
	EventCreate EventType = "create"
	EventUpdate EventType = "update"
	EventDelete EventType = "delete"
)

type TaskEvent struct {
	Type   EventType `json:"type"`
	TaskID int       `json:"taskId"`
	Data   any       `json:"data,omitempty"`
}

type Broker struct {
	clients   map[chan []byte]bool
	clientsMu sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[chan []byte]bool),
	}
}

func (b *Broker) RegisterClient() chan []byte {
	b.clientsMu.Lock()
	defer b.clientsMu.Unlock()

	ch := make(chan []byte, 100)
	b.clients[ch] = true
	return ch
}

func (b *Broker) UnregisterClient(ch chan []byte) {
	b.clientsMu.Lock()
	defer b.clientsMu.Unlock()
	defer fmt.Println("Client unregistered")
	delete(b.clients, ch)
}

func (b *Broker) Cleanup() {
	b.clientsMu.Lock()
	defer b.clientsMu.Unlock()

	for client := range b.clients {
		close(client)
		delete(b.clients, client)
	}
}

func (b *Broker) Broadcast(event TaskEvent) {
	log.Printf("Broadcasting event: %+v", event)

	b.clientsMu.RLock()
	defer b.clientsMu.RUnlock()

	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	for client := range b.clients {
		select {
		case client <- eventJSON:
		default:
			log.Printf("Couldn't send event to client, channel full")
		}
	}
}
