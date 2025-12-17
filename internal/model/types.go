package model

import "time"

// CryptoData represents a single price update from an exchange.
// We use JSON tags (e.g., `json:"exchange"`) for when we send this to the frontend later.
type CryptoData struct {
	ExchangeName string  `json:"exchange"`
	Symbol       string  `json:"symbol"` // e.g., "BTC", "ETH"
	Price        float64 `json:"price"`
	Timestamp    int64   `json:"timestamp"` // Unix timestamp for efficient sorting
}

// NewCryptoData is a helper function to create a new instance easily.
func NewCryptoData(exchange, symbol string, price float64) CryptoData {
	return CryptoData{
		ExchangeName: exchange,
		Symbol:       symbol,
		Price:        price,
		Timestamp:    time.Now().UnixMilli(),
	}
}