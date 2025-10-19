package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"100xtrader/go-core/internal/data"
)

func main() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	// Set up database path
	dbPath := filepath.Join(wd, "..", "temp_db.sqlite")

	fmt.Printf("Initializing database at: %s\n", dbPath)

	// Initialize database
	db, err := data.NewDB(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Set global database instance
	data.SetDB(db)

	// Create indexes for better performance
	if err := db.CreateIndexes(); err != nil {
		log.Fatal("Failed to create indexes:", err)
	}

	// Run migrations
	migrationRunner := data.NewMigrationRunner(db.GetConnection())
	migrationsDir := filepath.Join(wd, "migrations")

	if err := migrationRunner.RunMigrations(migrationsDir); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	fmt.Println("✅ Database initialized successfully!")
	fmt.Println("✅ All tables created with proper relationships")
	fmt.Println("✅ Indexes created for optimal performance")
	fmt.Println("✅ Migrations applied successfully")

	// Test database connection
	if err := testDatabase(db); err != nil {
		log.Fatal("Database test failed:", err)
	}

	fmt.Println("✅ Database test passed!")

	// Run example operations
	fmt.Println("\n--- Running Example Operations ---")
	Example()
}

// testDatabase performs basic database connectivity and structure tests
func testDatabase(db *data.DB) error {
	conn := db.GetConnection()

	// Test basic connectivity
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Test table existence
	tables := []string{
		"users", "trades", "trade_actions", "trade_journals",
		"tags", "trade_tags", "trade_screenshots", "notes", "trade_setups",
	}

	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		if err := conn.QueryRow(query).Scan(&count); err != nil {
			return fmt.Errorf("table %s test failed: %w", table, err)
		}
		fmt.Printf("✅ Table '%s' exists and is accessible\n", table)
	}

	return nil
}
