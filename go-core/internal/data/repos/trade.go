package repos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// TradeRepository handles trade database operations
type TradeRepository struct {
	db *sql.DB
}

// NewTradeRepository creates a new trade repository
func NewTradeRepository(db *sql.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

// CreateTrade creates a new trade
func (r *TradeRepository) CreateTrade(trade *data.Trade) error {
	// Convert slices to JSON
	rulesFollowedJSON, err := json.Marshal(trade.RulesFollowed)
	if err != nil {
		return fmt.Errorf("failed to marshal rules_followed: %w", err)
	}

	screenshotsJSON, err := json.Marshal(trade.Screenshots)
	if err != nil {
		return fmt.Errorf("failed to marshal screenshots: %w", err)
	}

	var psychologyJSON []byte
	if trade.Psychology != nil {
		psychologyJSON, err = json.Marshal(trade.Psychology)
		if err != nil {
			return fmt.Errorf("failed to marshal psychology: %w", err)
		}
	}

	query := `
		INSERT INTO trades (
			id, user_id, symbol, market_type, entry_date, entry_price, quantity, 
			total_amount, exit_price, direction, stop_loss, target, strategy, 
			outcome_summary, trade_analysis, rules_followed, screenshots, psychology,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(query,
		trade.ID, trade.UserID, trade.Symbol, trade.MarketType, trade.EntryDate,
		trade.EntryPrice, trade.Quantity, trade.TotalAmount, trade.ExitPrice,
		trade.Direction, trade.StopLoss, trade.Target, trade.Strategy,
		trade.OutcomeSummary, trade.TradeAnalysis, string(rulesFollowedJSON),
		string(screenshotsJSON), string(psychologyJSON), trade.CreatedAt, trade.UpdatedAt,
	)

	if err != nil {
		utils.LogError(err, "Failed to create trade", map[string]interface{}{
			"trade_id": trade.ID,
			"user_id":  trade.UserID,
		})
		return fmt.Errorf("failed to create trade: %w", err)
	}

	utils.LogInfo("Trade created successfully", map[string]interface{}{
		"trade_id": trade.ID,
		"user_id":  trade.UserID,
	})
	return nil
}

// UpdateTrade updates an existing trade
func (r *TradeRepository) UpdateTrade(trade *data.Trade) error {
	// Convert slices to JSON
	rulesFollowedJSON, err := json.Marshal(trade.RulesFollowed)
	if err != nil {
		return fmt.Errorf("failed to marshal rules_followed: %w", err)
	}

	screenshotsJSON, err := json.Marshal(trade.Screenshots)
	if err != nil {
		return fmt.Errorf("failed to marshal screenshots: %w", err)
	}

	var psychologyJSON []byte
	if trade.Psychology != nil {
		psychologyJSON, err = json.Marshal(trade.Psychology)
		if err != nil {
			return fmt.Errorf("failed to marshal psychology: %w", err)
		}
	}

	query := `
		UPDATE trades SET 
			symbol = ?, market_type = ?, entry_date = ?, entry_price = ?, 
			quantity = ?, total_amount = ?, exit_price = ?, direction = ?, 
			stop_loss = ?, target = ?, strategy = ?, outcome_summary = ?, 
			trade_analysis = ?, rules_followed = ?, screenshots = ?, 
			psychology = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(query,
		trade.Symbol, trade.MarketType, trade.EntryDate, trade.EntryPrice,
		trade.Quantity, trade.TotalAmount, trade.ExitPrice, trade.Direction,
		trade.StopLoss, trade.Target, trade.Strategy, trade.OutcomeSummary,
		trade.TradeAnalysis, string(rulesFollowedJSON), string(screenshotsJSON),
		string(psychologyJSON), trade.UpdatedAt, trade.ID, trade.UserID,
	)

	if err != nil {
		utils.LogError(err, "Failed to update trade", map[string]interface{}{
			"trade_id": trade.ID,
			"user_id":  trade.UserID,
		})
		return fmt.Errorf("failed to update trade: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("trade not found or not owned by user")
	}

	utils.LogInfo("Trade updated successfully", map[string]interface{}{
		"trade_id": trade.ID,
		"user_id":  trade.UserID,
	})
	return nil
}

// GetTradeByID retrieves a trade by ID
func (r *TradeRepository) GetTradeByID(tradeID string, userID int) (*data.Trade, error) {
	query := `
		SELECT id, user_id, symbol, market_type, entry_date, entry_price, quantity,
			   total_amount, exit_price, direction, stop_loss, target, strategy,
			   outcome_summary, trade_analysis, rules_followed, screenshots, psychology,
			   created_at, updated_at
		FROM trades 
		WHERE id = ? AND user_id = ?
	`

	row := r.db.QueryRow(query, tradeID, userID)
	trade, err := r.scanTrade(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trade not found")
		}
		utils.LogError(err, "Failed to get trade by ID", map[string]interface{}{
			"trade_id": tradeID,
			"user_id":  userID,
		})
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}

	return trade, nil
}

// GetTradesByUser retrieves all trades for a user
func (r *TradeRepository) GetTradesByUser(userID int, limit, offset int) ([]*data.Trade, error) {
	query := `
		SELECT id, user_id, symbol, market_type, entry_date, entry_price, quantity,
			   total_amount, exit_price, direction, stop_loss, target, strategy,
			   outcome_summary, trade_analysis, rules_followed, screenshots, psychology,
			   created_at, updated_at
		FROM trades 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		utils.LogError(err, "Failed to get trades by user", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get trades: %w", err)
	}
	defer rows.Close()

	var trades []*data.Trade
	for rows.Next() {
		trade, err := r.scanTrade(rows)
		if err != nil {
			utils.LogError(err, "Failed to scan trade", map[string]interface{}{
				"user_id": userID,
			})
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}
		trades = append(trades, trade)
	}

	return trades, nil
}

// DeleteTrade deletes a trade
func (r *TradeRepository) DeleteTrade(tradeID string, userID int) error {
	query := "DELETE FROM trades WHERE id = ? AND user_id = ?"

	result, err := r.db.Exec(query, tradeID, userID)
	if err != nil {
		utils.LogError(err, "Failed to delete trade", map[string]interface{}{
			"trade_id": tradeID,
			"user_id":  userID,
		})
		return fmt.Errorf("failed to delete trade: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("trade not found or not owned by user")
	}

	utils.LogInfo("Trade deleted successfully", map[string]interface{}{
		"trade_id": tradeID,
		"user_id":  userID,
	})
	return nil
}

// scanTrade scans a database row into a Trade struct
func (r *TradeRepository) scanTrade(scanner interface {
	Scan(dest ...interface{}) error
}) (*data.Trade, error) {
	var trade data.Trade
	var rulesFollowedJSON, screenshotsJSON, psychologyJSON string
	var entryDate, createdAt, updatedAt time.Time

	err := scanner.Scan(
		&trade.ID, &trade.UserID, &trade.Symbol, &trade.MarketType, &entryDate,
		&trade.EntryPrice, &trade.Quantity, &trade.TotalAmount, &trade.ExitPrice,
		&trade.Direction, &trade.StopLoss, &trade.Target, &trade.Strategy,
		&trade.OutcomeSummary, &trade.TradeAnalysis, &rulesFollowedJSON,
		&screenshotsJSON, &psychologyJSON, &createdAt, &updatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Convert JSON strings back to slices/structs
	if err := json.Unmarshal([]byte(rulesFollowedJSON), &trade.RulesFollowed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rules_followed: %w", err)
	}

	if err := json.Unmarshal([]byte(screenshotsJSON), &trade.Screenshots); err != nil {
		return nil, fmt.Errorf("failed to unmarshal screenshots: %w", err)
	}

	if psychologyJSON != "" {
		var psychology data.TradePsychology
		if err := json.Unmarshal([]byte(psychologyJSON), &psychology); err != nil {
			return nil, fmt.Errorf("failed to unmarshal psychology: %w", err)
		}
		trade.Psychology = &psychology
	}

	// Set time fields
	trade.EntryDate = entryDate
	trade.CreatedAt = createdAt
	trade.UpdatedAt = updatedAt

	return &trade, nil
}
