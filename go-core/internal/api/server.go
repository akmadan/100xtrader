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

		// Dhan broker routes (separate group to avoid route conflict)
		dhan := v1.Group("/users/:id/dhan")
		{
			dhan.POST("/save-credentials", handlers.SaveDhanCredentials(s.db)) // Save API key & secret
			dhan.POST("/generate-consent", handlers.GenerateDhanConsent(s.db)) // Step 1: Generate consent
			dhan.POST("/consume-consent", handlers.ConsumeDhanConsent(s.db))   // Step 3: Consume consent
			dhan.POST("/renew-token", handlers.RenewDhanToken(s.db))           // Renew token (auto-refresh)
			dhan.GET("/config", handlers.GetDhanBrokerConfig(s.db))
		}

		// Dhan OAuth callback webhook (separate route, accepts tokenId and userId as query params)
		v1.GET("/dhan/consent-callback", handlers.ConsumeDhanConsentCallback(s.db))

		// Trade routes
		trades := v1.Group("/trades")
		{
			trades.POST("", handlers.CreateTrade(s.db))                  // Create trade
			trades.GET("/:id", handlers.GetTrade(s.db))                  // Get trade
			trades.PUT("/:id", handlers.UpdateTrade(s.db))               // Update trade
			trades.DELETE("/:id", handlers.DeleteTrade(s.db))            // Delete trade
			trades.GET("", handlers.ListTrades(s.db))                    // List trades
			trades.GET("/user/:user_id", handlers.GetTradesByUser(s.db)) // Get user's trades
		}

		// User-specific trade routes (use :id to match other user routes)
		userTrades := v1.Group("/users/:id/trades")
		{
			userTrades.POST("/sync-dhan", handlers.SyncDhanTrades(s.db)) // Sync Dhan trades
		}

		// Strategy routes
		strategies := v1.Group("/strategies")
		{
			strategies.POST("", handlers.CreateStrategy(s.db))                   // Create strategy
			strategies.GET("/:id", handlers.GetStrategy(s.db))                   // Get strategy
			strategies.PUT("/:id", handlers.UpdateStrategy(s.db))                // Update strategy
			strategies.DELETE("/:id", handlers.DeleteStrategy(s.db))             // Delete strategy
			strategies.GET("", handlers.ListStrategies(s.db))                    // List strategies
			strategies.GET("/user/:user_id", handlers.GetStrategiesByUser(s.db)) // Get user's strategies
		}

		// Rule routes
		rules := v1.Group("/rules")
		{
			rules.POST("", handlers.CreateRule(s.db))                                         // Create rule
			rules.GET("/:id", handlers.GetRule(s.db))                                         // Get rule
			rules.PUT("/:id", handlers.UpdateRule(s.db))                                      // Update rule
			rules.DELETE("/:id", handlers.DeleteRule(s.db))                                   // Delete rule
			rules.GET("", handlers.ListRules(s.db))                                           // List rules
			rules.GET("/user/:user_id", handlers.GetRulesByUser(s.db))                        // Get user's rules
			rules.GET("/user/:user_id/category/:category", handlers.GetRulesByCategory(s.db)) // Get user's rules by category
		}

		// Algorithm routes
		algorithms := v1.Group("/algorithms")
		{
			algorithms.POST("", handlers.CreateAlgorithm(s.db))
			algorithms.GET("/:id", handlers.GetAlgorithm(s.db))
			algorithms.PUT("/:id", handlers.UpdateAlgorithm(s.db))
			algorithms.DELETE("/:id", handlers.DeleteAlgorithm(s.db))
		}

		// User-specific algorithm routes (use :id to match other user routes)
		userAlgorithms := v1.Group("/users/:id/algorithms")
		{
			userAlgorithms.GET("", handlers.GetAlgorithmsByUser(s.db))
		}

		// Mistake routes
		mistakes := v1.Group("/mistakes")
		{
			mistakes.POST("", handlers.CreateMistake(s.db))                                         // Create mistake
			mistakes.GET("/:id", handlers.GetMistake(s.db))                                         // Get mistake
			mistakes.PUT("/:id", handlers.UpdateMistake(s.db))                                      // Update mistake
			mistakes.DELETE("/:id", handlers.DeleteMistake(s.db))                                   // Delete mistake
			mistakes.GET("", handlers.ListMistakes(s.db))                                           // List mistakes
			mistakes.GET("/user/:user_id", handlers.GetMistakesByUser(s.db))                        // Get user's mistakes
			mistakes.GET("/user/:user_id/category/:category", handlers.GetMistakesByCategory(s.db)) // Get user's mistakes by category
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
