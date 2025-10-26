package dto

import (
	"time"

	"go-core/internal/data"
)

// CreateMistakeRequest represents the request to create a new mistake
type CreateMistakeRequest struct {
	UserID   int                  `json:"user_id" validate:"required"`
	Name     string               `json:"name" validate:"required,min=1,max=255"`
	Category data.MistakeCategory `json:"category" validate:"required"`
}

// UpdateMistakeRequest represents the request to update an existing mistake
type UpdateMistakeRequest struct {
	ID       string               `json:"id" validate:"required"`
	UserID   int                  `json:"user_id" validate:"required"`
	Name     string               `json:"name" validate:"required,min=1,max=255"`
	Category data.MistakeCategory `json:"category" validate:"required"`
}

// MistakeResponse represents the response for mistake operations
type MistakeResponse struct {
	ID        string               `json:"id"`
	UserID    int                  `json:"user_id"`
	Name      string               `json:"name"`
	Category  data.MistakeCategory `json:"category"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

// GetMistakesRequest represents the request to get mistakes with pagination
type GetMistakesRequest struct {
	UserID   int                   `json:"user_id" validate:"required"`
	Category *data.MistakeCategory `json:"category,omitempty"`
	Limit    int                   `json:"limit" validate:"min=1,max=100"`
	Offset   int                   `json:"offset" validate:"min=0"`
}

// GetMistakesResponse represents the response for getting mistakes
type GetMistakesResponse struct {
	Mistakes   []MistakeResponse  `json:"mistakes"`
	Pagination PaginationResponse `json:"pagination"`
}

// DeleteMistakeRequest represents the request to delete a mistake
type DeleteMistakeRequest struct {
	ID     string `json:"id" validate:"required"`
	UserID int    `json:"user_id" validate:"required"`
}
