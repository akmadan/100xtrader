package dto

type EnvironmentListItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Volatility  string `json:"volatility"`
	Trend       string `json:"trend"`
	Liquidity   string `json:"liquidity"`
}

type EnvironmentListResponse struct {
	Environments []EnvironmentListItem `json:"environments"`
}

type StartEnvironmentRequest struct {
	ID string `json:"id" binding:"required"`
}

type StartEnvironmentResponse struct {
	Status string `json:"status"`
}
