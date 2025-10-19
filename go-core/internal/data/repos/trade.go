package repos

import (
	"database/sql"
	"fmt"
	"time"

	"go-core/internal/data"
)

// TradeRepository handles trade-related database operations
type TradeRepository struct {
	db *sql.DB
}

// NewTradeRepository creates a new trade repository
func NewTradeRepository(db *sql.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

// Create creates a new trade
func (r *TradeRepository) Create(trade *data.Trade) error {
	query := `
		INSERT INTO trades (id, user_id, market, symbol, target, stoploss, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	trade.CreatedAt = now

	_, err := r.db.Exec(query, trade.ID, trade.UserID, trade.Market, trade.Symbol, trade.Target, trade.StopLoss, trade.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create trade: %w", err)
	}

	return nil
}

// GetByID retrieves a trade by ID
func (r *TradeRepository) GetByID(id string) (*data.Trade, error) {
	query := `
		SELECT id, user_id, market, symbol, target, stoploss, created_at
		FROM trades WHERE id = ?
	`

	trade := &data.Trade{}
	err := r.db.QueryRow(query, id).Scan(
		&trade.ID, &trade.UserID, &trade.Market, &trade.Symbol,
		&trade.Target, &trade.StopLoss, &trade.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trade not found")
		}
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}

	return trade, nil
}

// GetByUserID retrieves all trades for a user
func (r *TradeRepository) GetByUserID(userID int) ([]*data.Trade, error) {
	query := `
		SELECT id, user_id, market, symbol, target, stoploss, created_at
		FROM trades WHERE user_id = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades: %w", err)
	}
	defer rows.Close()

	var trades []*data.Trade
	for rows.Next() {
		trade := &data.Trade{}
		err := rows.Scan(
			&trade.ID, &trade.UserID, &trade.Market, &trade.Symbol,
			&trade.Target, &trade.StopLoss, &trade.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}
		trades = append(trades, trade)
	}

	return trades, nil
}

// GetByMarket retrieves trades by market type
func (r *TradeRepository) GetByMarket(market string) ([]*data.Trade, error) {
	query := `
		SELECT id, user_id, market, symbol, target, stoploss, created_at
		FROM trades WHERE market = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, market)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades by market: %w", err)
	}
	defer rows.Close()

	var trades []*data.Trade
	for rows.Next() {
		trade := &data.Trade{}
		err := rows.Scan(
			&trade.ID, &trade.UserID, &trade.Market, &trade.Symbol,
			&trade.Target, &trade.StopLoss, &trade.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}
		trades = append(trades, trade)
	}

	return trades, nil
}

// Update updates a trade
func (r *TradeRepository) Update(trade *data.Trade) error {
	query := `
		UPDATE trades 
		SET market = ?, symbol = ?, target = ?, stoploss = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, trade.Market, trade.Symbol, trade.Target, trade.StopLoss, trade.ID)
	if err != nil {
		return fmt.Errorf("failed to update trade: %w", err)
	}

	return nil
}

// Delete deletes a trade
func (r *TradeRepository) Delete(id string) error {
	query := `DELETE FROM trades WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trade: %w", err)
	}

	return nil
}

// List retrieves all trades
func (r *TradeRepository) List() ([]*data.Trade, error) {
	query := `
		SELECT id, user_id, market, symbol, target, stoploss, created_at
		FROM trades ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list trades: %w", err)
	}
	defer rows.Close()

	var trades []*data.Trade
	for rows.Next() {
		trade := &data.Trade{}
		err := rows.Scan(
			&trade.ID, &trade.UserID, &trade.Market, &trade.Symbol,
			&trade.Target, &trade.StopLoss, &trade.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}
		trades = append(trades, trade)
	}

	return trades, nil
}
