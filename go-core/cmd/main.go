package main

import (
	"fmt"
	"os"
	"path/filepath"

	_ "100xtrader/go-core/docs" // Import generated docs
	"100xtrader/go-core/internal/api"
	"100xtrader/go-core/internal/data"
	"100xtrader/go-core/internal/utils"
)

func main() {
	// Initialize logger
	utils.InitLogger()

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		utils.LogFatal("Failed to get working directory", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Set up database path
	dbPath := filepath.Join(wd, "..", "temp_db.sqlite")

	utils.LogInfo("Initializing database", map[string]interface{}{
		"path": dbPath,
	})

	// Initialize database
	db, err := data.NewDB(dbPath)
	if err != nil {
		utils.LogFatal("Failed to initialize database", map[string]interface{}{
			"error": err.Error(),
			"path":  dbPath,
		})
	}
	defer db.Close()

	// Set global database instance
	data.SetDB(db)

	// Create indexes for better performance
	if err := db.CreateIndexes(); err != nil {
		utils.LogFatal("Failed to create indexes", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Run migrations
	migrationRunner := data.NewMigrationRunner(db.GetConnection())
	migrationsDir := filepath.Join(wd, "migrations")

	if err := migrationRunner.RunMigrations(migrationsDir); err != nil {
		utils.LogFatal("Failed to run migrations", map[string]interface{}{
			"error": err.Error(),
			"dir":   migrationsDir,
		})
	}

	utils.LogInfo("Database initialization completed successfully")

	// Test database connection
	if err := testDatabase(db); err != nil {
		utils.LogFatal("Database test failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	utils.LogInfo("Database test passed")

	// Start API server
	utils.LogInfo("Starting API server")
	server := api.NewServer(db)

	// Start the server
	utils.LogInfo("API server ready to accept requests")
	if err := server.Run(":8080"); err != nil {
		utils.LogFatal("Failed to start API server", map[string]interface{}{
			"error": err.Error(),
		})
	}
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
			utils.LogError(err, "Table test failed", map[string]interface{}{
				"table": table,
			})
			return fmt.Errorf("table %s test failed: %w", table, err)
		}
		utils.LogInfo("Table exists and is accessible", map[string]interface{}{
			"table": table,
			"count": count,
		})
	}

	return nil
}
