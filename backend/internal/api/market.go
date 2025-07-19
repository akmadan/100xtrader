package api

import "github.com/gin-gonic/gin"

// GetMarketData handles GET /market-data/:symbol
func GetMarketData(c *gin.Context) {
	// TODO: Implement logic to get market data for a symbol
	c.JSON(200, gin.H{"market_data": nil})
}

// GetOrderBook handles GET /orderbook/:symbol
func GetOrderBook(c *gin.Context) {
	// TODO: Implement logic to get order book for a symbol
	c.JSON(200, gin.H{"orderbook": nil})
}

// RegisterMarketRoutes registers market data and orderbook routes
func RegisterMarketRoutes(r *gin.Engine) {
	r.GET("/market-data/:symbol", GetMarketData)
	r.GET("/orderbook/:symbol", GetOrderBook)
}
