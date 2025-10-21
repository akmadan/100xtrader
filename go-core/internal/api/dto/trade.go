package dto

import (
	"time"
)

// TradeCreateRequest represents the request to create a new trade
type TradeCreateRequest struct {
	UserID   int                        `json:"user_id" validate:"required,min=1"`
	Market   string                     `json:"market" validate:"required,oneof=stock option crypto futures forex index"`
	Symbol   string                     `json:"symbol" validate:"required,min=1,max=50"`
	Target   float64                    `json:"target" validate:"required,gt=0"`
	StopLoss float64                    `json:"stoploss" validate:"required,gt=0"`
	Tags     []string                   `json:"tags,omitempty"`    // Optional tags for the trade
	Journal  *TradeJournalCreateRequest `json:"journal,omitempty"` // Optional journal entry
	Actions  []TradeActionCreateRequest `json:"actions,omitempty"` // Optional trade actions (buy/sell)
}

// TradeUpdateRequest represents the request to update a trade
type TradeUpdateRequest struct {
	Market   *string  `json:"market,omitempty" validate:"omitempty,oneof=stock option crypto futures forex index"`
	Symbol   *string  `json:"symbol,omitempty" validate:"omitempty,min=1,max=50"`
	Target   *float64 `json:"target,omitempty" validate:"omitempty,gt=0"`
	StopLoss *float64 `json:"stoploss,omitempty" validate:"omitempty,gt=0"`
}

// TradeResponse represents the response for trade data
type TradeResponse struct {
	ID          string                    `json:"id"`
	UserID      int                       `json:"user_id"`
	Market      string                    `json:"market"`
	Symbol      string                    `json:"symbol"`
	Target      float64                   `json:"target"`
	StopLoss    float64                   `json:"stoploss"`
	CreatedAt   time.Time                 `json:"created_at"`
	Journal     *TradeJournalResponse     `json:"journal,omitempty"`     // Optional journal data
	Actions     []TradeActionResponse     `json:"actions,omitempty"`     // Trade actions (buy/sell)
	Tags        []TagResponse             `json:"tags,omitempty"`        // Tags associated with trade
	Screenshots []TradeScreenshotResponse `json:"screenshots,omitempty"` // Screenshots in journal
}

// TradeListResponse represents the response for listing trades
type TradeListResponse struct {
	Trades []TradeResponse `json:"trades"`
	Total  int             `json:"total"`
}

// TradeJournalCreateRequest represents the request to create a trade journal
type TradeJournalCreateRequest struct {
	TradeID    string `json:"trade_id,omitempty" validate:"omitempty"` // Optional when embedded in trade creation
	Notes      string `json:"notes" validate:"required,min=1"`
	Confidence int    `json:"confidence" validate:"required,min=0,max=10"`
}

// TradeJournalCreateStandaloneRequest represents the request to create a trade journal standalone
type TradeJournalCreateStandaloneRequest struct {
	TradeID    string `json:"trade_id" validate:"required"`
	Notes      string `json:"notes" validate:"required,min=1"`
	Confidence int    `json:"confidence" validate:"required,min=0,max=10"`
}

// TradeJournalUpdateRequest represents the request to update a trade journal
type TradeJournalUpdateRequest struct {
	Notes      *string `json:"notes,omitempty" validate:"omitempty,min=1"`
	Confidence *int    `json:"confidence,omitempty" validate:"omitempty,min=0,max=10"`
}

// TradeJournalResponse represents the response for trade journal data
type TradeJournalResponse struct {
	ID          int                       `json:"id"`
	TradeID     string                    `json:"trade_id"`
	Notes       string                    `json:"notes"`
	Confidence  int                       `json:"confidence"`
	CreatedAt   time.Time                 `json:"created_at"`
	Screenshots []TradeScreenshotResponse `json:"screenshots,omitempty"`
}

// TradeActionCreateRequest represents the request to create a trade action
type TradeActionCreateRequest struct {
	TradeID   string    `json:"trade_id" validate:"required"`
	Action    string    `json:"action" validate:"required,oneof=buy sell"`
	TradeTime time.Time `json:"trade_time" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
	Price     float64   `json:"price" validate:"required,gt=0"`
	Fee       float64   `json:"fee" validate:"gte=0"`
}

// TradeActionResponse represents the response for trade action data
type TradeActionResponse struct {
	ID        int       `json:"id"`
	TradeID   string    `json:"trade_id"`
	Action    string    `json:"action"`
	TradeTime time.Time `json:"trade_time"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	Fee       float64   `json:"fee"`
	CreatedAt time.Time `json:"created_at"`
}

// TradeScreenshotCreateRequest represents the request to create a trade screenshot
type TradeScreenshotCreateRequest struct {
	TradeJournalID int    `json:"trade_journal_id" validate:"required,min=1"`
	URL            string `json:"url" validate:"required,url"`
}

// TradeScreenshotResponse represents the response for trade screenshot data
type TradeScreenshotResponse struct {
	ID             int       `json:"id"`
	TradeJournalID int       `json:"trade_journal_id"`
	URL            string    `json:"url"`
	CreatedAt      time.Time `json:"created_at"`
}

// TagCreateRequest represents the request to create a tag
type TagCreateRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// TagResponse represents the response for tag data
type TagResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// TradeTagCreateRequest represents the request to add a tag to a trade
type TradeTagCreateRequest struct {
	TradeID string `json:"trade_id" validate:"required"`
	TagID   int    `json:"tag_id" validate:"required,min=1"`
}

// TradeWithJournalResponse represents a trade with its journal (optional)
type TradeWithJournalResponse struct {
	TradeResponse
	Journal *TradeJournalResponse `json:"journal,omitempty"`
}

// TradeActionRequest represents the request to add a trade action
type TradeActionRequest struct {
	Action    string    `json:"action" validate:"required,oneof=buy sell"`
	TradeTime time.Time `json:"trade_time" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
	Price     float64   `json:"price" validate:"required,gt=0"`
	Fee       float64   `json:"fee" validate:"gte=0"`
}

// TradeJournalRequest represents the request to update a trade journal
type TradeJournalRequest struct {
	Notes      string `json:"notes" validate:"required,min=1"`
	Confidence int    `json:"confidence" validate:"required,min=0,max=10"`
}

// TradeScreenshotRequest represents the request to add a screenshot
type TradeScreenshotRequest struct {
	URL string `json:"url" validate:"required,url"`
}

// TagRequest represents the request to create a tag (alias for TagCreateRequest)
type TagRequest = TagCreateRequest

// TagListResponse represents the response for listing tags
type TagListResponse struct {
	Tags  []TagResponse `json:"tags"`
	Total int           `json:"total"`
}

// TagTradeRequest represents the request to add a tag to a trade
type TagTradeRequest struct {
	TradeID string `json:"trade_id" validate:"required"`
	TagID   int    `json:"tag_id" validate:"required,min=1"`
}
