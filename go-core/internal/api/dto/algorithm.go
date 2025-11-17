package dto

import (
	"go-core/internal/data"
)

// CreateAlgorithmRequest represents the request to create a new algorithm
type CreateAlgorithmRequest struct {
	UserID      int                    `json:"user_id" validate:"required"`
	Name        string                 `json:"name" validate:"required,min=1,max=100"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=500"`
	Code        string                 `json:"code" validate:"required,min=1"`
	Status      data.AlgorithmStatus   `json:"status" validate:"required,oneof=draft live paused archived"`
	Symbol      string                 `json:"symbol" validate:"required,min=1"`
	Timeframe   data.Timeframe         `json:"timeframe" validate:"required,oneof=1m 5m 15m 30m 1h 4h 1d 1w"`
	ExecutionMode data.ExecutionMode   `json:"execution_mode" validate:"required,oneof=paper_trading live_trading"`
	Broker      *data.TradingBroker    `json:"broker,omitempty"`
	Enabled     bool                  `json:"enabled"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
}

// UpdateAlgorithmRequest represents the request to update an algorithm
type UpdateAlgorithmRequest struct {
	Name        *string                `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=500"`
	Code        *string                `json:"code,omitempty" validate:"omitempty,min=1"`
	Status      *data.AlgorithmStatus  `json:"status,omitempty" validate:"omitempty,oneof=draft live paused archived"`
	Symbol      *string                `json:"symbol,omitempty" validate:"omitempty,min=1"`
	Timeframe   *data.Timeframe        `json:"timeframe,omitempty" validate:"omitempty,oneof=1m 5m 15m 30m 1h 4h 1d 1w"`
	ExecutionMode *data.ExecutionMode `json:"execution_mode,omitempty" validate:"omitempty,oneof=paper_trading live_trading"`
	Broker      *data.TradingBroker    `json:"broker,omitempty"`
	Enabled     *bool                  `json:"enabled,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

// AlgorithmResponse represents the response for algorithm operations
type AlgorithmResponse struct {
	ID           string                 `json:"id"`
	UserID       int                    `json:"user_id"`
	Name         string                 `json:"name"`
	Description  *string                `json:"description,omitempty"`
	Code         string                 `json:"code"`
	Status       data.AlgorithmStatus  `json:"status"`
	Symbol       string                 `json:"symbol"`
	Timeframe    data.Timeframe         `json:"timeframe"`
	ExecutionMode data.ExecutionMode    `json:"execution_mode"`
	Broker       *data.TradingBroker    `json:"broker,omitempty"`
	Enabled      bool                   `json:"enabled"`
	Config       map[string]interface{} `json:"config,omitempty"`
	State        map[string]interface{} `json:"state,omitempty"`
	LastRunAt    *string                `json:"last_run_at,omitempty"`
	LastSignal   *string                `json:"last_signal,omitempty"`
	TotalTrades  int                    `json:"total_trades"`
	WinRate      float64                `json:"win_rate"`
	TotalPnL     float64                `json:"total_pnl"`
	Version      int                    `json:"version"`
	Tags         []string               `json:"tags,omitempty"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
}

// GetAlgorithmsResponse represents the response for getting multiple algorithms
type GetAlgorithmsResponse struct {
	Algorithms []AlgorithmResponse `json:"algorithms"`
	Pagination PaginationResponse  `json:"pagination"`
}

