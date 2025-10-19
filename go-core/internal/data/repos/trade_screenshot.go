package repos

import (
	"database/sql"
	"fmt"
	"time"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// TradeScreenshotRepository handles trade screenshot-related database operations
type TradeScreenshotRepository struct {
	db *sql.DB
}

// NewTradeScreenshotRepository creates a new trade screenshot repository
func NewTradeScreenshotRepository(db *sql.DB) *TradeScreenshotRepository {
	return &TradeScreenshotRepository{db: db}
}

// Create creates a new trade screenshot
func (r *TradeScreenshotRepository) Create(screenshot *data.TradeScreenshot) error {
	start := time.Now()
	query := `
		INSERT INTO trade_screenshots (trade_journal_id, url, created_at)
		VALUES (?, ?, ?)
	`

	now := time.Now()
	screenshot.CreatedAt = now

	result, err := r.db.Exec(query, screenshot.TradeJournalID, screenshot.URL, screenshot.CreatedAt)
	duration := time.Since(start)

	if err != nil {
		utils.LogDatabase("CREATE", "trade_screenshots", duration, err, map[string]interface{}{
			"trade_journal_id": screenshot.TradeJournalID,
		})
		return fmt.Errorf("failed to create trade screenshot: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.LogError(err, "Failed to get trade screenshot ID", map[string]interface{}{
			"trade_journal_id": screenshot.TradeJournalID,
		})
		return fmt.Errorf("failed to get trade screenshot ID: %w", err)
	}

	screenshot.ID = int(id)
	utils.LogDatabase("CREATE", "trade_screenshots", duration, nil, map[string]interface{}{
		"screenshot_id": screenshot.ID,
		"journal_id":    screenshot.TradeJournalID,
	})

	return nil
}

// GetByID retrieves a trade screenshot by ID
func (r *TradeScreenshotRepository) GetByID(id int) (*data.TradeScreenshot, error) {
	query := `
		SELECT id, trade_journal_id, url, created_at
		FROM trade_screenshots WHERE id = ?
	`

	screenshot := &data.TradeScreenshot{}
	err := r.db.QueryRow(query, id).Scan(
		&screenshot.ID, &screenshot.TradeJournalID, &screenshot.URL, &screenshot.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trade screenshot not found")
		}
		return nil, fmt.Errorf("failed to get trade screenshot: %w", err)
	}

	return screenshot, nil
}

// GetByJournalID retrieves all screenshots for a trade journal
func (r *TradeScreenshotRepository) GetByJournalID(journalID int) ([]*data.TradeScreenshot, error) {
	query := `
		SELECT id, trade_journal_id, url, created_at
		FROM trade_screenshots WHERE trade_journal_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, journalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade screenshots: %w", err)
	}
	defer rows.Close()

	var screenshots []*data.TradeScreenshot
	for rows.Next() {
		screenshot := &data.TradeScreenshot{}
		err := rows.Scan(
			&screenshot.ID, &screenshot.TradeJournalID, &screenshot.URL, &screenshot.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade screenshot: %w", err)
		}
		screenshots = append(screenshots, screenshot)
	}

	return screenshots, nil
}

// Delete deletes a trade screenshot
func (r *TradeScreenshotRepository) Delete(id int) error {
	query := `DELETE FROM trade_screenshots WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trade screenshot: %w", err)
	}

	return nil
}
