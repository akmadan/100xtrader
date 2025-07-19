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

type AddEnvironmentRequest struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Volatility  string `json:"volatility"`
	Trend       string `json:"trend"`
	Liquidity   string `json:"liquidity"`
}

type StartEnvironmentRequest struct {
	ID string `json:"id" binding:"required"`
}

type StartEnvironmentResponse struct {
	Status string `json:"status"`
}

type AddEnvironmentResponse struct {
	Status string `json:"status"`
}

type SetCurrentEnvironmentRequest struct {
	ID string `json:"id" binding:"required"`
}

type SetCurrentEnvironmentResponse struct {
	Status string `json:"status"`
}

type CurrentEnvironmentResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Volatility  string `json:"volatility"`
	Trend       string `json:"trend"`
	Liquidity   string `json:"liquidity"`
}
