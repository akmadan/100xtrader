package api

import (
	"net/http"
	"time"

	"github.com/akshitmadan/100xtrader/backend/internal/api/dto"
	"github.com/akshitmadan/100xtrader/backend/internal/engine"
	"github.com/gin-gonic/gin"
)

// GetMarketData godoc
// @Summary Get market data for a symbol
// @Description Returns the latest price, OHLC, volume, and recent trades for a symbol
// @Tags market
// @Param symbol path string true "Symbol"
// @Produce json
// @Success 200 {object} dto.MarketDataResponse
// @Failure 404 {object} map[string]string
// @Router /market-data/{symbol} [get]
func GetMarketData(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing symbol parameter"})
		return
	}
	lastPrice, lastTradeTime, ohlc, trades := engine.GetMarketDataSnapshot(symbol)
	if lastPrice == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no market data for symbol"})
		return
	}
	var tradeDTOs []dto.TradeTickResponse
	for _, t := range trades {
		tradeDTOs = append(tradeDTOs, dto.TradeTickResponse{
			Price:     t.Price,
			Quantity:  t.Quantity,
			Side:      t.Side,
			Timestamp: t.Timestamp.Format(time.RFC3339),
		})
	}
	resp := dto.MarketDataResponse{
		Symbol:        symbol,
		LastPrice:     lastPrice,
		LastTradeTime: lastTradeTime.Format(time.RFC3339),
		OHLC: dto.OHLCResponse{
			Open:   ohlc.Open,
			High:   ohlc.High,
			Low:    ohlc.Low,
			Close:  ohlc.Close,
			Volume: ohlc.Volume,
			Start:  ohlc.Start.Format(time.RFC3339),
		},
		RecentTrades: tradeDTOs,
	}
	c.JSON(http.StatusOK, resp)
}

// GetOrderBook godoc
// @Summary Get order book for a symbol
// @Description Returns the current order book (bids/asks) for a symbol
// @Tags market
// @Param symbol path string true "Symbol"
// @Produce json
// @Success 200 {object} dto.OrderBookResponse
// @Failure 404 {object} map[string]string
// @Router /orderbook/{symbol} [get]
func GetOrderBook(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing symbol parameter"})
		return
	}
	ob := engine.GetOrderBookManager().GetOrCreateOrderBook(symbol)
	// For now, return top 10 levels
	depth := ob.GetMarketDepth(10)
	var bids, asks []dto.OrderBookLevel
	for _, b := range depth.Bids {
		bids = append(bids, dto.OrderBookLevel{Price: b.Price, Quantity: b.Quantity})
	}
	for _, a := range depth.Asks {
		asks = append(asks, dto.OrderBookLevel{Price: a.Price, Quantity: a.Quantity})
	}
	resp := dto.OrderBookResponse{
		Symbol: symbol,
		Bids:   bids,
		Asks:   asks,
		Time:   time.Now().UTC().Format(time.RFC3339),
	}
	c.JSON(http.StatusOK, resp)
}

// RegisterMarketRoutes registers market data and orderbook routes
func RegisterMarketRoutes(r *gin.Engine) {
	r.GET("/market-data/:symbol", GetMarketData)
	r.GET("/orderbook/:symbol", GetOrderBook)
}
