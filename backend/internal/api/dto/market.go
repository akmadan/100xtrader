package dto

type OrderBookLevel struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

type OrderBookResponse struct {
	Symbol string           `json:"symbol"`
	Bids   []OrderBookLevel `json:"bids"`
	Asks   []OrderBookLevel `json:"asks"`
	Time   string           `json:"time"`
}

type TradeTickResponse struct {
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Side      string  `json:"side"`
	Timestamp string  `json:"timestamp"`
}

type OHLCResponse struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
	Start  string  `json:"start"`
}

type MarketDataResponse struct {
	Symbol        string              `json:"symbol"`
	LastPrice     float64             `json:"last_price"`
	LastTradeTime string              `json:"last_trade_time"`
	OHLC          OHLCResponse        `json:"ohlc"`
	RecentTrades  []TradeTickResponse `json:"recent_trades"`
}
