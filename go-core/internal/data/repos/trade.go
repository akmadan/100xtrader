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
			trading_broker, trader_broker_id, exchange_order_id, order_id, product_type, transaction_type,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var tradingBroker, traderBrokerID, exchangeOrderID, orderID, productType, transactionType interface{}
	if trade.TradingBroker != nil {
		tradingBroker = string(*trade.TradingBroker)
	}
	if trade.TraderBrokerID != nil {
		traderBrokerID = *trade.TraderBrokerID
	}
	if trade.ExchangeOrderID != nil {
		exchangeOrderID = *trade.ExchangeOrderID
	}
	if trade.OrderID != nil {
		orderID = *trade.OrderID
	}
	if trade.ProductType != nil {
		productType = string(*trade.ProductType)
	}
	if trade.TransactionType != nil {
		transactionType = *trade.TransactionType
	}

	_, err = r.db.Exec(query,
		trade.ID, trade.UserID, trade.Symbol, trade.MarketType, trade.EntryDate,
		trade.EntryPrice, trade.Quantity, trade.TotalAmount, trade.ExitPrice,
		trade.Direction, trade.StopLoss, trade.Target, trade.Strategy,
		trade.OutcomeSummary, trade.TradeAnalysis, string(rulesFollowedJSON),
		string(screenshotsJSON), string(psychologyJSON),
		tradingBroker, traderBrokerID, exchangeOrderID, orderID, productType, transactionType,
		trade.CreatedAt, trade.UpdatedAt,
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
			psychology = ?, trading_broker = ?, trader_broker_id = ?, 
			exchange_order_id = ?, order_id = ?, product_type = ?, transaction_type = ?,
			updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	var tradingBroker, traderBrokerID, exchangeOrderID, orderID, productType, transactionType interface{}
	if trade.TradingBroker != nil {
		tradingBroker = string(*trade.TradingBroker)
	}
	if trade.TraderBrokerID != nil {
		traderBrokerID = *trade.TraderBrokerID
	}
	if trade.ExchangeOrderID != nil {
		exchangeOrderID = *trade.ExchangeOrderID
	}
	if trade.OrderID != nil {
		orderID = *trade.OrderID
	}
	if trade.ProductType != nil {
		productType = string(*trade.ProductType)
	}
	if trade.TransactionType != nil {
		transactionType = *trade.TransactionType
	}

	result, err := r.db.Exec(query,
		trade.Symbol, trade.MarketType, trade.EntryDate, trade.EntryPrice,
		trade.Quantity, trade.TotalAmount, trade.ExitPrice, trade.Direction,
		trade.StopLoss, trade.Target, trade.Strategy, trade.OutcomeSummary,
		trade.TradeAnalysis, string(rulesFollowedJSON), string(screenshotsJSON),
		string(psychologyJSON),
		tradingBroker, traderBrokerID, exchangeOrderID, orderID, productType, transactionType,
		trade.UpdatedAt, trade.ID, trade.UserID,
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
			   trading_broker, trader_broker_id, exchange_order_id, order_id, product_type, transaction_type,
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
			   trading_broker, trader_broker_id, exchange_order_id, order_id, product_type, transaction_type,
			   created_at, updated_at
		FROM trades 
		WHERE user_id = ?
		ORDER BY entry_date DESC, created_at DESC
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
	var tradingBroker, traderBrokerID, exchangeOrderID, orderID, productType, transactionType sql.NullString

	err := scanner.Scan(
		&trade.ID, &trade.UserID, &trade.Symbol, &trade.MarketType, &entryDate,
		&trade.EntryPrice, &trade.Quantity, &trade.TotalAmount, &trade.ExitPrice,
		&trade.Direction, &trade.StopLoss, &trade.Target, &trade.Strategy,
		&trade.OutcomeSummary, &trade.TradeAnalysis, &rulesFollowedJSON,
		&screenshotsJSON, &psychologyJSON,
		&tradingBroker, &traderBrokerID, &exchangeOrderID, &orderID, &productType, &transactionType,
		&createdAt, &updatedAt,
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

	// Set broker-specific fields
	if tradingBroker.Valid {
		broker := data.TradingBroker(tradingBroker.String)
		trade.TradingBroker = &broker
	}
	if traderBrokerID.Valid {
		trade.TraderBrokerID = &traderBrokerID.String
	}
	if exchangeOrderID.Valid {
		trade.ExchangeOrderID = &exchangeOrderID.String
	}
	if orderID.Valid {
		trade.OrderID = &orderID.String
	}
	if productType.Valid {
		product := data.ProductType(productType.String)
		trade.ProductType = &product
	}
	if transactionType.Valid {
		trade.TransactionType = &transactionType.String
	}

	// Set time fields
	trade.EntryDate = entryDate
	trade.CreatedAt = createdAt
	trade.UpdatedAt = updatedAt

	return &trade, nil
}

// TradeExistsByBrokerID checks if a trade already exists for a user by broker-specific ID
// This is used to prevent duplicate imports from broker APIs
func (r *TradeRepository) TradeExistsByBrokerID(userID int, tradingBroker data.TradingBroker, exchangeOrderID, orderID string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM trades 
		WHERE user_id = ? 
		AND trading_broker = ? 
		AND (
			(exchange_order_id IS NOT NULL AND exchange_order_id = ?) OR
			(order_id IS NOT NULL AND order_id = ?)
		)
		LIMIT 1
	`

	var count int
	err := r.db.QueryRow(query, userID, string(tradingBroker), exchangeOrderID, orderID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		utils.LogError(err, "Failed to check if trade exists by broker ID", map[string]interface{}{
			"user_id":           userID,
			"trading_broker":    tradingBroker,
			"exchange_order_id": exchangeOrderID,
			"order_id":          orderID,
		})
		return false, fmt.Errorf("failed to check if trade exists: %w", err)
	}

	return count > 0, nil
}

// GetLatestTradeDateByBroker gets the latest entry_date for trades from a specific broker for a user
// Returns nil if no trades exist for that broker
func (r *TradeRepository) GetLatestTradeDateByBroker(userID int, tradingBroker data.TradingBroker) (*time.Time, error) {
	query := `
		SELECT MAX(entry_date) 
		FROM trades 
		WHERE user_id = ? AND trading_broker = ?
	`

	var latestDateStr sql.NullString
	err := r.db.QueryRow(query, userID, string(tradingBroker)).Scan(&latestDateStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No trades found
		}
		utils.LogError(err, "Failed to get latest trade date by broker", map[string]interface{}{
			"user_id":        userID,
			"trading_broker": tradingBroker,
		})
		return nil, fmt.Errorf("failed to get latest trade date: %w", err)
	}

	if !latestDateStr.Valid || latestDateStr.String == "" {
		return nil, nil // No trades found
	}

	// Parse the date string (SQLite stores dates as strings)
	latestDate, err := time.Parse("2006-01-02T15:04:05Z07:00", latestDateStr.String)
	if err != nil {
		// Try alternative formats
		latestDate, err = time.Parse("2006-01-02 15:04:05", latestDateStr.String)
		if err != nil {
			latestDate, err = time.Parse("2006-01-02", latestDateStr.String)
			if err != nil {
				utils.LogError(err, "Failed to parse latest trade date", map[string]interface{}{
					"date_string": latestDateStr.String,
				})
				return nil, fmt.Errorf("failed to parse latest trade date: %w", err)
			}
		}
	}

	return &latestDate, nil
}
