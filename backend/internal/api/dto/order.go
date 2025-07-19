package dto

// OrderCreateRequest is the request body for creating an order
// Validation tags are used by Gin and go-playground/validator

type OrderCreateRequest struct {
	User      string  `json:"user" binding:"required"`
	Symbol    string  `json:"symbol" binding:"required"`
	Side      string  `json:"side" binding:"required,oneof=buy sell"`
	Type      string  `json:"type" binding:"required,oneof=market limit stop"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	Price     float64 `json:"price" binding:"required,gt=0"`
	SessionID string  `json:"session_id" binding:"required"`
	Source    string  `json:"source" binding:"required,oneof=user ai"`
}

// OrderResponse is the response body for an order

type OrderResponse struct {
	ID        int64   `json:"id"`
	User      string  `json:"user"`
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Type      string  `json:"type"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
