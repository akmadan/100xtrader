package handlers

import (
	"net/http"
	"strconv"
	"time"

	"go-core/internal/api/dto"
	"go-core/internal/data"
	"go-core/internal/data/repos"
	"go-core/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateTrade creates a new trade with optional journal and actions
// @Summary Create a new trade
// @Description Create a new trade with optional journal entry, trade actions, and tags
// @Tags trades
// @Accept json
// @Produce json
// @Param trade body dto.TradeCreateRequest true "Trade creation data"
// @Success 201 {object} dto.SuccessResponse{data=dto.TradeResponse} "Trade created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/trades [post]
func CreateTrade(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TradeCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind trade creation request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid request",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Validate request
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			utils.LogError(err, "Trade creation validation failed")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation failed",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Generate trade ID
		tradeID := generateTradeID()

		// Create trade
		trade := &data.Trade{
			ID:       tradeID,
			UserID:   req.UserID,
			Market:   req.Market,
			Symbol:   req.Symbol,
			Target:   req.Target,
			StopLoss: req.StopLoss,
		}

		tradeRepo := repos.NewTradeRepository(db.GetConnection())
		if err := tradeRepo.Create(trade); err != nil {
			utils.LogError(err, "Failed to create trade", map[string]interface{}{
				"trade_id": tradeID,
				"user_id":  req.UserID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Failed to create trade",
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		utils.LogTrade("CREATE", tradeID, req.Symbol, req.UserID, map[string]interface{}{
			"market": req.Market,
		})

		// Create journal if provided
		if req.Journal != nil {
			journalRepo := repos.NewTradeJournalRepository(db.GetConnection())
			journal := &data.TradeJournal{
				TradeID:    tradeID,
				Notes:      req.Journal.Notes,
				Confidence: req.Journal.Confidence,
			}
			if err := journalRepo.Create(journal); err != nil {
				utils.LogError(err, "Failed to create trade journal", map[string]interface{}{
					"trade_id": tradeID,
				})
				// Don't fail the entire request, just log the error
			}
		}

		// Create trade actions if provided
		if len(req.Actions) > 0 {
			actionRepo := repos.NewTradeActionRepository(db.GetConnection())
			for _, actionReq := range req.Actions {
				action := &data.TradeAction{
					TradeID:   tradeID,
					Action:    actionReq.Action,
					TradeTime: actionReq.TradeTime,
					Quantity:  actionReq.Quantity,
					Price:     actionReq.Price,
					Fee:       actionReq.Fee,
				}
				if err := actionRepo.Create(action); err != nil {
					utils.LogError(err, "Failed to create trade action", map[string]interface{}{
						"trade_id": tradeID,
						"action":   actionReq.Action,
					})
					// Don't fail the entire request, just log the error
				}
			}
		}

		// Add tags if provided
		if len(req.Tags) > 0 {
			tagRepo := repos.NewTagRepository(db.GetConnection())
			for _, tagName := range req.Tags {
				// Create or get tag
				tag, err := tagRepo.GetOrCreate(tagName)
				if err != nil {
					utils.LogError(err, "Failed to get/create tag", map[string]interface{}{
						"tag_name": tagName,
					})
					continue
				}
				// Associate tag with trade
				if err := tagRepo.AddToTrade(tradeID, tag.ID); err != nil {
					utils.LogError(err, "Failed to add tag to trade", map[string]interface{}{
						"trade_id": tradeID,
						"tag_id":   tag.ID,
					})
				}
			}
		}

		// Return the created trade with all related data
		response := dto.TradeResponse{
			ID:        trade.ID,
			UserID:    trade.UserID,
			Market:    trade.Market,
			Symbol:    trade.Symbol,
			Target:    trade.Target,
			StopLoss:  trade.StopLoss,
			CreatedAt: trade.CreatedAt,
		}

		c.JSON(http.StatusCreated, dto.SuccessResponse{
			Message: "Trade created successfully",
			Data:    response,
		})
	}
}

// GetTrade retrieves a trade with its journal, actions, and tags
// @Summary Get trade by ID
// @Description Retrieve a specific trade with its journal, actions, tags, and screenshots
// @Tags trades
// @Accept json
// @Produce json
// @Param id path string true "Trade ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.TradeResponse} "Trade retrieved successfully"
// @Failure 404 {object} dto.ErrorResponse "Trade not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/trades/{id} [get]
func GetTrade(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tradeID := c.Param("id")

		tradeRepo := repos.NewTradeRepository(db.GetConnection())
		trade, err := tradeRepo.GetByID(tradeID)
		if err != nil {
			utils.LogError(err, "Failed to get trade", map[string]interface{}{
				"trade_id": tradeID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Trade not found",
				Message: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}

		// Get journal if exists
		journalRepo := repos.NewTradeJournalRepository(db.GetConnection())
		journal, _ := journalRepo.GetByTradeID(tradeID)

		// Get actions
		actionRepo := repos.NewTradeActionRepository(db.GetConnection())
		actions, _ := actionRepo.GetByTradeID(tradeID)

		// Get tags
		tagRepo := repos.NewTagRepository(db.GetConnection())
		tags, _ := tagRepo.GetByTradeID(tradeID)

		// Get screenshots if journal exists
		var screenshots []dto.TradeScreenshotResponse
		if journal != nil {
			screenshotRepo := repos.NewTradeScreenshotRepository(db.GetConnection())
			screenshotData, _ := screenshotRepo.GetByJournalID(journal.ID)
			for _, screenshot := range screenshotData {
				screenshots = append(screenshots, dto.TradeScreenshotResponse{
					ID:             screenshot.ID,
					TradeJournalID: screenshot.TradeJournalID,
					URL:            screenshot.URL,
					CreatedAt:      screenshot.CreatedAt,
				})
			}
		}

		// Build response
		response := dto.TradeResponse{
			ID:        trade.ID,
			UserID:    trade.UserID,
			Market:    trade.Market,
			Symbol:    trade.Symbol,
			Target:    trade.Target,
			StopLoss:  trade.StopLoss,
			CreatedAt: trade.CreatedAt,
		}

		if journal != nil {
			response.Journal = &dto.TradeJournalResponse{
				ID:         journal.ID,
				TradeID:    journal.TradeID,
				Notes:      journal.Notes,
				Confidence: journal.Confidence,
				CreatedAt:  journal.CreatedAt,
			}
		}

		// Convert actions
		for _, action := range actions {
			response.Actions = append(response.Actions, dto.TradeActionResponse{
				ID:        action.ID,
				TradeID:   action.TradeID,
				Action:    action.Action,
				TradeTime: action.TradeTime,
				Quantity:  action.Quantity,
				Price:     action.Price,
				Fee:       action.Fee,
				CreatedAt: action.CreatedAt,
			})
		}

		// Convert tags
		for _, tag := range tags {
			response.Tags = append(response.Tags, dto.TagResponse{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
			})
		}

		response.Screenshots = screenshots

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Trade retrieved successfully",
			Data:    response,
		})
	}
}

// ListTrades retrieves all trades with optional filtering
// @Summary List trades
// @Description Retrieve a list of trades with optional filtering by user_id, market, and include_journal
// @Tags trades
// @Accept json
// @Produce json
// @Param user_id query int false "Filter by user ID"
// @Param market query string false "Filter by market type"
// @Param include_journal query boolean false "Include journal data in response"
// @Success 200 {object} dto.TradeListResponse "Trades retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/trades [get]
func ListTrades(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		userIDStr := c.Query("user_id")
		market := c.Query("market")
		includeJournal := c.Query("include_journal") == "true"

		tradeRepo := repos.NewTradeRepository(db.GetConnection())
		var trades []*data.Trade
		var err error

		if userIDStr != "" {
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "Invalid user_id",
					Message: "user_id must be a valid integer",
					Code:    http.StatusBadRequest,
				})
				return
			}
			trades, err = tradeRepo.GetByUserID(userID)
		} else if market != "" {
			trades, err = tradeRepo.GetByMarket(market)
		} else {
			trades, err = tradeRepo.List()
		}

		if err != nil {
			utils.LogError(err, "Failed to list trades")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Failed to retrieve trades",
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Convert to response format
		var responses []dto.TradeResponse
		for _, trade := range trades {
			response := dto.TradeResponse{
				ID:        trade.ID,
				UserID:    trade.UserID,
				Market:    trade.Market,
				Symbol:    trade.Symbol,
				Target:    trade.Target,
				StopLoss:  trade.StopLoss,
				CreatedAt: trade.CreatedAt,
			}

			// Include journal if requested
			if includeJournal {
				journalRepo := repos.NewTradeJournalRepository(db.GetConnection())
				journal, _ := journalRepo.GetByTradeID(trade.ID)
				if journal != nil {
					response.Journal = &dto.TradeJournalResponse{
						ID:         journal.ID,
						TradeID:    journal.TradeID,
						Notes:      journal.Notes,
						Confidence: journal.Confidence,
						CreatedAt:  journal.CreatedAt,
					}
				}
			}

			responses = append(responses, response)
		}

		c.JSON(http.StatusOK, dto.TradeListResponse{
			Trades: responses,
			Total:  len(responses),
		})
	}
}

// generateTradeID generates a unique trade ID
func generateTradeID() string {
	return "TRADE-" + time.Now().Format("20060102150405") + "-" + randomString(4)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
