package data

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"go-core/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

// DB represents the database connection
type DB struct {
	conn *sql.DB
}

// NewDB creates a new database connection
func NewDB(dbPath string) (*DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}

	// Initialize the database with tables
	if err := db.InitTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// GetConnection returns the underlying database connection
func (db *DB) GetConnection() *sql.DB {
	return db.conn
}

// InitTables runs database migrations to set up the schema
func (db *DB) InitTables() error {
	// Get the migrations directory path
	migrationsDir := "migrations"

	// Create migration runner
	migrationRunner := NewMigrationRunner(db.conn)

	// Run all pending migrations
	if err := migrationRunner.RunMigrations(migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	utils.LogInfo("Database migrations completed successfully")
	return nil
}

// CreateIndexes creates useful indexes for better performance
// Note: Indexes are now created as part of the migration process
func (db *DB) CreateIndexes() error {
	// Indexes are now created in the migration files
	// This method is kept for backward compatibility but does nothing
	utils.LogInfo("Database indexes are managed by migrations")
	return nil
}

// GetDBInstance returns a singleton database instance
var dbInstance *DB

// GetDB returns the database instance
func GetDB() *DB {
	return dbInstance
}

// SetDB sets the database instance
func SetDB(db *DB) {
	dbInstance = db
}
