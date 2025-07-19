package api

import (
	"context"
	"net/http"
	"time"

	"github.com/akshitmadan/100xtrader/backend/internal/api/dto"
	"github.com/akshitmadan/100xtrader/backend/internal/data/repos"
	"github.com/akshitmadan/100xtrader/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

// PortfolioHandler holds dependencies for portfolio endpoints

type PortfolioHandler struct {
	PosRepo   repos.PositionRepository
	TradeRepo repos.TradeRepository
}

func NewPortfolioHandler(posRepo repos.PositionRepository, tradeRepo repos.TradeRepository) *PortfolioHandler {
	return &PortfolioHandler{
		PosRepo:   posRepo,
		TradeRepo: tradeRepo,
	}
}

// RegisterPortfolioRoutes registers portfolio, positions, and trades routes with DI
func RegisterPortfolioRoutes(r *gin.Engine, handler *PortfolioHandler) {
	r.GET("/portfolio", handler.GetPortfolio)
	r.GET("/positions", handler.GetPositions)
	r.GET("/trades", handler.GetTrades)
}

// GetPortfolio godoc
// @Summary Get user portfolio
// @Description Returns portfolio summary for a user
// @Tags portfolio
// @Param user query string true "User ID"
// @Produce json
// @Success 200 {object} dto.PortfolioResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /portfolio [get]
func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	user := c.Query("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user parameter"})
		return
	}
	positions, err := h.PosRepo.GetPositionsByUser(context.Background(), user)
	if err != nil {
		utils.Logger.WithError(err).Error("failed to fetch positions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch positions"})
		return
	}
	// For demo, mock cash balance and use last trade price for market value
	cash := 100000.0
	totalValue := cash
	var posDTOs []dto.PositionResponse
	for _, p := range positions {
		// For now, use average price as market price (mock)
		marketValue := p.Quantity * p.AveragePrice
		pnl := (p.AveragePrice - p.AveragePrice) * p.Quantity // PnL mock
		totalValue += marketValue
		posDTOs = append(posDTOs, dto.PositionResponse{
			Symbol:       p.Symbol,
			Quantity:     p.Quantity,
			AveragePrice: p.AveragePrice,
			MarketValue:  marketValue,
			PnL:          pnl,
		})
	}
	resp := dto.PortfolioResponse{
		User:        user,
		CashBalance: cash,
		TotalValue:  totalValue,
		Positions:   posDTOs,
		PnL:         0, // mock
	}
	c.JSON(http.StatusOK, resp)
}

// GetPositions godoc
// @Summary Get user positions
// @Description Returns all open positions for a user
// @Tags portfolio
// @Param user query string true "User ID"
// @Produce json
// @Success 200 {object} dto.PositionsListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /positions [get]
func (h *PortfolioHandler) GetPositions(c *gin.Context) {
	user := c.Query("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user parameter"})
		return
	}
	positions, err := h.PosRepo.GetPositionsByUser(context.Background(), user)
	if err != nil {
		utils.Logger.WithError(err).Error("failed to fetch positions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch positions"})
		return
	}
	var posDTOs []dto.PositionResponse
	for _, p := range positions {
		marketValue := p.Quantity * p.AveragePrice            // mock
		pnl := (p.AveragePrice - p.AveragePrice) * p.Quantity // mock
		posDTOs = append(posDTOs, dto.PositionResponse{
			Symbol:       p.Symbol,
			Quantity:     p.Quantity,
			AveragePrice: p.AveragePrice,
			MarketValue:  marketValue,
			PnL:          pnl,
		})
	}
	c.JSON(http.StatusOK, dto.PositionsListResponse{Positions: posDTOs})
}

// GetTrades godoc
// @Summary Get user trade history
// @Description Returns trade history for a user
// @Tags portfolio
// @Param user query string true "User ID"
// @Produce json
// @Success 200 {object} dto.TradesListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trades [get]
func (h *PortfolioHandler) GetTrades(c *gin.Context) {
	user := c.Query("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user parameter"})
		return
	}
	trades, err := h.TradeRepo.GetTradesByUser(context.Background(), user)
	if err != nil {
		utils.Logger.WithError(err).Error("failed to fetch trades")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch trades"})
		return
	}
	var tradeDTOs []dto.TradeResponse
	for _, t := range trades {
		// For now, side is not stored in trade, so mock as "buy"
		tradeDTOs = append(tradeDTOs, dto.TradeResponse{
			ID:        t.ID,
			Symbol:    t.Symbol,
			Side:      "buy", // mock
			Quantity:  t.Quantity,
			Price:     t.Price,
			Timestamp: t.Timestamp.Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, dto.TradesListResponse{Trades: tradeDTOs})
}
