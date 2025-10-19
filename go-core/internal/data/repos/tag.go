package repos

import (
	"database/sql"
	"fmt"
	"time"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// TagRepository handles tag-related database operations
type TagRepository struct {
	db *sql.DB
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Create creates a new tag
func (r *TagRepository) Create(tag *data.Tag) error {
	start := time.Now()
	query := `
		INSERT INTO tags (name, created_at)
		VALUES (?, ?)
	`

	now := time.Now()
	tag.CreatedAt = now

	result, err := r.db.Exec(query, tag.Name, tag.CreatedAt)
	duration := time.Since(start)

	if err != nil {
		utils.LogDatabase("CREATE", "tags", duration, err, map[string]interface{}{
			"name": tag.Name,
		})
		return fmt.Errorf("failed to create tag: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.LogError(err, "Failed to get tag ID", map[string]interface{}{
			"name": tag.Name,
		})
		return fmt.Errorf("failed to get tag ID: %w", err)
	}

	tag.ID = int(id)
	utils.LogDatabase("CREATE", "tags", duration, nil, map[string]interface{}{
		"tag_id": tag.ID,
		"name":   tag.Name,
	})

	return nil
}

// GetByID retrieves a tag by ID
func (r *TagRepository) GetByID(id int) (*data.Tag, error) {
	query := `
		SELECT id, name, created_at
		FROM tags WHERE id = ?
	`

	tag := &data.Tag{}
	err := r.db.QueryRow(query, id).Scan(
		&tag.ID, &tag.Name, &tag.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found")
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return tag, nil
}

// GetByName retrieves a tag by name
func (r *TagRepository) GetByName(name string) (*data.Tag, error) {
	query := `
		SELECT id, name, created_at
		FROM tags WHERE name = ?
	`

	tag := &data.Tag{}
	err := r.db.QueryRow(query, name).Scan(
		&tag.ID, &tag.Name, &tag.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Tag not found, not an error
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return tag, nil
}

// GetOrCreate gets an existing tag or creates a new one
func (r *TagRepository) GetOrCreate(name string) (*data.Tag, error) {
	// Try to get existing tag
	tag, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}

	// If tag exists, return it
	if tag != nil {
		return tag, nil
	}

	// Create new tag
	tag = &data.Tag{Name: name}
	if err := r.Create(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// GetByTradeID retrieves all tags for a trade
func (r *TagRepository) GetByTradeID(tradeID string) ([]*data.Tag, error) {
	query := `
		SELECT t.id, t.name, t.created_at
		FROM tags t
		JOIN trade_tags tt ON t.id = tt.tag_id
		WHERE tt.trade_id = ?
		ORDER BY t.name
	`

	rows, err := r.db.Query(query, tradeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags for trade: %w", err)
	}
	defer rows.Close()

	var tags []*data.Tag
	for rows.Next() {
		tag := &data.Tag{}
		err := rows.Scan(
			&tag.ID, &tag.Name, &tag.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// AddToTrade adds a tag to a trade
func (r *TagRepository) AddToTrade(tradeID string, tagID int) error {
	query := `
		INSERT INTO trade_tags (trade_id, tag_id, created_at)
		VALUES (?, ?, ?)
	`

	_, err := r.db.Exec(query, tradeID, tagID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add tag to trade: %w", err)
	}

	return nil
}

// RemoveFromTrade removes a tag from a trade
func (r *TagRepository) RemoveFromTrade(tradeID string, tagID int) error {
	query := `DELETE FROM trade_tags WHERE trade_id = ? AND tag_id = ?`

	_, err := r.db.Exec(query, tradeID, tagID)
	if err != nil {
		return fmt.Errorf("failed to remove tag from trade: %w", err)
	}

	return nil
}

// List retrieves all tags
func (r *TagRepository) List() ([]*data.Tag, error) {
	query := `
		SELECT id, name, created_at
		FROM tags ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	defer rows.Close()

	var tags []*data.Tag
	for rows.Next() {
		tag := &data.Tag{}
		err := rows.Scan(
			&tag.ID, &tag.Name, &tag.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
