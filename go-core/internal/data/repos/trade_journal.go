package repos

import (
	"database/sql"
	"fmt"
	"time"

	"100xtrader/go-core/internal/data"
	"100xtrader/go-core/internal/utils"
)

// TradeJournalRepository handles trade journal-related database operations
type TradeJournalRepository struct {
	db *sql.DB
}

// NewTradeJournalRepository creates a new trade journal repository
func NewTradeJournalRepository(db *sql.DB) *TradeJournalRepository {
	return &TradeJournalRepository{db: db}
}

// Create creates a new trade journal
func (r *TradeJournalRepository) Create(journal *data.TradeJournal) error {
	start := time.Now()
	query := `
		INSERT INTO trade_journals (trade_id, notes, confidence, created_at)
		VALUES (?, ?, ?, ?)
	`

	now := time.Now()
	journal.CreatedAt = now

	result, err := r.db.Exec(query, journal.TradeID, journal.Notes, journal.Confidence, journal.CreatedAt)
	duration := time.Since(start)

	if err != nil {
		utils.LogDatabase("CREATE", "trade_journals", duration, err, map[string]interface{}{
			"trade_id": journal.TradeID,
		})
		return fmt.Errorf("failed to create trade journal: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.LogError(err, "Failed to get trade journal ID", map[string]interface{}{
			"trade_id": journal.TradeID,
		})
		return fmt.Errorf("failed to get trade journal ID: %w", err)
	}

	journal.ID = int(id)
	utils.LogDatabase("CREATE", "trade_journals", duration, nil, map[string]interface{}{
		"journal_id": journal.ID,
		"trade_id":   journal.TradeID,
	})

	return nil
}

// GetByID retrieves a trade journal by ID
func (r *TradeJournalRepository) GetByID(id int) (*data.TradeJournal, error) {
	query := `
		SELECT id, trade_id, notes, confidence, created_at
		FROM trade_journals WHERE id = ?
	`

	journal := &data.TradeJournal{}
	err := r.db.QueryRow(query, id).Scan(
		&journal.ID, &journal.TradeID, &journal.Notes, &journal.Confidence, &journal.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trade journal not found")
		}
		return nil, fmt.Errorf("failed to get trade journal: %w", err)
	}

	return journal, nil
}

// GetByTradeID retrieves journal for a specific trade
func (r *TradeJournalRepository) GetByTradeID(tradeID string) (*data.TradeJournal, error) {
	query := `
		SELECT id, trade_id, notes, confidence, created_at
		FROM trade_journals WHERE trade_id = ?
	`

	journal := &data.TradeJournal{}
	err := r.db.QueryRow(query, tradeID).Scan(
		&journal.ID, &journal.TradeID, &journal.Notes, &journal.Confidence, &journal.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No journal found, not an error
		}
		return nil, fmt.Errorf("failed to get trade journal: %w", err)
	}

	return journal, nil
}

// Update updates a trade journal
func (r *TradeJournalRepository) Update(journal *data.TradeJournal) error {
	query := `
		UPDATE trade_journals 
		SET notes = ?, confidence = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, journal.Notes, journal.Confidence, journal.ID)
	if err != nil {
		return fmt.Errorf("failed to update trade journal: %w", err)
	}

	return nil
}

// Delete deletes a trade journal
func (r *TradeJournalRepository) Delete(id int) error {
	query := `DELETE FROM trade_journals WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trade journal: %w", err)
	}

	return nil
}
