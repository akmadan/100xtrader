package data

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"100xtrader/go-core/internal/utils"

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

// InitTables creates all the necessary tables
func (db *DB) InitTables() error {
	queries := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			phone VARCHAR(20),
			last_signed_in TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Trades table
		`CREATE TABLE IF NOT EXISTS trades (
			id VARCHAR(255) PRIMARY KEY,
			user_id INTEGER NOT NULL,
			market VARCHAR(50) NOT NULL CHECK (market IN ('stock', 'option', 'crypto', 'futures', 'forex', 'index')),
			symbol VARCHAR(50) NOT NULL,
			target REAL NOT NULL,
			stoploss REAL NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,

		// Trade actions table
		`CREATE TABLE IF NOT EXISTS trade_actions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			trade_id VARCHAR(255) NOT NULL,
			action VARCHAR(10) NOT NULL CHECK (action IN ('buy', 'sell')),
			trade_time TIMESTAMP NOT NULL,
			quantity INTEGER NOT NULL,
			price REAL NOT NULL,
			fee REAL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (trade_id) REFERENCES trades(id) ON DELETE CASCADE
		)`,

		// Trade journals table
		`CREATE TABLE IF NOT EXISTS trade_journals (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			trade_id VARCHAR(255) NOT NULL,
			notes TEXT,
			confidence INTEGER CHECK (confidence >= 0 AND confidence <= 10),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (trade_id) REFERENCES trades(id) ON DELETE CASCADE
		)`,

		// Tags table
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Trade tags table (many-to-many)
		`CREATE TABLE IF NOT EXISTS trade_tags (
			trade_id VARCHAR(255) NOT NULL,
			tag_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (trade_id, tag_id),
			FOREIGN KEY (trade_id) REFERENCES trades(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)`,

		// Trade screenshots table
		`CREATE TABLE IF NOT EXISTS trade_screenshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			trade_journal_id INTEGER NOT NULL,
			url VARCHAR(500) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (trade_journal_id) REFERENCES trade_journals(id) ON DELETE CASCADE
		)`,

		// Notes table
		`CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			mood VARCHAR(20) NOT NULL CHECK (mood IN ('excited', 'neutral', 'low')),
			market_condition VARCHAR(20) NOT NULL CHECK (market_condition IN ('up', 'down', 'sideways')),
			market_volatility VARCHAR(20) NOT NULL CHECK (market_volatility IN ('high', 'medium', 'low')),
			summary TEXT,
			day TIMESTAMP NOT NULL,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,

		// Trade setups table
		`CREATE TABLE IF NOT EXISTS trade_setups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			market VARCHAR(50) NOT NULL CHECK (market IN ('stock', 'option', 'crypto', 'futures', 'forex', 'index')),
			side VARCHAR(10) NOT NULL CHECK (side IN ('long', 'short')),
			symbol VARCHAR(50) NOT NULL,
			entry REAL NOT NULL,
			target REAL NOT NULL,
			stoploss REAL NOT NULL,
			note TEXT,
			risk_reward_ratio REAL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}

	for _, query := range queries {
		if _, err := db.conn.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	utils.LogInfo("Database tables initialized successfully")
	return nil
}

// CreateIndexes creates useful indexes for better performance
func (db *DB) CreateIndexes() error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_trades_user_id ON trades(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_trades_market ON trades(market)",
		"CREATE INDEX IF NOT EXISTS idx_trades_symbol ON trades(symbol)",
		"CREATE INDEX IF NOT EXISTS idx_trade_actions_trade_id ON trade_actions(trade_id)",
		"CREATE INDEX IF NOT EXISTS idx_trade_journals_trade_id ON trade_journals(trade_id)",
		"CREATE INDEX IF NOT EXISTS idx_trade_tags_trade_id ON trade_tags(trade_id)",
		"CREATE INDEX IF NOT EXISTS idx_trade_tags_tag_id ON trade_tags(tag_id)",
		"CREATE INDEX IF NOT EXISTS idx_trade_screenshots_journal_id ON trade_screenshots(trade_journal_id)",
		"CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_notes_day ON notes(day)",
		"CREATE INDEX IF NOT EXISTS idx_trade_setups_user_id ON trade_setups(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_trade_setups_market ON trade_setups(market)",
		"CREATE INDEX IF NOT EXISTS idx_trade_setups_symbol ON trade_setups(symbol)",
	}

	for _, index := range indexes {
		if _, err := db.conn.Exec(index); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	utils.LogInfo("Database indexes created successfully")
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
