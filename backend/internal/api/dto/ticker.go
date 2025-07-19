package dto

type TickerListResponse struct {
	Tickers []string `json:"tickers"`
}

type AddTickerRequest struct {
	Symbol string `json:"symbol" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

type AddTickerResponse struct {
	Status string `json:"status"`
}

type RemoveTickerResponse struct {
	Status string `json:"status"`
}
