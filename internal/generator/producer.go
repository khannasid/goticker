package generator

import (
	"encoding/json"
	"fmt"
	"goticker/internal/model"
	"net/http"
	"strconv"
	"time"
)

// httpClient is a shared client with a timeout to prevent hanging requests.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// fetchBinancePrice fetches real data from Binance API.
func fetchBinancePrice(symbol string) (float64, error) {
	// Binance uses symbols like "BTCUSDT"
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol)
	
	resp, err := httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("binance API returned status: %d", resp.StatusCode)
	}

	var result BinanceAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(result.Price, 64)
}

// fetchCoinbasePrice fetches real data from Coinbase API.
func fetchCoinbasePrice(symbol string) (float64, error) {
	// Coinbase uses pairs like "BTC-USD"
	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-USD/spot", symbol)
	
	resp, err := httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("coinbase API returned status: %d", resp.StatusCode)
	}

	var result CoinbaseAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(result.Data.Amount, 64)
}

// StartRealExchangeWorker polls the real API and pushes data to the channel.
func StartRealExchangeWorker(exchangeName, symbol string, ch chan<- model.CryptoData) {
	// Ticker to control rate limiting (e.g., 1 request every 2 seconds)
	ticker := time.NewTicker(2 * time.Second) 
	defer ticker.Stop()

	fmt.Printf("🔌 Connected to %s API for %s\n", exchangeName, symbol)

	for range ticker.C {
		var price float64
		var err error

		// 1. Fetch from the specific API
		switch exchangeName {
		case "Binance":
			price, err = fetchBinancePrice(symbol)
		case "Coinbase":
			price, err = fetchCoinbasePrice(symbol)
		}

		// 2. Error Handling (Log but don't crash)
		if err != nil {
			fmt.Printf("⚠️ Error fetching from %s: %v\n", exchangeName, err)
			continue // Skip this tick, try again next time
		}

		// 3. Send to Channel (Fan-In)
		ch <- model.NewCryptoData(exchangeName, symbol, price)
	}
}