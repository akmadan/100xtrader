package data

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           int        `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Email        string     `json:"email" db:"email"`
	Phone        *string    `json:"phone" db:"phone"`
	LastSignedIn *time.Time `json:"last_signed_in" db:"last_signed_in"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// Trade represents a trading position
type Trade struct {
	ID             string           `json:"id" db:"id"`
	UserID         int              `json:"user_id" db:"user_id"`
	Symbol         string           `json:"symbol" db:"symbol"`
	MarketType     MarketType       `json:"market_type" db:"market_type"`
	EntryDate      time.Time        `json:"entry_date" db:"entry_date"`
	EntryPrice     float64          `json:"entry_price" db:"entry_price"`
	Quantity       int              `json:"quantity" db:"quantity"`
	TotalAmount    float64          `json:"total_amount" db:"total_amount"`
	ExitPrice      *float64         `json:"exit_price" db:"exit_price"`
	Direction      TradeDirection   `json:"direction" db:"direction"`
	StopLoss       *float64         `json:"stop_loss" db:"stop_loss"`
	Target         *float64         `json:"target" db:"target"`
	Strategy       string           `json:"strategy" db:"strategy"`
	OutcomeSummary OutcomeSummary   `json:"outcome_summary" db:"outcome_summary"`
	TradeAnalysis  *string          `json:"trade_analysis" db:"trade_analysis"`
	RulesFollowed  []string         `json:"rules_followed" db:"rules_followed"`
	Screenshots    []string         `json:"screenshots" db:"screenshots"`
	Psychology     *TradePsychology `json:"psychology" db:"psychology"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
}

// TradePsychology represents psychology information for a trade
type TradePsychology struct {
	EntryConfidence    int      `json:"entry_confidence" db:"entry_confidence"`       // 1-10 scale
	SatisfactionRating int      `json:"satisfaction_rating" db:"satisfaction_rating"` // 1-10 scale
	EmotionalState     string   `json:"emotional_state" db:"emotional_state"`
	MistakesMade       []string `json:"mistakes_made" db:"mistakes_made"`
	LessonsLearned     *string  `json:"lessons_learned" db:"lessons_learned"`
}

// Strategy represents a trading strategy
type Strategy struct {
	ID          string    `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Rule represents a trading rule
type Rule struct {
	ID          string       `json:"id" db:"id"`
	UserID      int          `json:"user_id" db:"user_id"`
	Name        string       `json:"name" db:"name"`
	Description string       `json:"description" db:"description"`
	Category    RuleCategory `json:"category" db:"category"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

// Mistake represents a trading mistake
type Mistake struct {
	ID        string          `json:"id" db:"id"`
	UserID    int             `json:"user_id" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	Category  MistakeCategory `json:"category" db:"category"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// MarketType represents the different market types
type MarketType string

const (
	MarketTypeIndian      MarketType = "indian"
	MarketTypeUS          MarketType = "us"
	MarketTypeCrypto      MarketType = "crypto"
	MarketTypeForex       MarketType = "forex"
	MarketTypeCommodities MarketType = "commodities"
)

// TradeDirection represents trade direction
type TradeDirection string

const (
	TradeDirectionLong  TradeDirection = "long"
	TradeDirectionShort TradeDirection = "short"
)

// TradeDuration represents trade duration
type TradeDuration string

const (
	TradeDurationIntraday   TradeDuration = "intraday"
	TradeDurationSwing      TradeDuration = "swing"
	TradeDurationPositional TradeDuration = "positional"
)

// OutcomeSummary represents trade outcome
type OutcomeSummary string

const (
	OutcomeSummaryProfitable    OutcomeSummary = "profitable"
	OutcomeSummaryLoss          OutcomeSummary = "loss"
	OutcomeSummaryBreakeven     OutcomeSummary = "breakeven"
	OutcomeSummaryPartialProfit OutcomeSummary = "partial_profit"
	OutcomeSummaryPartialLoss   OutcomeSummary = "partial_loss"
)

// RuleCategory represents rule categories
type RuleCategory string

const (
	RuleCategoryEntry          RuleCategory = "entry"
	RuleCategoryExit           RuleCategory = "exit"
	RuleCategoryStopLoss       RuleCategory = "stop_loss"
	RuleCategoryTakeProfit     RuleCategory = "take_profit"
	RuleCategoryRiskManagement RuleCategory = "risk_management"
	RuleCategoryPsychology     RuleCategory = "psychology"
	RuleCategoryOther          RuleCategory = "other"
)

// MistakeCategory represents mistake categories
type MistakeCategory string

const (
	MistakeCategoryPsychological MistakeCategory = "psychological"
	MistakeCategoryBehavioral    MistakeCategory = "behavioral"
)
