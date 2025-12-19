package api

import (
	"encoding/json"
	"goticker/internal/model"
	"goticker/internal/store"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPriceHandler(t *testing.T) {
	// 1. Setup Dependencies
	mockStore := store.NewStorage()
	mockStore.UpdatePrice(model.NewCryptoData("Binance", "BTC", 100.00))
	
	// We don't need a real Hub for this HTTP test, so we pass nil
	server := NewServer(mockStore, nil)

	// 2. Create a Request
	// "GET /price?symbol=BTC"
	req, _ := http.NewRequest("GET", "/price?symbol=BTC", nil)
	
	// 3. Create a ResponseRecorder (Acting as the Browser)
	rr := httptest.NewRecorder()

	// 4. Serve the Handler
	handler := http.HandlerFunc(server.GetPriceHandler)
	handler.ServeHTTP(rr, req)

	// 5. Assert Status Code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 6. Assert JSON Body
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response["symbol"] != "BTC" {
		t.Errorf("Expected symbol BTC, got %v", response["symbol"])
	}
	if response["price"].(float64) != 100.00 {
		t.Errorf("Expected price 100.00, got %v", response["price"])
	}
}