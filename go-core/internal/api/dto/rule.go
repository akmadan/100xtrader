package dto

import (
	"time"

	"go-core/internal/data"
)

// CreateRuleRequest represents the request to create a new rule
type CreateRuleRequest struct {
	UserID      int               `json:"user_id" validate:"required"`
	Name        string            `json:"name" validate:"required,min=1,max=255"`
	Description string            `json:"description" validate:"required,min=1,max=1000"`
	Category    data.RuleCategory `json:"category" validate:"required"`
}

// UpdateRuleRequest represents the request to update an existing rule
type UpdateRuleRequest struct {
	ID          string            `json:"id" validate:"required"`
	UserID      int               `json:"user_id" validate:"required"`
	Name        string            `json:"name" validate:"required,min=1,max=255"`
	Description string            `json:"description" validate:"required,min=1,max=1000"`
	Category    data.RuleCategory `json:"category" validate:"required"`
}

// RuleResponse represents the response for rule operations
type RuleResponse struct {
	ID          string            `json:"id"`
	UserID      int               `json:"user_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    data.RuleCategory `json:"category"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// GetRulesRequest represents the request to get rules with pagination
type GetRulesRequest struct {
	UserID   int                `json:"user_id" validate:"required"`
	Category *data.RuleCategory `json:"category,omitempty"`
	Limit    int                `json:"limit" validate:"min=1,max=100"`
	Offset   int                `json:"offset" validate:"min=0"`
}

// GetRulesResponse represents the response for getting rules
type GetRulesResponse struct {
	Rules      []RuleResponse     `json:"rules"`
	Pagination PaginationResponse `json:"pagination"`
}

// DeleteRuleRequest represents the request to delete a rule
type DeleteRuleRequest struct {
	ID     string `json:"id" validate:"required"`
	UserID int    `json:"user_id" validate:"required"`
}
