package api

import (
	"context"
	"net/http"
	"time"

	"database/sql"
	"strconv"

	"github.com/akshitmadan/100xtrader/backend/internal/api/dto"
	"github.com/akshitmadan/100xtrader/backend/internal/data"
	"github.com/akshitmadan/100xtrader/backend/internal/data/repos"
	"github.com/akshitmadan/100xtrader/backend/internal/engine"
	"github.com/akshitmadan/100xtrader/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

// OrderHandler holds dependencies for order endpoints

type OrderHandler struct {
	Repo      repos.OrderRepository
	TradeRepo repos.TradeRepository
	PosRepo   repos.PositionRepository
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(repo repos.OrderRepository, tradeRepo repos.TradeRepository, posRepo repos.PositionRepository) *OrderHandler {
	return &OrderHandler{
		Repo:      repo,
		TradeRepo: tradeRepo,
		PosRepo:   posRepo,
	}
}

// RegisterOrderRoutes registers order-related routes with dependency injection
func RegisterOrderRoutes(r *gin.Engine, handler *OrderHandler) {
	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders", handler.GetOrders)
	r.DELETE("/orders/:id", handler.DeleteOrder) // TODO: implement
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Places a new order (market, limit, or stop)
// @Tags orders
// @Accept json
// @Produce json
// @Param order body dto.OrderCreateRequest true "Order create request"
// @Success 201 {object} dto.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.WithError(err).Warn("invalid order create request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	order := &data.Order{
		User:      req.User,
		Symbol:    req.Symbol,
		Side:      req.Side,
		Type:      req.Type,
		Quantity:  req.Quantity,
		Price:     req.Price,
		Status:    "open",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	id, err := h.Repo.CreateOrder(context.Background(), order)
	if err != nil {
		utils.Logger.WithError(err).Error("failed to create order")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}
	order.ID = id

	var trades []*data.Trade
	if order.Type == "market" {
		trades = engine.GetOrderBookManager().AddMarketOrderAndMatch(order)
	} else if order.Type == "stop" {
		// Add to stop order list, do not match immediately
		engine.GetOrderBookManager().GetOrCreateOrderBook(order.Symbol).AddStopOrder(order)
		trades = nil
	} else {
		trades = engine.GetOrderBookManager().AddOrderAndMatch(order)
	}

	for _, trade := range trades {
		if err := h.TradeRepo.RecordTrade(context.Background(), trade); err != nil {
			utils.Logger.WithError(err).Error("failed to record trade")
		}
		// Update position for buyer (the user placing the order)
		if err := h.PosRepo.UpdatePosition(context.Background(), order.User, order.Symbol, trade.Quantity, trade.Price); err != nil {
			utils.Logger.WithError(err).Error("failed to update position for buyer")
		}
		// Update position for seller (fetch user from sell order)
		sellOrder, err := h.Repo.GetOrderByID(context.Background(), trade.SellOrderID)
		if err != nil {
			utils.Logger.WithError(err).Error("failed to fetch sell order for seller position update")
			continue
		}
		if err := h.PosRepo.UpdatePosition(context.Background(), sellOrder.User, sellOrder.Symbol, -trade.Quantity, trade.Price); err != nil {
			utils.Logger.WithError(err).Error("failed to update position for seller")
		}
		utils.Logger.WithField("trade", trade).Info("trade executed")

		// After each trade, check and trigger stop orders for this symbol
		triggered := engine.GetOrderBookManager().GetOrCreateOrderBook(order.Symbol).CheckAndTriggerStopOrders(trade.Price)
		for _, stopOrder := range triggered {
			// Process triggered stop order as a market order
			stopTrades := engine.GetOrderBookManager().AddMarketOrderAndMatch(stopOrder)
			for _, stopTrade := range stopTrades {
				if err := h.TradeRepo.RecordTrade(context.Background(), stopTrade); err != nil {
					utils.Logger.WithError(err).Error("failed to record triggered stop trade")
				}
				if err := h.PosRepo.UpdatePosition(context.Background(), stopOrder.User, stopOrder.Symbol, stopTrade.Quantity, stopTrade.Price); err != nil {
					utils.Logger.WithError(err).Error("failed to update position for triggered stop buyer")
				}
				sellOrder, err := h.Repo.GetOrderByID(context.Background(), stopTrade.SellOrderID)
				if err != nil {
					utils.Logger.WithError(err).Error("failed to fetch sell order for triggered stop seller position update")
					continue
				}
				if err := h.PosRepo.UpdatePosition(context.Background(), sellOrder.User, sellOrder.Symbol, -stopTrade.Quantity, stopTrade.Price); err != nil {
					utils.Logger.WithError(err).Error("failed to update position for triggered stop seller")
				}
				utils.Logger.WithField("trade", stopTrade).Info("triggered stop trade executed")
			}
		}
	}

	resp := dto.OrderResponse{
		ID:        order.ID,
		User:      order.User,
		Symbol:    order.Symbol,
		Side:      order.Side,
		Type:      order.Type,
		Quantity:  order.Quantity,
		Price:     order.Price,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
		UpdatedAt: order.UpdatedAt.Format(time.RFC3339),
	}
	c.JSON(http.StatusCreated, resp)
}

// GetOrders godoc
// @Summary List all orders
// @Description Returns all orders (optionally filter by user)
// @Tags orders
// @Produce json
// @Success 200 {array} dto.OrderResponse
// @Failure 500 {object} map[string]string
// @Router /orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.Repo.ListOrders(context.Background())
	if err != nil {
		utils.Logger.WithError(err).Error("failed to list orders")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list orders"})
		return
	}
	var resp []dto.OrderResponse
	for _, order := range orders {
		resp = append(resp, dto.OrderResponse{
			ID:        order.ID,
			User:      order.User,
			Symbol:    order.Symbol,
			Side:      order.Side,
			Type:      order.Type,
			Quantity:  order.Quantity,
			Price:     order.Price,
			Status:    order.Status,
			CreatedAt: order.CreatedAt.Format(time.RFC3339),
			UpdatedAt: order.UpdatedAt.Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteOrder godoc
// @Summary Delete an order
// @Description Deletes an order by ID
// @Tags orders
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Logger.WithError(err).Warn("invalid order id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	err = h.Repo.DeleteOrder(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		utils.Logger.WithError(err).Error("failed to delete order")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
