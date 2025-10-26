package dto

import (
	"time"

	"go-core/internal/data"
)

// CreateTradeRequest represents the request to create a new trade
type CreateTradeRequest struct {
	UserID         int                      `json:"user_id" validate:"required"`
	Symbol         string                   `json:"symbol" validate:"required"`
	MarketType     data.MarketType          `json:"market_type" validate:"required"`
	EntryDate      string                   `json:"entry_date" validate:"required"`
	EntryPrice     float64                  `json:"entry_price" validate:"required,gt=0"`
	Quantity       int                      `json:"quantity" validate:"required,gt=0"`
	TotalAmount    float64                  `json:"total_amount" validate:"required,gt=0"`
	ExitPrice      *float64                 `json:"exit_price,omitempty"`
	Direction      data.TradeDirection      `json:"direction" validate:"required"`
	StopLoss       *float64                 `json:"stop_loss,omitempty"`
	Target         *float64                 `json:"target,omitempty"`
	Strategy       string                   `json:"strategy" validate:"required"`
	OutcomeSummary data.OutcomeSummary      `json:"outcome_summary" validate:"required"`
	TradeAnalysis  *string                  `json:"trade_analysis,omitempty"`
	RulesFollowed  []string                 `json:"rules_followed,omitempty"`
	Screenshots    []string                 `json:"screenshots,omitempty"`
	Psychology     *CreatePsychologyRequest `json:"psychology,omitempty"`
}

// UpdateTradeRequest represents the request to update an existing trade
type UpdateTradeRequest struct {
	ID             string                   `json:"id" validate:"required"`
	UserID         int                      `json:"user_id" validate:"required"`
	Symbol         string                   `json:"symbol" validate:"required"`
	MarketType     data.MarketType          `json:"market_type" validate:"required"`
	EntryDate      string                   `json:"entry_date" validate:"required"`
	EntryPrice     float64                  `json:"entry_price" validate:"required,gt=0"`
	Quantity       int                      `json:"quantity" validate:"required,gt=0"`
	TotalAmount    float64                  `json:"total_amount" validate:"required,gt=0"`
	ExitPrice      *float64                 `json:"exit_price,omitempty"`
	Direction      data.TradeDirection      `json:"direction" validate:"required"`
	StopLoss       *float64                 `json:"stop_loss,omitempty"`
	Target         *float64                 `json:"target,omitempty"`
	Strategy       string                   `json:"strategy" validate:"required"`
	OutcomeSummary data.OutcomeSummary      `json:"outcome_summary" validate:"required"`
	TradeAnalysis  *string                  `json:"trade_analysis,omitempty"`
	RulesFollowed  []string                 `json:"rules_followed,omitempty"`
	Screenshots    []string                 `json:"screenshots,omitempty"`
	Psychology     *CreatePsychologyRequest `json:"psychology,omitempty"`
}

// CreatePsychologyRequest represents psychology data for a trade
type CreatePsychologyRequest struct {
	EntryConfidence    int      `json:"entry_confidence" validate:"required,min=1,max=10"`
	SatisfactionRating int      `json:"satisfaction_rating" validate:"required,min=1,max=10"`
	EmotionalState     string   `json:"emotional_state" validate:"required"`
	MistakesMade       []string `json:"mistakes_made,omitempty"`
	LessonsLearned     *string  `json:"lessons_learned,omitempty"`
}

// TradeResponse represents the response for trade operations
type TradeResponse struct {
	ID             string              `json:"id"`
	UserID         int                 `json:"user_id"`
	Symbol         string              `json:"symbol"`
	MarketType     data.MarketType     `json:"market_type"`
	EntryDate      time.Time           `json:"entry_date"`
	EntryPrice     float64             `json:"entry_price"`
	Quantity       int                 `json:"quantity"`
	TotalAmount    float64             `json:"total_amount"`
	ExitPrice      *float64            `json:"exit_price,omitempty"`
	Direction      data.TradeDirection `json:"direction"`
	StopLoss       *float64            `json:"stop_loss,omitempty"`
	Target         *float64            `json:"target,omitempty"`
	Strategy       string              `json:"strategy"`
	OutcomeSummary data.OutcomeSummary `json:"outcome_summary"`
	TradeAnalysis  *string             `json:"trade_analysis,omitempty"`
	RulesFollowed  []string            `json:"rules_followed,omitempty"`
	Screenshots    []string            `json:"screenshots,omitempty"`
	Psychology     *PsychologyResponse `json:"psychology,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

// PsychologyResponse represents psychology data in responses
type PsychologyResponse struct {
	EntryConfidence    int      `json:"entry_confidence"`
	SatisfactionRating int      `json:"satisfaction_rating"`
	EmotionalState     string   `json:"emotional_state"`
	MistakesMade       []string `json:"mistakes_made,omitempty"`
	LessonsLearned     *string  `json:"lessons_learned,omitempty"`
}

// GetTradesRequest represents the request to get trades with pagination
type GetTradesRequest struct {
	UserID int `json:"user_id" validate:"required"`
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

// GetTradesResponse represents the response for getting trades
type GetTradesResponse struct {
	Trades     []TradeResponse    `json:"trades"`
	Pagination PaginationResponse `json:"pagination"`
}

// DeleteTradeRequest represents the request to delete a trade
type DeleteTradeRequest struct {
	ID     string `json:"id" validate:"required"`
	UserID int    `json:"user_id" validate:"required"`
}
