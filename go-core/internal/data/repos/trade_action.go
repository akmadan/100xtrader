package repos

import (
	"database/sql"
	"fmt"
	"time"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// TradeActionRepository handles trade action-related database operations
type TradeActionRepository struct {
	db *sql.DB
}

// NewTradeActionRepository creates a new trade action repository
func NewTradeActionRepository(db *sql.DB) *TradeActionRepository {
	return &TradeActionRepository{db: db}
}

// Create creates a new trade action
func (r *TradeActionRepository) Create(action *data.TradeAction) error {
	start := time.Now()
	query := `
		INSERT INTO trade_actions (trade_id, action, trade_time, quantity, price, fee, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	action.CreatedAt = now

	result, err := r.db.Exec(query, action.TradeID, action.Action, action.TradeTime, action.Quantity, action.Price, action.Fee, action.CreatedAt)
	duration := time.Since(start)

	if err != nil {
		utils.LogDatabase("CREATE", "trade_actions", duration, err, map[string]interface{}{
			"trade_id": action.TradeID,
			"action":   action.Action,
		})
		return fmt.Errorf("failed to create trade action: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.LogError(err, "Failed to get trade action ID", map[string]interface{}{
			"trade_id": action.TradeID,
		})
		return fmt.Errorf("failed to get trade action ID: %w", err)
	}

	action.ID = int(id)
	utils.LogDatabase("CREATE", "trade_actions", duration, nil, map[string]interface{}{
		"action_id": action.ID,
		"trade_id":  action.TradeID,
	})

	return nil
}

// GetByID retrieves a trade action by ID
func (r *TradeActionRepository) GetByID(id int) (*data.TradeAction, error) {
	query := `
		SELECT id, trade_id, action, trade_time, quantity, price, fee, created_at
		FROM trade_actions WHERE id = ?
	`

	action := &data.TradeAction{}
	err := r.db.QueryRow(query, id).Scan(
		&action.ID, &action.TradeID, &action.Action, &action.TradeTime,
		&action.Quantity, &action.Price, &action.Fee, &action.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trade action not found")
		}
		return nil, fmt.Errorf("failed to get trade action: %w", err)
	}

	return action, nil
}

// GetByTradeID retrieves all actions for a specific trade
func (r *TradeActionRepository) GetByTradeID(tradeID string) ([]*data.TradeAction, error) {
	query := `
		SELECT id, trade_id, action, trade_time, quantity, price, fee, created_at
		FROM trade_actions WHERE trade_id = ? ORDER BY trade_time ASC
	`

	rows, err := r.db.Query(query, tradeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade actions: %w", err)
	}
	defer rows.Close()

	var actions []*data.TradeAction
	for rows.Next() {
		action := &data.TradeAction{}
		err := rows.Scan(
			&action.ID, &action.TradeID, &action.Action, &action.TradeTime,
			&action.Quantity, &action.Price, &action.Fee, &action.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade action: %w", err)
		}
		actions = append(actions, action)
	}

	return actions, nil
}

// Delete deletes a trade action
func (r *TradeActionRepository) Delete(id int) error {
	query := `DELETE FROM trade_actions WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trade action: %w", err)
	}

	return nil
}
