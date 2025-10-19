package dto

import (
	"time"
)

// NoteCreateRequest represents the request to create a new note
type NoteCreateRequest struct {
	UserID           int       `json:"user_id" validate:"required,min=1"`
	Mood             string    `json:"mood" validate:"required,oneof=excited neutral low"`
	MarketCondition  string    `json:"market_condition" validate:"required,oneof=up down sideways"`
	MarketVolatility string    `json:"market_volatility" validate:"required,oneof=high medium low"`
	Summary          string    `json:"summary,omitempty"`
	Day              time.Time `json:"day" validate:"required"`
	Notes            string    `json:"notes,omitempty"`
}

// NoteUpdateRequest represents the request to update a note
type NoteUpdateRequest struct {
	Mood             *string    `json:"mood,omitempty" validate:"omitempty,oneof=excited neutral low"`
	MarketCondition  *string    `json:"market_condition,omitempty" validate:"omitempty,oneof=up down sideways"`
	MarketVolatility *string    `json:"market_volatility,omitempty" validate:"omitempty,oneof=high medium low"`
	Summary          *string    `json:"summary,omitempty"`
	Day              *time.Time `json:"day,omitempty"`
	Notes            *string    `json:"notes,omitempty"`
}

// NoteResponse represents the response for note data
type NoteResponse struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	Mood             string    `json:"mood"`
	MarketCondition  string    `json:"market_condition"`
	MarketVolatility string    `json:"market_volatility"`
	Summary          string    `json:"summary"`
	Day              time.Time `json:"day"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
}

// NoteListResponse represents the response for listing notes
type NoteListResponse struct {
	Notes []NoteResponse `json:"notes"`
	Total int            `json:"total"`
}

// NoteFilterRequest represents the request to filter notes
type NoteFilterRequest struct {
	UserID           *int       `json:"user_id,omitempty" validate:"omitempty,min=1"`
	Mood             *string    `json:"mood,omitempty" validate:"omitempty,oneof=excited neutral low"`
	MarketCondition  *string    `json:"market_condition,omitempty" validate:"omitempty,oneof=up down sideways"`
	MarketVolatility *string    `json:"market_volatility,omitempty" validate:"omitempty,oneof=high medium low"`
	StartDate        *time.Time `json:"start_date,omitempty"`
	EndDate          *time.Time `json:"end_date,omitempty"`
}

// NoteStatsResponse represents statistics for notes
type NoteStatsResponse struct {
	TotalNotes          int                 `json:"total_notes"`
	MoodBreakdown       map[string]int      `json:"mood_breakdown"`
	MarketCondition     map[string]int      `json:"market_condition_breakdown"`
	VolatilityBreakdown map[string]int      `json:"volatility_breakdown"`
	AvgNotesPerDay      float64             `json:"avg_notes_per_day"`
	RecentTrends        []NoteTrendResponse `json:"recent_trends"`
}

// NoteTrendResponse represents trend data for notes
type NoteTrendResponse struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
	Mood  string    `json:"mood"`
}

// DailyNoteRequest represents the request for daily note operations
type DailyNoteRequest struct {
	UserID int       `json:"user_id" validate:"required,min=1"`
	Date   time.Time `json:"date" validate:"required"`
}

// DailyNoteResponse represents the response for daily note data
type DailyNoteResponse struct {
	Date  time.Time      `json:"date"`
	Notes []NoteResponse `json:"notes"`
	Count int            `json:"count"`
	Mood  string         `json:"mood,omitempty"`
}
