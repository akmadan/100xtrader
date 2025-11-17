// @title 100xtrader API
// @version 1.0
// @description A comprehensive trading journal API for tracking trades, setups, and market analysis
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

package main

import (
	"fmt"
	"os"
	"path/filepath"

	_ "go-core/docs" // Import docs for swagger
	"go-core/internal/api"
	"go-core/internal/data"
	"go-core/internal/utils"
)

func main() {
	// Initialize logger
	utils.InitLogger()

	// Set up database path - use environment variable or default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Get the current working directory
		wd, err := os.Getwd()
		if err != nil {
			utils.LogFatal("Failed to get working directory", map[string]interface{}{
				"error": err.Error(),
			})
		}
		dbPath = filepath.Join(wd, "..", "db.sqlite")
	}

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
		"users", "trades", "strategies", "rules", "mistakes",
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
