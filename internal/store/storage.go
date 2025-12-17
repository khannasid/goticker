package store

import (
	"goticker/internal/model"
	"sync"
)

// Storage handles the thread-safe storage of prices.
type Storage struct {
	mu sync.RWMutex
	// Map of Symbol -> Exchange -> Price
	// Example: prices["BTC"]["Binance"] = 65000.00
	prices map[string]map[string]float64
}

// NewStorage initializes the map.
func NewStorage() *Storage {
	return &Storage{
		prices: make(map[string]map[string]float64),
	}
}

// UpdatePrice safely updates the price for a specific exchange and symbol.
func (s *Storage) UpdatePrice(data model.CryptoData) {
	// 1. LOCK: We are writing data. Block everyone else.
	s.mu.Lock()
	// defer ensures Unlock runs even if the function crashes or returns early.
	defer s.mu.Unlock() 

	// Initialize map for symbol if it doesn't exist
	if _, ok := s.prices[data.Symbol]; !ok {
		s.prices[data.Symbol] = make(map[string]float64)
	}

	// Update the price
	s.prices[data.Symbol][data.ExchangeName] = data.Price
}

// GetAveragePrice calculates the average price for a symbol across all exchanges.
func (s *Storage) GetAveragePrice(symbol string) (float64, bool) {
	// 1. RLOCK: We are only reading. Allow other readers, but block writers.
	s.mu.RLock()
	defer s.mu.RUnlock()

	exchanges, ok := s.prices[symbol]
	if !ok || len(exchanges) == 0 {
		return 0, false
	}

	var sum float64
	for _, price := range exchanges {
		sum += price
	}

	return sum / float64(len(exchanges)), true
}