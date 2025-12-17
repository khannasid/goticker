package api

import (
	"fmt"
	"goticker/internal/model"
	"sync"
	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	// Registered clients.
	// We use a map because deleting a client (disconnecting) is O(1) fast.
	clients map[*websocket.Conn]bool

	// Inbound messages from the Processor.
	broadcast chan model.CryptoData

	// Mutex to protect the clients map (Thread Safety!).
	mu sync.Mutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		broadcast: make(chan model.CryptoData),
		clients:   make(map[*websocket.Conn]bool),
	}
}

// Run starts the Hub's main loop.
// It waits for messages and then broadcasts them.
func (h *Hub) Run() {
	for {
		// Wait for the next message from the Processor
		msg := <-h.broadcast

		// Lock the map so no one can connect/disconnect while we are broadcasting
		h.mu.Lock()
		for client := range h.clients {
			// Write JSON to the websocket
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Printf("⚠️ Client disconnected error: %v\n", err)
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}

// BroadcastToClients is a helper function for the Processor to call.
func (h *Hub) BroadcastToClients(data model.CryptoData) {
	h.broadcast <- data
}

// AddClient registers a new websocket connection safely.
func (h *Hub) AddClient(conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()
}