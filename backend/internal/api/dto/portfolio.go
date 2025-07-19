package dto

type PositionResponse struct {
	Symbol       string  `json:"symbol"`
	Quantity     float64 `json:"quantity"`
	AveragePrice float64 `json:"average_price"`
	MarketValue  float64 `json:"market_value"`
	PnL          float64 `json:"pnl"`
}

type TradeResponse struct {
	ID        int64   `json:"id"`
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
}

type PortfolioResponse struct {
	User        string             `json:"user"`
	CashBalance float64            `json:"cash_balance"`
	TotalValue  float64            `json:"total_value"`
	Positions   []PositionResponse `json:"positions"`
	PnL         float64            `json:"pnl"`
}

type PositionsListResponse struct {
	Positions []PositionResponse `json:"positions"`
}

type TradesListResponse struct {
	Trades []TradeResponse `json:"trades"`
}
