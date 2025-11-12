package repos

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(user *data.User) error {
	// Convert configuredBrokers to JSON
	var configuredBrokersJSON []byte
	var err error
	if user.ConfiguredBrokers == nil {
		configuredBrokersJSON = []byte("{}")
	} else {
		configuredBrokersJSON, err = json.Marshal(user.ConfiguredBrokers)
		if err != nil {
			return fmt.Errorf("failed to marshal configured_brokers: %w", err)
		}
	}

	query := `INSERT INTO users (name, email, phone, configured_brokers, created_at) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, user.Name, user.Email, user.Phone, string(configuredBrokersJSON), user.CreatedAt)
	if err != nil {
		utils.LogError(err, "Failed to create user")
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.LogError(err, "Failed to get last insert ID for user")
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	user.ID = int(id)
	return nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id string) (*data.User, error) {
	query := `SELECT id, name, email, phone, last_signed_in, configured_brokers, created_at FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	user := &data.User{}
	var configuredBrokersJSON string
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.LastSignedIn, &configuredBrokersJSON, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		utils.LogError(err, "Failed to get user by ID", map[string]interface{}{
			"user_id": id,
		})
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse configuredBrokers JSON
	if configuredBrokersJSON != "" && configuredBrokersJSON != "{}" {
		if err := json.Unmarshal([]byte(configuredBrokersJSON), &user.ConfiguredBrokers); err != nil {
			utils.LogError(err, "Failed to unmarshal configured_brokers", map[string]interface{}{
				"user_id": id,
			})
			// Set empty map if unmarshal fails
			user.ConfiguredBrokers = make(map[string]data.BrokerConfig)
		}
	} else {
		user.ConfiguredBrokers = make(map[string]data.BrokerConfig)
	}

	return user, nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(user *data.User) error {
	// Convert configuredBrokers to JSON
	var configuredBrokersJSON []byte
	var err error
	if user.ConfiguredBrokers == nil {
		configuredBrokersJSON = []byte("{}")
	} else {
		configuredBrokersJSON, err = json.Marshal(user.ConfiguredBrokers)
		if err != nil {
			return fmt.Errorf("failed to marshal configured_brokers: %w", err)
		}
	}

	query := `UPDATE users SET name = ?, email = ?, phone = ?, configured_brokers = ? WHERE id = ?`
	_, err = r.db.Exec(query, user.Name, user.Email, user.Phone, string(configuredBrokersJSON), user.ID)
	if err != nil {
		utils.LogError(err, "Failed to update user", map[string]interface{}{
			"user_id": user.ID,
		})
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdateUserBrokerConfig updates the broker configuration for a user
func (r *UserRepository) UpdateUserBrokerConfig(userID int, brokerName string, config data.BrokerConfig) error {
	// Get current user
	user, err := r.GetUserByID(fmt.Sprintf("%d", userID))
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Initialize map if nil
	if user.ConfiguredBrokers == nil {
		user.ConfiguredBrokers = make(map[string]data.BrokerConfig)
	}

	// Update broker config
	user.ConfiguredBrokers[brokerName] = config

	// Save updated user
	return r.UpdateUser(user)
}

// DeleteUser deletes a user by ID
func (r *UserRepository) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		utils.LogError(err, "Failed to delete user", map[string]interface{}{
			"user_id": id,
		})
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogError(err, "Failed to get rows affected for user deletion")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetUsers retrieves users with pagination
func (r *UserRepository) GetUsers(limit, offset int) ([]*data.User, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM users`
	var total int
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		utils.LogError(err, "Failed to get user count")
		return nil, 0, fmt.Errorf("failed to get user count: %w", err)
	}

	// Get users with pagination
	query := `SELECT id, name, email, phone, last_signed_in, configured_brokers, created_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		utils.LogError(err, "Failed to get users")
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*data.User
	for rows.Next() {
		user := &data.User{}
		var configuredBrokersJSON string
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.LastSignedIn, &configuredBrokersJSON, &user.CreatedAt)
		if err != nil {
			utils.LogError(err, "Failed to scan user row")
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}

		// Parse configuredBrokers JSON
		if configuredBrokersJSON != "" && configuredBrokersJSON != "{}" {
			if err := json.Unmarshal([]byte(configuredBrokersJSON), &user.ConfiguredBrokers); err != nil {
				utils.LogError(err, "Failed to unmarshal configured_brokers")
				user.ConfiguredBrokers = make(map[string]data.BrokerConfig)
			}
		} else {
			user.ConfiguredBrokers = make(map[string]data.BrokerConfig)
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		utils.LogError(err, "Error iterating user rows")
		return nil, 0, fmt.Errorf("failed to iterate users: %w", err)
	}

	return users, total, nil
}
