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
	ID        string    `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Market    string    `json:"market" db:"market"` // stock, option, crypto, futures, forex, index
	Symbol    string    `json:"symbol" db:"symbol"`
	Target    float64   `json:"target" db:"target"`
	StopLoss  float64   `json:"stoploss" db:"stoploss"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TradeAction represents a buy or sell action for a trade
type TradeAction struct {
	ID        int       `json:"id" db:"id"`
	TradeID   string    `json:"trade_id" db:"trade_id"`
	Action    string    `json:"action" db:"action"` // buy, sell
	TradeTime time.Time `json:"trade_time" db:"trade_time"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`
	Fee       float64   `json:"fee" db:"fee"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TradeJournal represents journal entries for trades
type TradeJournal struct {
	ID         int       `json:"id" db:"id"`
	TradeID    string    `json:"trade_id" db:"trade_id"`
	Notes      string    `json:"notes" db:"notes"`
	Confidence int       `json:"confidence" db:"confidence"` // Scale 0-10
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Tag represents a tag for categorizing trades
type Tag struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TradeTag represents the many-to-many relationship between trades and tags
type TradeTag struct {
	TradeID   string    `json:"trade_id" db:"trade_id"`
	TagID     int       `json:"tag_id" db:"tag_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TradeScreenshot represents screenshots attached to trade journals
type TradeScreenshot struct {
	ID             int       `json:"id" db:"id"`
	TradeJournalID int       `json:"trade_journal_id" db:"trade_journal_id"`
	URL            string    `json:"url" db:"url"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// Note represents daily trading notes
type Note struct {
	ID               int       `json:"id" db:"id"`
	UserID           int       `json:"user_id" db:"user_id"`
	Mood             string    `json:"mood" db:"mood"`                           // excited, neutral, low
	MarketCondition  string    `json:"market_condition" db:"market_condition"`   // up, down, sideways
	MarketVolatility string    `json:"market_volatility" db:"market_volatility"` // high, medium, low
	Summary          string    `json:"summary" db:"summary"`
	Day              time.Time `json:"day" db:"day"`
	Notes            string    `json:"notes" db:"notes"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// TradeSetup represents planned trade setups
type TradeSetup struct {
	ID              int       `json:"id" db:"id"`
	UserID          int       `json:"user_id" db:"user_id"`
	Market          string    `json:"market" db:"market"` // stock, option, crypto, futures, forex, index
	Side            string    `json:"side" db:"side"`     // long, short
	Symbol          string    `json:"symbol" db:"symbol"`
	Entry           float64   `json:"entry" db:"entry"`
	Target          float64   `json:"target" db:"target"`
	StopLoss        float64   `json:"stoploss" db:"stoploss"`
	Note            string    `json:"note" db:"note"`
	RiskRewardRatio float64   `json:"risk_reward_ratio" db:"risk_reward_ratio"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// MarketType represents the different market types
type MarketType string

const (
	MarketStock   MarketType = "stock"
	MarketOption  MarketType = "option"
	MarketCrypto  MarketType = "crypto"
	MarketFutures MarketType = "futures"
	MarketForex   MarketType = "forex"
	MarketIndex   MarketType = "index"
)

// ActionType represents trade action types
type ActionType string

const (
	ActionBuy  ActionType = "buy"
	ActionSell ActionType = "sell"
)

// MoodType represents user mood types
type MoodType string

const (
	MoodExcited MoodType = "excited"
	MoodNeutral MoodType = "neutral"
	MoodLow     MoodType = "low"
)

// MarketConditionType represents market condition types
type MarketConditionType string

const (
	ConditionUp       MarketConditionType = "up"
	ConditionDown     MarketConditionType = "down"
	ConditionSideways MarketConditionType = "sideways"
)

// VolatilityType represents market volatility types
type VolatilityType string

const (
	VolatilityHigh   VolatilityType = "high"
	VolatilityMedium VolatilityType = "medium"
	VolatilityLow    VolatilityType = "low"
)

// SideType represents trade side types
type SideType string

const (
	SideLong  SideType = "long"
	SideShort SideType = "short"
)
