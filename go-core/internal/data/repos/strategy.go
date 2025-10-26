package repos

import (
	"database/sql"
	"fmt"
	"time"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// StrategyRepository handles strategy database operations
type StrategyRepository struct {
	db *sql.DB
}

// NewStrategyRepository creates a new strategy repository
func NewStrategyRepository(db *sql.DB) *StrategyRepository {
	return &StrategyRepository{db: db}
}

// CreateStrategy creates a new strategy
func (r *StrategyRepository) CreateStrategy(strategy *data.Strategy) error {
	query := `
		INSERT INTO strategies (id, user_id, name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		strategy.ID, strategy.UserID, strategy.Name, strategy.Description,
		strategy.CreatedAt, strategy.UpdatedAt,
	)

	if err != nil {
		utils.LogError(err, "Failed to create strategy", map[string]interface{}{
			"strategy_id": strategy.ID,
			"user_id":     strategy.UserID,
		})
		return fmt.Errorf("failed to create strategy: %w", err)
	}

	utils.LogInfo("Strategy created successfully", map[string]interface{}{
		"strategy_id": strategy.ID,
		"user_id":     strategy.UserID,
	})
	return nil
}

// UpdateStrategy updates an existing strategy
func (r *StrategyRepository) UpdateStrategy(strategy *data.Strategy) error {
	query := `
		UPDATE strategies SET 
			name = ?, description = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(query,
		strategy.Name, strategy.Description, strategy.UpdatedAt,
		strategy.ID, strategy.UserID,
	)

	if err != nil {
		utils.LogError(err, "Failed to update strategy", map[string]interface{}{
			"strategy_id": strategy.ID,
			"user_id":     strategy.UserID,
		})
		return fmt.Errorf("failed to update strategy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("strategy not found or not owned by user")
	}

	utils.LogInfo("Strategy updated successfully", map[string]interface{}{
		"strategy_id": strategy.ID,
		"user_id":     strategy.UserID,
	})
	return nil
}

// GetStrategyByID retrieves a strategy by ID
func (r *StrategyRepository) GetStrategyByID(strategyID string, userID int) (*data.Strategy, error) {
	query := `
		SELECT id, user_id, name, description, created_at, updated_at
		FROM strategies 
		WHERE id = ? AND user_id = ?
	`

	row := r.db.QueryRow(query, strategyID, userID)
	strategy, err := r.scanStrategy(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("strategy not found")
		}
		utils.LogError(err, "Failed to get strategy by ID", map[string]interface{}{
			"strategy_id": strategyID,
			"user_id":     userID,
		})
		return nil, fmt.Errorf("failed to get strategy: %w", err)
	}

	return strategy, nil
}

// GetStrategiesByUser retrieves all strategies for a user
func (r *StrategyRepository) GetStrategiesByUser(userID int, limit, offset int) ([]*data.Strategy, error) {
	query := `
		SELECT id, user_id, name, description, created_at, updated_at
		FROM strategies 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		utils.LogError(err, "Failed to get strategies by user", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get strategies: %w", err)
	}
	defer rows.Close()

	var strategies []*data.Strategy
	for rows.Next() {
		strategy, err := r.scanStrategy(rows)
		if err != nil {
			utils.LogError(err, "Failed to scan strategy", map[string]interface{}{
				"user_id": userID,
			})
			return nil, fmt.Errorf("failed to scan strategy: %w", err)
		}
		strategies = append(strategies, strategy)
	}

	return strategies, nil
}

// DeleteStrategy deletes a strategy
func (r *StrategyRepository) DeleteStrategy(strategyID string, userID int) error {
	query := "DELETE FROM strategies WHERE id = ? AND user_id = ?"

	result, err := r.db.Exec(query, strategyID, userID)
	if err != nil {
		utils.LogError(err, "Failed to delete strategy", map[string]interface{}{
			"strategy_id": strategyID,
			"user_id":     userID,
		})
		return fmt.Errorf("failed to delete strategy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("strategy not found or not owned by user")
	}

	utils.LogInfo("Strategy deleted successfully", map[string]interface{}{
		"strategy_id": strategyID,
		"user_id":     userID,
	})
	return nil
}

// scanStrategy scans a database row into a Strategy struct
func (r *StrategyRepository) scanStrategy(scanner interface {
	Scan(dest ...interface{}) error
}) (*data.Strategy, error) {
	var strategy data.Strategy
	var createdAt, updatedAt time.Time

	err := scanner.Scan(
		&strategy.ID, &strategy.UserID, &strategy.Name, &strategy.Description,
		&createdAt, &updatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Set time fields
	strategy.CreatedAt = createdAt
	strategy.UpdatedAt = updatedAt

	return &strategy, nil
}
