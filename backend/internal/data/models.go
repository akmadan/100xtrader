package data

import "time"

// Order represents a trading order
type Order struct {
	ID        int64     `db:"id"`
	User      string    `db:"user"`
	Symbol    string    `db:"symbol"`
	Side      string    `db:"side"` // buy or sell
	Type      string    `db:"type"` // market, limit, stop
	Quantity  float64   `db:"quantity"`
	Price     float64   `db:"price"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Trade represents a matched trade
type Trade struct {
	ID          int64     `db:"id"`
	BuyOrderID  int64     `db:"buy_order_id"`
	SellOrderID int64     `db:"sell_order_id"`
	Symbol      string    `db:"symbol"`
	Quantity    float64   `db:"quantity"`
	Price       float64   `db:"price"`
	Timestamp   time.Time `db:"timestamp"`
}

// Position represents a user's position in a symbol
type Position struct {
	ID           int64   `db:"id"`
	User         string  `db:"user"`
	Symbol       string  `db:"symbol"`
	Quantity     float64 `db:"quantity"`
	AveragePrice float64 `db:"average_price"`
}

// Environment represents a trading environment definition
type Environment struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}
