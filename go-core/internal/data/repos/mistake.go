package repos

import (
	"database/sql"
	"fmt"
	"time"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// MistakeRepository handles mistake database operations
type MistakeRepository struct {
	db *sql.DB
}

// NewMistakeRepository creates a new mistake repository
func NewMistakeRepository(db *sql.DB) *MistakeRepository {
	return &MistakeRepository{db: db}
}

// CreateMistake creates a new mistake
func (r *MistakeRepository) CreateMistake(mistake *data.Mistake) error {
	query := `
		INSERT INTO mistakes (id, user_id, name, category, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		mistake.ID, mistake.UserID, mistake.Name, mistake.Category,
		mistake.CreatedAt, mistake.UpdatedAt,
	)

	if err != nil {
		utils.LogError(err, "Failed to create mistake", map[string]interface{}{
			"mistake_id": mistake.ID,
			"user_id":    mistake.UserID,
		})
		return fmt.Errorf("failed to create mistake: %w", err)
	}

	utils.LogInfo("Mistake created successfully", map[string]interface{}{
		"mistake_id": mistake.ID,
		"user_id":    mistake.UserID,
	})
	return nil
}

// UpdateMistake updates an existing mistake
func (r *MistakeRepository) UpdateMistake(mistake *data.Mistake) error {
	query := `
		UPDATE mistakes SET 
			name = ?, category = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(query,
		mistake.Name, mistake.Category, mistake.UpdatedAt,
		mistake.ID, mistake.UserID,
	)

	if err != nil {
		utils.LogError(err, "Failed to update mistake", map[string]interface{}{
			"mistake_id": mistake.ID,
			"user_id":    mistake.UserID,
		})
		return fmt.Errorf("failed to update mistake: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("mistake not found or not owned by user")
	}

	utils.LogInfo("Mistake updated successfully", map[string]interface{}{
		"mistake_id": mistake.ID,
		"user_id":    mistake.UserID,
	})
	return nil
}

// GetMistakeByID retrieves a mistake by ID
func (r *MistakeRepository) GetMistakeByID(mistakeID string, userID int) (*data.Mistake, error) {
	query := `
		SELECT id, user_id, name, category, created_at, updated_at
		FROM mistakes 
		WHERE id = ? AND user_id = ?
	`

	row := r.db.QueryRow(query, mistakeID, userID)
	mistake, err := r.scanMistake(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("mistake not found")
		}
		utils.LogError(err, "Failed to get mistake by ID", map[string]interface{}{
			"mistake_id": mistakeID,
			"user_id":    userID,
		})
		return nil, fmt.Errorf("failed to get mistake: %w", err)
	}

	return mistake, nil
}

// GetMistakesByUser retrieves all mistakes for a user
func (r *MistakeRepository) GetMistakesByUser(userID int, limit, offset int) ([]*data.Mistake, error) {
	query := `
		SELECT id, user_id, name, category, created_at, updated_at
		FROM mistakes 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		utils.LogError(err, "Failed to get mistakes by user", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get mistakes: %w", err)
	}
	defer rows.Close()

	var mistakes []*data.Mistake
	for rows.Next() {
		mistake, err := r.scanMistake(rows)
		if err != nil {
			utils.LogError(err, "Failed to scan mistake", map[string]interface{}{
				"user_id": userID,
			})
			return nil, fmt.Errorf("failed to scan mistake: %w", err)
		}
		mistakes = append(mistakes, mistake)
	}

	return mistakes, nil
}

// GetMistakesByCategory retrieves mistakes by category for a user
func (r *MistakeRepository) GetMistakesByCategory(userID int, category data.MistakeCategory, limit, offset int) ([]*data.Mistake, error) {
	query := `
		SELECT id, user_id, name, category, created_at, updated_at
		FROM mistakes 
		WHERE user_id = ? AND category = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, category, limit, offset)
	if err != nil {
		utils.LogError(err, "Failed to get mistakes by category", map[string]interface{}{
			"user_id":  userID,
			"category": category,
		})
		return nil, fmt.Errorf("failed to get mistakes: %w", err)
	}
	defer rows.Close()

	var mistakes []*data.Mistake
	for rows.Next() {
		mistake, err := r.scanMistake(rows)
		if err != nil {
			utils.LogError(err, "Failed to scan mistake", map[string]interface{}{
				"user_id":  userID,
				"category": category,
			})
			return nil, fmt.Errorf("failed to scan mistake: %w", err)
		}
		mistakes = append(mistakes, mistake)
	}

	return mistakes, nil
}

// DeleteMistake deletes a mistake
func (r *MistakeRepository) DeleteMistake(mistakeID string, userID int) error {
	query := "DELETE FROM mistakes WHERE id = ? AND user_id = ?"

	result, err := r.db.Exec(query, mistakeID, userID)
	if err != nil {
		utils.LogError(err, "Failed to delete mistake", map[string]interface{}{
			"mistake_id": mistakeID,
			"user_id":    userID,
		})
		return fmt.Errorf("failed to delete mistake: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("mistake not found or not owned by user")
	}

	utils.LogInfo("Mistake deleted successfully", map[string]interface{}{
		"mistake_id": mistakeID,
		"user_id":    userID,
	})
	return nil
}

// scanMistake scans a database row into a Mistake struct
func (r *MistakeRepository) scanMistake(scanner interface {
	Scan(dest ...interface{}) error
}) (*data.Mistake, error) {
	var mistake data.Mistake
	var createdAt, updatedAt time.Time

	err := scanner.Scan(
		&mistake.ID, &mistake.UserID, &mistake.Name, &mistake.Category,
		&createdAt, &updatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Set time fields
	mistake.CreatedAt = createdAt
	mistake.UpdatedAt = updatedAt

	return &mistake, nil
}
