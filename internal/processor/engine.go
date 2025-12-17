package processor

import (
	"fmt"
	"goticker/internal/api" // Import api to see the Hub
	"goticker/internal/model"
	"goticker/internal/store"
)

// Update function signature to accept *api.Hub
func StartProcessor(ch <-chan model.CryptoData, db *store.Storage, hub *api.Hub) {
	for data := range ch {
		// 1. Update Store
		db.UpdatePrice(data)

		// 2. Calculate Average
		avgPrice, _ := db.GetAveragePrice(data.Symbol)

		// 3. Log to Console
		fmt.Printf("📊 %s: $%.2f (Avg: $%.2f)\n", data.Symbol, data.Price, avgPrice)

		// 4. BROADCAST TO WEB! (The new part)
		// We create a new payload that includes the average
		payload := data
		payload.Price = avgPrice // Let's send the average price to the frontend
		
		hub.BroadcastToClients(payload)
	}
}