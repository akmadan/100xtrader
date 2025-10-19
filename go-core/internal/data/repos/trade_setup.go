package repos

import (
	"database/sql"
	"fmt"
	"time"

	"go-core/internal/data"
)

// TradeSetupRepository handles trade setup-related database operations
type TradeSetupRepository struct {
	db *sql.DB
}

// NewTradeSetupRepository creates a new trade setup repository
func NewTradeSetupRepository(db *sql.DB) *TradeSetupRepository {
	return &TradeSetupRepository{db: db}
}

// Create creates a new trade setup
func (r *TradeSetupRepository) Create(setup *data.TradeSetup) error {
	query := `
		INSERT INTO trade_setups (user_id, market, side, symbol, entry, target, stoploss, note, risk_reward_ratio, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	setup.CreatedAt = now

	result, err := r.db.Exec(query, setup.UserID, setup.Market, setup.Side, setup.Symbol, setup.Entry, setup.Target, setup.StopLoss, setup.Note, setup.RiskRewardRatio, setup.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create trade setup: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get trade setup ID: %w", err)
	}

	setup.ID = int(id)
	return nil
}

// GetByID retrieves a trade setup by ID
func (r *TradeSetupRepository) GetByID(id int) (*data.TradeSetup, error) {
	query := `
		SELECT id, user_id, market, side, symbol, entry, target, stoploss, note, risk_reward_ratio, created_at
		FROM trade_setups WHERE id = ?
	`

	setup := &data.TradeSetup{}
	err := r.db.QueryRow(query, id).Scan(
		&setup.ID, &setup.UserID, &setup.Market, &setup.Side, &setup.Symbol,
		&setup.Entry, &setup.Target, &setup.StopLoss, &setup.Note, &setup.RiskRewardRatio, &setup.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trade setup not found")
		}
		return nil, fmt.Errorf("failed to get trade setup: %w", err)
	}

	return setup, nil
}

// GetByUserID retrieves all trade setups for a user
func (r *TradeSetupRepository) GetByUserID(userID int) ([]*data.TradeSetup, error) {
	query := `
		SELECT id, user_id, market, side, symbol, entry, target, stoploss, note, risk_reward_ratio, created_at
		FROM trade_setups WHERE user_id = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade setups: %w", err)
	}
	defer rows.Close()

	var setups []*data.TradeSetup
	for rows.Next() {
		setup := &data.TradeSetup{}
		err := rows.Scan(
			&setup.ID, &setup.UserID, &setup.Market, &setup.Side, &setup.Symbol,
			&setup.Entry, &setup.Target, &setup.StopLoss, &setup.Note, &setup.RiskRewardRatio, &setup.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade setup: %w", err)
		}
		setups = append(setups, setup)
	}

	return setups, nil
}

// GetByMarket retrieves trade setups by market type
func (r *TradeSetupRepository) GetByMarket(market string) ([]*data.TradeSetup, error) {
	query := `
		SELECT id, user_id, market, side, symbol, entry, target, stoploss, note, risk_reward_ratio, created_at
		FROM trade_setups WHERE market = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, market)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade setups by market: %w", err)
	}
	defer rows.Close()

	var setups []*data.TradeSetup
	for rows.Next() {
		setup := &data.TradeSetup{}
		err := rows.Scan(
			&setup.ID, &setup.UserID, &setup.Market, &setup.Side, &setup.Symbol,
			&setup.Entry, &setup.Target, &setup.StopLoss, &setup.Note, &setup.RiskRewardRatio, &setup.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade setup: %w", err)
		}
		setups = append(setups, setup)
	}

	return setups, nil
}

// Update updates a trade setup
func (r *TradeSetupRepository) Update(setup *data.TradeSetup) error {
	query := `
		UPDATE trade_setups 
		SET market = ?, side = ?, symbol = ?, entry = ?, target = ?, stoploss = ?, note = ?, risk_reward_ratio = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, setup.Market, setup.Side, setup.Symbol, setup.Entry, setup.Target, setup.StopLoss, setup.Note, setup.RiskRewardRatio, setup.ID)
	if err != nil {
		return fmt.Errorf("failed to update trade setup: %w", err)
	}

	return nil
}

// Delete deletes a trade setup
func (r *TradeSetupRepository) Delete(id int) error {
	query := `DELETE FROM trade_setups WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trade setup: %w", err)
	}

	return nil
}

// List retrieves all trade setups
func (r *TradeSetupRepository) List() ([]*data.TradeSetup, error) {
	query := `
		SELECT id, user_id, market, side, symbol, entry, target, stoploss, note, risk_reward_ratio, created_at
		FROM trade_setups ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list trade setups: %w", err)
	}
	defer rows.Close()

	var setups []*data.TradeSetup
	for rows.Next() {
		setup := &data.TradeSetup{}
		err := rows.Scan(
			&setup.ID, &setup.UserID, &setup.Market, &setup.Side, &setup.Symbol,
			&setup.Entry, &setup.Target, &setup.StopLoss, &setup.Note, &setup.RiskRewardRatio, &setup.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade setup: %w", err)
		}
		setups = append(setups, setup)
	}

	return setups, nil
}
