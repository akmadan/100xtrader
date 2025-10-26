package repos

import (
	"database/sql"
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
	query := `INSERT INTO users (name, email, phone, created_at) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, user.Name, user.Email, user.Phone, user.CreatedAt)
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
	query := `SELECT id, name, email, phone, last_signed_in, created_at FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	user := &data.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.LastSignedIn, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		utils.LogError(err, "Failed to get user by ID", map[string]interface{}{
			"user_id": id,
		})
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(user *data.User) error {
	query := `UPDATE users SET name = ?, email = ?, phone = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Phone, user.ID)
	if err != nil {
		utils.LogError(err, "Failed to update user", map[string]interface{}{
			"user_id": user.ID,
		})
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
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
	query := `SELECT id, name, email, phone, last_signed_in, created_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		utils.LogError(err, "Failed to get users")
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*data.User
	for rows.Next() {
		user := &data.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.LastSignedIn, &user.CreatedAt)
		if err != nil {
			utils.LogError(err, "Failed to scan user row")
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		utils.LogError(err, "Error iterating user rows")
		return nil, 0, fmt.Errorf("failed to iterate users: %w", err)
	}

	return users, total, nil
}
