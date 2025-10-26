package repos

import (
	"database/sql"
	"fmt"
	"time"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// RuleRepository handles rule database operations
type RuleRepository struct {
	db *sql.DB
}

// NewRuleRepository creates a new rule repository
func NewRuleRepository(db *sql.DB) *RuleRepository {
	return &RuleRepository{db: db}
}

// CreateRule creates a new rule
func (r *RuleRepository) CreateRule(rule *data.Rule) error {
	query := `
		INSERT INTO rules (id, user_id, name, description, category, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		rule.ID, rule.UserID, rule.Name, rule.Description, rule.Category,
		rule.CreatedAt, rule.UpdatedAt,
	)

	if err != nil {
		utils.LogError(err, "Failed to create rule", map[string]interface{}{
			"rule_id": rule.ID,
			"user_id": rule.UserID,
		})
		return fmt.Errorf("failed to create rule: %w", err)
	}

	utils.LogInfo("Rule created successfully", map[string]interface{}{
		"rule_id": rule.ID,
		"user_id": rule.UserID,
	})
	return nil
}

// UpdateRule updates an existing rule
func (r *RuleRepository) UpdateRule(rule *data.Rule) error {
	query := `
		UPDATE rules SET 
			name = ?, description = ?, category = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(query,
		rule.Name, rule.Description, rule.Category, rule.UpdatedAt,
		rule.ID, rule.UserID,
	)

	if err != nil {
		utils.LogError(err, "Failed to update rule", map[string]interface{}{
			"rule_id": rule.ID,
			"user_id": rule.UserID,
		})
		return fmt.Errorf("failed to update rule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rule not found or not owned by user")
	}

	utils.LogInfo("Rule updated successfully", map[string]interface{}{
		"rule_id": rule.ID,
		"user_id": rule.UserID,
	})
	return nil
}

// GetRuleByID retrieves a rule by ID
func (r *RuleRepository) GetRuleByID(ruleID string, userID int) (*data.Rule, error) {
	query := `
		SELECT id, user_id, name, description, category, created_at, updated_at
		FROM rules 
		WHERE id = ? AND user_id = ?
	`

	row := r.db.QueryRow(query, ruleID, userID)
	rule, err := r.scanRule(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("rule not found")
		}
		utils.LogError(err, "Failed to get rule by ID", map[string]interface{}{
			"rule_id": ruleID,
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get rule: %w", err)
	}

	return rule, nil
}

// GetRulesByUser retrieves all rules for a user
func (r *RuleRepository) GetRulesByUser(userID int, limit, offset int) ([]*data.Rule, error) {
	query := `
		SELECT id, user_id, name, description, category, created_at, updated_at
		FROM rules 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		utils.LogError(err, "Failed to get rules by user", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get rules: %w", err)
	}
	defer rows.Close()

	var rules []*data.Rule
	for rows.Next() {
		rule, err := r.scanRule(rows)
		if err != nil {
			utils.LogError(err, "Failed to scan rule", map[string]interface{}{
				"user_id": userID,
			})
			return nil, fmt.Errorf("failed to scan rule: %w", err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// GetRulesByCategory retrieves rules by category for a user
func (r *RuleRepository) GetRulesByCategory(userID int, category data.RuleCategory, limit, offset int) ([]*data.Rule, error) {
	query := `
		SELECT id, user_id, name, description, category, created_at, updated_at
		FROM rules 
		WHERE user_id = ? AND category = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, category, limit, offset)
	if err != nil {
		utils.LogError(err, "Failed to get rules by category", map[string]interface{}{
			"user_id":  userID,
			"category": category,
		})
		return nil, fmt.Errorf("failed to get rules: %w", err)
	}
	defer rows.Close()

	var rules []*data.Rule
	for rows.Next() {
		rule, err := r.scanRule(rows)
		if err != nil {
			utils.LogError(err, "Failed to scan rule", map[string]interface{}{
				"user_id":  userID,
				"category": category,
			})
			return nil, fmt.Errorf("failed to scan rule: %w", err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// DeleteRule deletes a rule
func (r *RuleRepository) DeleteRule(ruleID string, userID int) error {
	query := "DELETE FROM rules WHERE id = ? AND user_id = ?"

	result, err := r.db.Exec(query, ruleID, userID)
	if err != nil {
		utils.LogError(err, "Failed to delete rule", map[string]interface{}{
			"rule_id": ruleID,
			"user_id": userID,
		})
		return fmt.Errorf("failed to delete rule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rule not found or not owned by user")
	}

	utils.LogInfo("Rule deleted successfully", map[string]interface{}{
		"rule_id": ruleID,
		"user_id": userID,
	})
	return nil
}

// scanRule scans a database row into a Rule struct
func (r *RuleRepository) scanRule(scanner interface {
	Scan(dest ...interface{}) error
}) (*data.Rule, error) {
	var rule data.Rule
	var createdAt, updatedAt time.Time

	err := scanner.Scan(
		&rule.ID, &rule.UserID, &rule.Name, &rule.Description, &rule.Category,
		&createdAt, &updatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Set time fields
	rule.CreatedAt = createdAt
	rule.UpdatedAt = updatedAt

	return &rule, nil
}
