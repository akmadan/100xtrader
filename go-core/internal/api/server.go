// @title 100xTrader API
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

package api

import (
	"net/http"
	"time"

	"go-core/internal/api/dto"
	"go-core/internal/api/handlers"
	"go-core/internal/api/middleware"
	"go-core/internal/data"
	"go-core/internal/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	db     *data.DB
}

// NewServer creates a new API server
func NewServer(db *data.DB) *Server {
	// Set Gin mode based on environment
	if utils.GetLogger().Level.String() == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	server := &Server{
		router: router,
		db:     db,
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	// Global middleware
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recovery())
	s.router.Use(middleware.CORS())
	s.router.Use(middleware.RequestID())
}

// setupRoutes configures API routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

	// Swagger documentation
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("", handlers.CreateUser(s.db))
			users.GET("/:id", handlers.GetUser(s.db))
			users.PUT("/:id", handlers.UpdateUser(s.db))
			users.DELETE("/:id", handlers.DeleteUser(s.db))
			users.GET("", handlers.ListUsers(s.db))
			users.POST("/signin", handlers.SignInUser(s.db))
		}

		// Trade routes (includes optional journal and actions)
		trades := v1.Group("/trades")
		{
			trades.POST("", handlers.CreateTrade(s.db))                  // Create trade with optional journal and actions
			trades.GET("/:id", handlers.GetTrade(s.db))                  // Get trade with journal and actions
			trades.PUT("/:id", handlers.UpdateTrade(s.db))               // Update trade with journal and actions
			trades.DELETE("/:id", handlers.DeleteTrade(s.db))            // Delete trade
			trades.GET("", handlers.ListTrades(s.db))                    // List trades with optional journal data
			trades.GET("/user/:user_id", handlers.GetTradesByUser(s.db)) // Get user's trades

			// Trade-specific sub-resources
			trades.POST("/:id/actions", handlers.AddTradeAction(s.db))                 // Add action to existing trade
			trades.DELETE("/:id/actions/:action_id", handlers.RemoveTradeAction(s.db)) // Remove action from trade
			trades.POST("/:id/journal", handlers.UpdateTradeJournal(s.db))             // Update journal for trade
			trades.POST("/:id/screenshots", handlers.AddScreenshot(s.db))              // Add screenshot to trade journal
		}

		// Trade Setup routes
		setups := v1.Group("/trade-setups")
		{
			setups.POST("", handlers.CreateTradeSetup(s.db))
			setups.GET("/:id", handlers.GetTradeSetup(s.db))
			setups.PUT("/:id", handlers.UpdateTradeSetup(s.db))
			setups.DELETE("/:id", handlers.DeleteTradeSetup(s.db))
			setups.GET("", handlers.ListTradeSetups(s.db))
			setups.GET("/user/:user_id", handlers.GetSetupsByUser(s.db))
		}

		// Notes routes
		notes := v1.Group("/notes")
		{
			notes.POST("", handlers.CreateNote(s.db))
			notes.GET("/:id", handlers.GetNote(s.db))
			notes.PUT("/:id", handlers.UpdateNote(s.db))
			notes.DELETE("/:id", handlers.DeleteNote(s.db))
			notes.GET("", handlers.ListNotes(s.db))
			notes.GET("/user/:user_id", handlers.GetNotesByUser(s.db))
			notes.GET("/user/:user_id/daily", handlers.GetDailyNotes(s.db))
		}

		// Tags routes
		tags := v1.Group("/tags")
		{
			tags.POST("", handlers.CreateTag(s.db))
			tags.GET("/:id", handlers.GetTag(s.db))
			tags.GET("", handlers.ListTags(s.db))
			tags.POST("/trade", handlers.AddTagToTrade(s.db))
			tags.DELETE("/trade/:trade_id/:tag_id", handlers.RemoveTagFromTrade(s.db))
		}
	}
}

// healthCheck handles health check requests
// @Summary Health check
// @Description Check if the API service is healthy and database is accessible
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} dto.SuccessResponse "Service is healthy"
// @Failure 503 {object} dto.ErrorResponse "Service unavailable"
// @Router /health [get]
func (s *Server) healthCheck(c *gin.Context) {
	// Test database connection
	if err := s.db.GetConnection().Ping(); err != nil {
		utils.LogError(err, "Health check failed - database connection error")
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Error:   "Service Unavailable",
			Message: "Database connection failed",
			Code:    http.StatusServiceUnavailable,
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Service is healthy",
		Data: map[string]interface{}{
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
			"status":    "ok",
		},
	})
}

// Run starts the server
func (s *Server) Run(addr string) error {
	utils.LogInfo("Starting API server", map[string]interface{}{
		"address": addr,
		"mode":    gin.Mode(),
	})

	return s.router.Run(addr)
}

// GetRouter returns the router for testing
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
