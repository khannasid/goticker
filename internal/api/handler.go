package api

import (
	"encoding/json"
	"time"
	"fmt"
	"goticker/internal/store"
	"net/http"
	"github.com/gorilla/websocket" // <--- Import this
)

// 1. Define the Upgrader
// This helper checks if the request is valid and upgrades connection to WS.
var upgrader = websocket.Upgrader{
	// CheckOrigin allows requests from localhost (important for development)
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

// Update the Server struct to hold the Hub
type Server struct {
	store *store.Storage
	hub   *Hub // <--- Add Hub dependency
}

// Update NewServer to accept the Hub
func NewServer(store *store.Storage, hub *Hub) *Server {
	return &Server{store: store, hub: hub}
}

// ... (Keep your existing GetPriceHandler here) ...

// GetPriceHandler handles GET /price?symbol=BTC
// It reads from the store and returns JSON.
func (s *Server) GetPriceHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse Query Params
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Missing 'symbol' query parameter", http.StatusBadRequest)
		return
	}

	// 2. Fetch data from the Store (Thread-Safe Read)
	price, ok := s.store.GetAveragePrice(symbol)
	if !ok {
		http.Error(w, "Symbol not found or no data yet", http.StatusNotFound)
		return
	}

	// 3. Create the Response Struct
	// We define a quick struct here just for the JSON response
	response := struct {
		Symbol    string  `json:"symbol"`
		Price     float64 `json:"price"`
		Timestamp int64   `json:"timestamp"`
	}{
		Symbol:    symbol,
		Price:     price,
		Timestamp: time.Now().UnixMilli(),
	}

	// 4. Send JSON Response
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder writes directly to the HTTP response stream (efficient!)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}	


// HandleWebSocket upgrades the HTTP connection to a persistent WebSocket.
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. Upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}

	// 2. Register the new client to the Hub
	s.hub.AddClient(conn)
	
	fmt.Println("✅ New Client Connected via WebSocket!")
}