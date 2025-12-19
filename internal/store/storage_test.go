package store

import (
	"goticker/internal/model"
	"testing"
)

// 1. UNIT TEST: Test correctness
func TestGetAveragePrice(t *testing.T) {
	// Setup
	s := NewStorage()
	
	// Inject fake data
	s.UpdatePrice(model.NewCryptoData("Binance", "BTC", 50000.00))
	s.UpdatePrice(model.NewCryptoData("Coinbase", "BTC", 50200.00))
	
	// Execute
	avg, ok := s.GetAveragePrice("BTC")
	
	// Assert
	if !ok {
		t.Fatalf("Expected price data, got none")
	}
	
	expected := 50100.00 // (50000 + 50200) / 2
	if avg != expected {
		t.Errorf("Expected average %.2f, got %.2f", expected, avg)
	}
}

// 2. BENCHMARK TEST: Test speed
// Go runs this function thousands of times to measure performance.
func BenchmarkGetAveragePrice(b *testing.B) {
	// Setup: Fill the store with some dummy data
	s := NewStorage()
	s.UpdatePrice(model.NewCryptoData("Binance", "BTC", 50000.00))
	s.UpdatePrice(model.NewCryptoData("Coinbase", "BTC", 50200.00))
	s.UpdatePrice(model.NewCryptoData("Kraken", "BTC", 50100.00))

	// Reset timer so setup doesn't count towards the speed test
	b.ResetTimer()

	// The loop that tests performance
	for i := 0; i < b.N; i++ {
		s.GetAveragePrice("BTC")
	}
}