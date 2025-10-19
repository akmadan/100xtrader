package repos

import (
	"database/sql"
	"fmt"
	"time"

	"100xtrader/go-core/internal/data"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *data.User) error {
	query := `
		INSERT INTO users (name, email, phone, last_signed_in, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	user.CreatedAt = now

	result, err := r.db.Exec(query, user.Name, user.Email, user.Phone, user.LastSignedIn, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	user.ID = int(id)
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*data.User, error) {
	query := `
		SELECT id, name, email, phone, last_signed_in, created_at
		FROM users WHERE id = ?
	`

	user := &data.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone,
		&user.LastSignedIn, &user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*data.User, error) {
	query := `
		SELECT id, name, email, phone, last_signed_in, created_at
		FROM users WHERE email = ?
	`

	user := &data.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone,
		&user.LastSignedIn, &user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *data.User) error {
	query := `
		UPDATE users 
		SET name = ?, email = ?, phone = ?, last_signed_in = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, user.Name, user.Email, user.Phone, user.LastSignedIn, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List retrieves all users
func (r *UserRepository) List() ([]*data.User, error) {
	query := `
		SELECT id, name, email, phone, last_signed_in, created_at
		FROM users ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*data.User
	for rows.Next() {
		user := &data.User{}
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Phone,
			&user.LastSignedIn, &user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
