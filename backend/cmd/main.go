// @title 100xtrader API
// @version 1.0
// @description Local Trading Simulation Platform API
// @BasePath /
package main

import (
	"net/http"

	"github.com/akshitmadan/100xtrader/backend/internal/api"
	"github.com/akshitmadan/100xtrader/backend/internal/config"
	"github.com/akshitmadan/100xtrader/backend/internal/data"
	"github.com/akshitmadan/100xtrader/backend/internal/data/repos"
	"github.com/akshitmadan/100xtrader/backend/internal/engine"
	"github.com/akshitmadan/100xtrader/backend/internal/utils"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/akshitmadan/100xtrader/backend/docs"
)

func main() {
	utils.InitLogger()
	config.LoadConfig("../../")

	db := data.InitDB(config.AppConfig.DBSource)
	data.RunMigrations(db, "migrations/001_init.sql")

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Wire up environment repository and handler
	envRepo := repos.NewEnvironmentRepository(db)
	envMgr := engine.GetEnvironmentManager(envRepo)
	envHandler := api.NewEnvironmentHandler(envRepo, envMgr)
	// Register all environment routes with the unified handler
	api.RegisterEnvironmentRoutes(r, envHandler)

	// Wire up order, trade, and position repositories and handler
	orderRepo := repos.NewOrderRepository(db)
	tradeRepo := repos.NewTradeRepository(db)
	posRepo := repos.NewPositionRepository(db)
	orderHandler := api.NewOrderHandler(orderRepo, tradeRepo, posRepo)
	api.RegisterOrderRoutes(r, orderHandler)

	// Wire up portfolio handler
	portfolioHandler := api.NewPortfolioHandler(posRepo, tradeRepo)
	api.RegisterPortfolioRoutes(r, portfolioHandler)

	sessionRepo := repos.NewSessionRepository(db)
	r.POST("/sessions", api.CreateSession(sessionRepo))
	r.POST("/sessions/end", api.EndSession(sessionRepo))
	r.GET("/sessions", api.ListSessions(sessionRepo))

	api.RegisterMarketRoutes(r)

	r.GET("/tickers", api.ListTickers)
	r.POST("/tickers", api.AddTicker)
	r.DELETE("/tickers/:symbol", api.RemoveTicker)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(config.AppConfig.ServerPort) // listen and serve on configured port
}
