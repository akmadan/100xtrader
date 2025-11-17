package repos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// AlgorithmRepository handles algorithm database operations
type AlgorithmRepository struct {
	db *sql.DB
}

// NewAlgorithmRepository creates a new algorithm repository
func NewAlgorithmRepository(db *sql.DB) *AlgorithmRepository {
	return &AlgorithmRepository{db: db}
}

// CreateAlgorithm creates a new algorithm
func (r *AlgorithmRepository) CreateAlgorithm(algo *data.Algorithm) error {
	// Convert JSON fields
	configJSON, err := json.Marshal(algo.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if len(configJSON) == 0 {
		configJSON = []byte("{}")
	}

	stateJSON, err := json.Marshal(algo.State)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}
	if len(stateJSON) == 0 {
		stateJSON = []byte("{}")
	}

	tagsJSON, err := json.Marshal(algo.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}
	if len(tagsJSON) == 0 {
		tagsJSON = []byte("[]")
	}

	var brokerValue interface{}
	if algo.Broker != nil {
		brokerValue = string(*algo.Broker)
	}

	var lastRunAtValue interface{}
	if algo.LastRunAt != nil {
		lastRunAtValue = algo.LastRunAt.Format(time.RFC3339)
	}

	query := `
		INSERT INTO algorithms (
			id, user_id, name, description, code, status,
			symbol, timeframe, execution_mode, broker, enabled,
			config, state, last_run_at, last_signal,
			total_trades, win_rate, total_pnl, version, tags,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(
		query,
		algo.ID,
		algo.UserID,
		algo.Name,
		algo.Description,
		algo.Code,
		string(algo.Status),
		algo.Symbol,
		string(algo.Timeframe),
		string(algo.ExecutionMode),
		brokerValue,
		algo.Enabled,
		string(configJSON),
		string(stateJSON),
		lastRunAtValue,
		algo.LastSignal,
		algo.TotalTrades,
		algo.WinRate,
		algo.TotalPnL,
		algo.Version,
		string(tagsJSON),
		algo.CreatedAt.Format(time.RFC3339),
		algo.UpdatedAt.Format(time.RFC3339),
	)

	if err != nil {
		utils.LogError(err, "Failed to create algorithm")
		return fmt.Errorf("failed to create algorithm: %w", err)
	}

	return nil
}

// GetAlgorithmByID retrieves an algorithm by ID
func (r *AlgorithmRepository) GetAlgorithmByID(id string, userID int) (*data.Algorithm, error) {
	query := `
		SELECT id, user_id, name, description, code, status,
		       symbol, timeframe, execution_mode, broker, enabled,
		       config, state, last_run_at, last_signal,
		       total_trades, win_rate, total_pnl, version, tags,
		       created_at, updated_at
		FROM algorithms
		WHERE id = ? AND user_id = ?
	`

	var algo data.Algorithm
	var statusStr, timeframeStr, executionModeStr, brokerStr sql.NullString
	var description, lastSignal sql.NullString
	var lastRunAtStr sql.NullString
	var configJSON, stateJSON, tagsJSON string
	var createdAtStr, updatedAtStr string

	err := r.db.QueryRow(query, id, userID).Scan(
		&algo.ID,
		&algo.UserID,
		&algo.Name,
		&description,
		&algo.Code,
		&statusStr,
		&algo.Symbol,
		&timeframeStr,
		&executionModeStr,
		&brokerStr,
		&algo.Enabled,
		&configJSON,
		&stateJSON,
		&lastRunAtStr,
		&lastSignal,
		&algo.TotalTrades,
		&algo.WinRate,
		&algo.TotalPnL,
		&algo.Version,
		&tagsJSON,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("algorithm not found")
		}
		utils.LogError(err, "Failed to get algorithm")
		return nil, fmt.Errorf("failed to get algorithm: %w", err)
	}

	// Parse enums
	algo.Status = data.AlgorithmStatus(statusStr.String)
	algo.Timeframe = data.Timeframe(timeframeStr.String)
	algo.ExecutionMode = data.ExecutionMode(executionModeStr.String)

	if brokerStr.Valid {
		broker := data.TradingBroker(brokerStr.String)
		algo.Broker = &broker
	}

	if description.Valid {
		algo.Description = &description.String
	}

	if lastSignal.Valid {
		algo.LastSignal = &lastSignal.String
	}

	// Parse JSON fields
	if err := json.Unmarshal([]byte(configJSON), &algo.Config); err != nil {
		algo.Config = make(map[string]interface{})
	}

	if err := json.Unmarshal([]byte(stateJSON), &algo.State); err != nil {
		algo.State = make(map[string]interface{})
	}

	if err := json.Unmarshal([]byte(tagsJSON), &algo.Tags); err != nil {
		algo.Tags = []string{}
	}

	// Parse timestamps
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}
	algo.CreatedAt = createdAt

	updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}
	algo.UpdatedAt = updatedAt

	if lastRunAtStr.Valid && lastRunAtStr.String != "" {
		lastRunAt, err := time.Parse(time.RFC3339, lastRunAtStr.String)
		if err == nil {
			algo.LastRunAt = &lastRunAt
		}
	}

	return &algo, nil
}

// GetAlgorithmsByUser retrieves all algorithms for a user
func (r *AlgorithmRepository) GetAlgorithmsByUser(userID int, limit int, offset int) ([]*data.Algorithm, error) {
	query := `
		SELECT id, user_id, name, description, code, status,
		       symbol, timeframe, execution_mode, broker, enabled,
		       config, state, last_run_at, last_signal,
		       total_trades, win_rate, total_pnl, version, tags,
		       created_at, updated_at
		FROM algorithms
		WHERE user_id = ?
		ORDER BY updated_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		utils.LogError(err, "Failed to get algorithms")
		return nil, fmt.Errorf("failed to get algorithms: %w", err)
	}
	defer rows.Close()

	var algorithms []*data.Algorithm
	for rows.Next() {
		algo, err := r.scanAlgorithm(rows)
		if err != nil {
			utils.LogError(err, "Failed to scan algorithm")
			continue
		}
		algorithms = append(algorithms, algo)
	}

	return algorithms, nil
}

// UpdateAlgorithm updates an existing algorithm
func (r *AlgorithmRepository) UpdateAlgorithm(algo *data.Algorithm) error {
	// Convert JSON fields
	configJSON, err := json.Marshal(algo.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if len(configJSON) == 0 {
		configJSON = []byte("{}")
	}

	stateJSON, err := json.Marshal(algo.State)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}
	if len(stateJSON) == 0 {
		stateJSON = []byte("{}")
	}

	tagsJSON, err := json.Marshal(algo.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}
	if len(tagsJSON) == 0 {
		tagsJSON = []byte("[]")
	}

	var brokerValue interface{}
	if algo.Broker != nil {
		brokerValue = string(*algo.Broker)
	}

	var lastRunAtValue interface{}
	if algo.LastRunAt != nil {
		lastRunAtValue = algo.LastRunAt.Format(time.RFC3339)
	}

	query := `
		UPDATE algorithms
		SET name = ?, description = ?, code = ?, status = ?,
		    symbol = ?, timeframe = ?, execution_mode = ?, broker = ?, enabled = ?,
		    config = ?, state = ?, last_run_at = ?, last_signal = ?,
		    total_trades = ?, win_rate = ?, total_pnl = ?, version = ?, tags = ?,
		    updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(
		query,
		algo.Name,
		algo.Description,
		algo.Code,
		string(algo.Status),
		algo.Symbol,
		string(algo.Timeframe),
		string(algo.ExecutionMode),
		brokerValue,
		algo.Enabled,
		string(configJSON),
		string(stateJSON),
		lastRunAtValue,
		algo.LastSignal,
		algo.TotalTrades,
		algo.WinRate,
		algo.TotalPnL,
		algo.Version,
		string(tagsJSON),
		algo.UpdatedAt.Format(time.RFC3339),
		algo.ID,
		algo.UserID,
	)

	if err != nil {
		utils.LogError(err, "Failed to update algorithm")
		return fmt.Errorf("failed to update algorithm: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("algorithm not found or user mismatch")
	}

	return nil
}

// DeleteAlgorithm deletes an algorithm
func (r *AlgorithmRepository) DeleteAlgorithm(id string, userID int) error {
	query := `DELETE FROM algorithms WHERE id = ? AND user_id = ?`

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		utils.LogError(err, "Failed to delete algorithm")
		return fmt.Errorf("failed to delete algorithm: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("algorithm not found or user mismatch")
	}

	return nil
}

// scanAlgorithm scans a row into an Algorithm struct
func (r *AlgorithmRepository) scanAlgorithm(rows *sql.Rows) (*data.Algorithm, error) {
	var algo data.Algorithm
	var statusStr, timeframeStr, executionModeStr, brokerStr sql.NullString
	var description, lastSignal sql.NullString
	var lastRunAtStr sql.NullString
	var configJSON, stateJSON, tagsJSON string
	var createdAtStr, updatedAtStr string

	err := rows.Scan(
		&algo.ID,
		&algo.UserID,
		&algo.Name,
		&description,
		&algo.Code,
		&statusStr,
		&algo.Symbol,
		&timeframeStr,
		&executionModeStr,
		&brokerStr,
		&algo.Enabled,
		&configJSON,
		&stateJSON,
		&lastRunAtStr,
		&lastSignal,
		&algo.TotalTrades,
		&algo.WinRate,
		&algo.TotalPnL,
		&algo.Version,
		&tagsJSON,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		return nil, err
	}

	// Parse enums
	algo.Status = data.AlgorithmStatus(statusStr.String)
	algo.Timeframe = data.Timeframe(timeframeStr.String)
	algo.ExecutionMode = data.ExecutionMode(executionModeStr.String)

	if brokerStr.Valid {
		broker := data.TradingBroker(brokerStr.String)
		algo.Broker = &broker
	}

	if description.Valid {
		algo.Description = &description.String
	}

	if lastSignal.Valid {
		algo.LastSignal = &lastSignal.String
	}

	// Parse JSON fields
	if err := json.Unmarshal([]byte(configJSON), &algo.Config); err != nil {
		algo.Config = make(map[string]interface{})
	}

	if err := json.Unmarshal([]byte(stateJSON), &algo.State); err != nil {
		algo.State = make(map[string]interface{})
	}

	if err := json.Unmarshal([]byte(tagsJSON), &algo.Tags); err != nil {
		algo.Tags = []string{}
	}

	// Parse timestamps
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}
	algo.CreatedAt = createdAt

	updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}
	algo.UpdatedAt = updatedAt

	if lastRunAtStr.Valid && lastRunAtStr.String != "" {
		lastRunAt, err := time.Parse(time.RFC3339, lastRunAtStr.String)
		if err == nil {
			algo.LastRunAt = &lastRunAt
		}
	}

	return &algo, nil
}

