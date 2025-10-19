package dto

import (
	"time"
)

// TradeSetupCreateRequest represents the request to create a new trade setup
type TradeSetupCreateRequest struct {
	UserID          int     `json:"user_id" validate:"required,min=1"`
	Market          string  `json:"market" validate:"required,oneof=stock option crypto futures forex index"`
	Side            string  `json:"side" validate:"required,oneof=long short"`
	Symbol          string  `json:"symbol" validate:"required,min=1,max=50"`
	Entry           float64 `json:"entry" validate:"required,gt=0"`
	Target          float64 `json:"target" validate:"required,gt=0"`
	StopLoss        float64 `json:"stoploss" validate:"required,gt=0"`
	Note            string  `json:"note,omitempty"`
	RiskRewardRatio float64 `json:"risk_reward_ratio,omitempty" validate:"omitempty,gte=0"`
}

// TradeSetupUpdateRequest represents the request to update a trade setup
type TradeSetupUpdateRequest struct {
	Market          *string  `json:"market,omitempty" validate:"omitempty,oneof=stock option crypto futures forex index"`
	Side            *string  `json:"side,omitempty" validate:"omitempty,oneof=long short"`
	Symbol          *string  `json:"symbol,omitempty" validate:"omitempty,min=1,max=50"`
	Entry           *float64 `json:"entry,omitempty" validate:"omitempty,gt=0"`
	Target          *float64 `json:"target,omitempty" validate:"omitempty,gt=0"`
	StopLoss        *float64 `json:"stoploss,omitempty" validate:"omitempty,gt=0"`
	Note            *string  `json:"note,omitempty"`
	RiskRewardRatio *float64 `json:"risk_reward_ratio,omitempty" validate:"omitempty,gte=0"`
}

// TradeSetupResponse represents the response for trade setup data
type TradeSetupResponse struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Market          string    `json:"market"`
	Side            string    `json:"side"`
	Symbol          string    `json:"symbol"`
	Entry           float64   `json:"entry"`
	Target          float64   `json:"target"`
	StopLoss        float64   `json:"stoploss"`
	Note            string    `json:"note"`
	RiskRewardRatio float64   `json:"risk_reward_ratio"`
	CreatedAt       time.Time `json:"created_at"`
}

// TradeSetupListResponse represents the response for listing trade setups
type TradeSetupListResponse struct {
	Setups []TradeSetupResponse `json:"setups"`
	Total  int                  `json:"total"`
}

// TradeSetupFilterRequest represents the request to filter trade setups
type TradeSetupFilterRequest struct {
	UserID *int    `json:"user_id,omitempty" validate:"omitempty,min=1"`
	Market *string `json:"market,omitempty" validate:"omitempty,oneof=stock option crypto futures forex index"`
	Side   *string `json:"side,omitempty" validate:"omitempty,oneof=long short"`
	Symbol *string `json:"symbol,omitempty"`
}

// TradeSetupStatsResponse represents statistics for trade setups
type TradeSetupStatsResponse struct {
	TotalSetups     int            `json:"total_setups"`
	LongSetups      int            `json:"long_setups"`
	ShortSetups     int            `json:"short_setups"`
	AvgRiskReward   float64        `json:"avg_risk_reward"`
	MarketBreakdown map[string]int `json:"market_breakdown"`
}
