package dto

import (
	"time"
)

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Code    int               `json:"code"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// DateRangeRequest represents a date range filter
type DateRangeRequest struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// SortRequest represents sorting parameters
type SortRequest struct {
	Field string `json:"field" validate:"required"`
	Order string `json:"order" validate:"required,oneof=asc desc"`
}

// SearchRequest represents search parameters
type SearchRequest struct {
	Query string `json:"query" validate:"required,min=1"`
	Field string `json:"field,omitempty"`
}

// FilterRequest represents generic filter parameters
type FilterRequest struct {
	Field string      `json:"field" validate:"required"`
	Value interface{} `json:"value" validate:"required"`
	Op    string      `json:"op" validate:"required,oneof=eq ne gt gte lt lte like in"`
}

// StatsResponse represents generic statistics
type StatsResponse struct {
	Total     int                    `json:"total"`
	Breakdown map[string]interface{} `json:"breakdown"`
	Trends    []TrendData            `json:"trends,omitempty"`
}

// TrendData represents trend information
type TrendData struct {
	Date  time.Time   `json:"date"`
	Value interface{} `json:"value"`
	Label string      `json:"label,omitempty"`
}

// BulkOperationRequest represents bulk operations
type BulkOperationRequest struct {
	IDs    []int       `json:"ids" validate:"required,min=1"`
	Action string      `json:"action" validate:"required,oneof=delete update"`
	Data   interface{} `json:"data,omitempty"`
}

// BulkOperationResponse represents bulk operation results
type BulkOperationResponse struct {
	SuccessCount int      `json:"success_count"`
	ErrorCount   int      `json:"error_count"`
	Errors       []string `json:"errors,omitempty"`
	Message      string   `json:"message"`
}
