package repos

import (
	"database/sql"
	"fmt"
	"time"

	"go-core/internal/data"
)

// NotesRepository handles notes-related database operations
type NotesRepository struct {
	db *sql.DB
}

// NewNotesRepository creates a new notes repository
func NewNotesRepository(db *sql.DB) *NotesRepository {
	return &NotesRepository{db: db}
}

// Create creates a new note
func (r *NotesRepository) Create(note *data.Note) error {
	query := `
		INSERT INTO notes (user_id, mood, market_condition, market_volatility, summary, day, notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	note.CreatedAt = now

	result, err := r.db.Exec(query, note.UserID, note.Mood, note.MarketCondition, note.MarketVolatility, note.Summary, note.Day, note.Notes, note.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get note ID: %w", err)
	}

	note.ID = int(id)
	return nil
}

// GetByID retrieves a note by ID
func (r *NotesRepository) GetByID(id int) (*data.Note, error) {
	query := `
		SELECT id, user_id, mood, market_condition, market_volatility, summary, day, notes, created_at
		FROM notes WHERE id = ?
	`

	note := &data.Note{}
	err := r.db.QueryRow(query, id).Scan(
		&note.ID, &note.UserID, &note.Mood, &note.MarketCondition, &note.MarketVolatility,
		&note.Summary, &note.Day, &note.Notes, &note.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("note not found")
		}
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	return note, nil
}

// GetByUserID retrieves all notes for a user
func (r *NotesRepository) GetByUserID(userID int) ([]*data.Note, error) {
	query := `
		SELECT id, user_id, mood, market_condition, market_volatility, summary, day, notes, created_at
		FROM notes WHERE user_id = ? ORDER BY day DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}
	defer rows.Close()

	var notes []*data.Note
	for rows.Next() {
		note := &data.Note{}
		err := rows.Scan(
			&note.ID, &note.UserID, &note.Mood, &note.MarketCondition, &note.MarketVolatility,
			&note.Summary, &note.Day, &note.Notes, &note.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// GetByDate retrieves notes by date
func (r *NotesRepository) GetByDate(userID int, date time.Time) ([]*data.Note, error) {
	query := `
		SELECT id, user_id, mood, market_condition, market_volatility, summary, day, notes, created_at
		FROM notes WHERE user_id = ? AND DATE(day) = DATE(?)
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes by date: %w", err)
	}
	defer rows.Close()

	var notes []*data.Note
	for rows.Next() {
		note := &data.Note{}
		err := rows.Scan(
			&note.ID, &note.UserID, &note.Mood, &note.MarketCondition, &note.MarketVolatility,
			&note.Summary, &note.Day, &note.Notes, &note.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// Update updates a note
func (r *NotesRepository) Update(note *data.Note) error {
	query := `
		UPDATE notes 
		SET mood = ?, market_condition = ?, market_volatility = ?, summary = ?, day = ?, notes = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, note.Mood, note.MarketCondition, note.MarketVolatility, note.Summary, note.Day, note.Notes, note.ID)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	return nil
}

// Delete deletes a note
func (r *NotesRepository) Delete(id int) error {
	query := `DELETE FROM notes WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	return nil
}
