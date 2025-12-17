package generator

// BinanceAPIResponse matches: {"symbol":"BTCUSDT","price":"65000.00"}
// Note: Binance sends price as a STRING, not a float! We must parse it.
type BinanceAPIResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"` 
}

// CoinbaseAPIResponse matches: {"data":{"base":"BTC","currency":"USD","amount":"64000.50"}}
type CoinbaseAPIResponse struct {
	Data struct {
		Base   string `json:"base"`
		Amount string `json:"amount"` // Coinbase also sends string!
	} `json:"data"`
}