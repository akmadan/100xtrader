package dto

import (
	"time"
)

// CreateStrategyRequest represents the request to create a new strategy
type CreateStrategyRequest struct {
	UserID      int    `json:"user_id" validate:"required"`
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,min=1,max=1000"`
}

// UpdateStrategyRequest represents the request to update an existing strategy
type UpdateStrategyRequest struct {
	ID          string `json:"id" validate:"required"`
	UserID      int    `json:"user_id" validate:"required"`
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,min=1,max=1000"`
}

// StrategyResponse represents the response for strategy operations
type StrategyResponse struct {
	ID          string    `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetStrategiesRequest represents the request to get strategies with pagination
type GetStrategiesRequest struct {
	UserID int `json:"user_id" validate:"required"`
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

// GetStrategiesResponse represents the response for getting strategies
type GetStrategiesResponse struct {
	Strategies []StrategyResponse `json:"strategies"`
	Pagination PaginationResponse `json:"pagination"`
}

// DeleteStrategyRequest represents the request to delete a strategy
type DeleteStrategyRequest struct {
	ID     string `json:"id" validate:"required"`
	UserID int    `json:"user_id" validate:"required"`
}
