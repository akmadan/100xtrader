package api

import (
	"net/http"

	"github.com/akshitmadan/100xtrader/backend/internal/api/dto"
	"github.com/akshitmadan/100xtrader/backend/internal/data"
	"github.com/akshitmadan/100xtrader/backend/internal/engine"
	"github.com/gin-gonic/gin"
)

// ListTickers godoc
// @Summary List all active tickers
// @Description Returns all currently active tickers
// @Tags tickers
// @Produce json
// @Success 200 {object} dto.TickerListResponse
// @Router /tickers [get]
func ListTickers(c *gin.Context) {
	tickers := engine.GetTickerManager(nil).ListTickers() // Use DI in main.go
	var resp dto.TickerListResponse
	for _, t := range tickers {
		resp.Tickers = append(resp.Tickers, t.Symbol)
	}
	c.JSON(http.StatusOK, resp)
}

// AddTicker godoc
// @Summary Add a new ticker
// @Description Adds a new ticker and initializes its market data and order book
// @Tags tickers
// @Accept json
// @Produce json
// @Param request body dto.AddTickerRequest true "Add ticker request"
// @Success 200 {object} dto.AddTickerResponse
// @Failure 400 {object} map[string]string
// @Router /tickers [post]
func AddTicker(c *gin.Context) {
	var req dto.AddTickerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	ticker := &data.Ticker{Symbol: req.Symbol, Name: req.Name}
	if err := engine.GetTickerManager(nil).AddTicker(ticker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add ticker"})
		return
	}
	engine.StartMarketSimulation(req.Symbol)
	c.JSON(http.StatusOK, dto.AddTickerResponse{Status: "added"})
}

// RemoveTicker godoc
// @Summary Remove a ticker
// @Description Removes a ticker and cleans up its resources
// @Tags tickers
// @Param symbol path string true "Ticker symbol"
// @Produce json
// @Success 200 {object} dto.RemoveTickerResponse
// @Router /tickers/{symbol} [delete]
func RemoveTicker(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing symbol parameter"})
		return
	}
	if err := engine.GetTickerManager(nil).RemoveTicker(symbol); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove ticker"})
		return
	}
	c.JSON(http.StatusOK, dto.RemoveTickerResponse{Status: "removed"})
}
