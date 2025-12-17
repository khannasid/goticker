package main

import (
	"fmt"
	"goticker/internal/api"
	"goticker/internal/generator"
	"goticker/internal/model"
	"goticker/internal/processor"
	"goticker/internal/store"
	"net/http"
)

func main() {
	fmt.Println("🚀 Starting GoTicker PRO (Real APIs Mode)...")

	// 1. Infrastructure
	priceChannel := make(chan model.CryptoData, 10)
	priceDb := store.NewStorage()
	hub := api.NewHub()
	go hub.Run()

	// 2. Start Processor
	go processor.StartProcessor(priceChannel, priceDb, hub)

	// 3. Start REAL Data Producers
	// We are now hitting the actual internet.
	// Note: We use "BTC" and "ETH". The worker handles the specific API format (BTCUSDT vs BTC-USD).
	
	// Worker 1: Binance BTC
	go generator.StartRealExchangeWorker("Binance", "BTC", priceChannel)
	// Worker 2: Coinbase BTC
	go generator.StartRealExchangeWorker("Coinbase", "BTC", priceChannel)
	
	// Worker 3: Binance ETH
	go generator.StartRealExchangeWorker("Binance", "ETH", priceChannel)
	// Worker 4: Coinbase ETH
	go generator.StartRealExchangeWorker("Coinbase", "ETH", priceChannel)

	// 4. API Server
	apiServer := api.NewServer(priceDb, hub)
	http.HandleFunc("/price", apiServer.GetPriceHandler)
	http.HandleFunc("/ws", apiServer.HandleWebSocket)

	go func() {
		fmt.Println("🌐 Server running on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()

	select {}
}