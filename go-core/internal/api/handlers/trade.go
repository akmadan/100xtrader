package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-core/internal/api/dto"
	"go-core/internal/data"
	"go-core/internal/data/repos"
	"go-core/internal/services/brokers"
	"go-core/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateTrade creates a new trade
// @Summary Create a new trade
// @Description Create a new trade with all required and optional fields
// @Tags trades
// @Accept json
// @Produce json
// @Param trade body dto.CreateTradeRequest true "Trade data"
// @Success 201 {object} dto.SuccessResponse{data=dto.TradeResponse} "Trade created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/trades [post]
func CreateTrade(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateTradeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind trade request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Invalid JSON data",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Validate request
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			utils.LogError(err, "Validation failed for trade request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Parse entry date
		entryDate, err := time.Parse("2006-01-02", req.EntryDate)
		if err != nil {
			utils.LogError(err, "Failed to parse entry date")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Date",
				Message: "Entry date must be in YYYY-MM-DD format",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Convert DTO to model
		trade := &data.Trade{
			ID:             utils.GenerateID(),
			UserID:         req.UserID,
			Symbol:         req.Symbol,
			MarketType:     req.MarketType,
			EntryDate:      entryDate,
			EntryPrice:     req.EntryPrice,
			Quantity:       req.Quantity,
			TotalAmount:    req.TotalAmount,
			ExitPrice:      req.ExitPrice,
			Direction:      req.Direction,
			StopLoss:       req.StopLoss,
			Target:         req.Target,
			Strategy:       req.Strategy,
			OutcomeSummary: req.OutcomeSummary,
			TradeAnalysis:  req.TradeAnalysis,
			RulesFollowed:  req.RulesFollowed,
			Screenshots:    req.Screenshots,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Add psychology if provided
		if req.Psychology != nil {
			trade.Psychology = &data.TradePsychology{
				EntryConfidence:    req.Psychology.EntryConfidence,
				SatisfactionRating: req.Psychology.SatisfactionRating,
				EmotionalState:     req.Psychology.EmotionalState,
				MistakesMade:       req.Psychology.MistakesMade,
				LessonsLearned:     req.Psychology.LessonsLearned,
			}
		}

		// Create trade in database
		repo := repos.NewTradeRepository(db.GetConnection())
		if err := repo.CreateTrade(trade); err != nil {
			utils.LogError(err, "Failed to create trade")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to create trade",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Convert to response DTO
		response := convertTradeToResponse(trade)

		utils.LogInfo("Trade created successfully", map[string]interface{}{
			"trade_id": trade.ID,
			"user_id":  trade.UserID,
		})

		c.JSON(http.StatusCreated, dto.SuccessResponse{
			Message: "Trade created successfully",
			Data:    response,
		})
	}
}

// GetTrade retrieves a trade by ID
// @Summary Get a trade by ID
// @Description Retrieve a specific trade with all its details
// @Tags trades
// @Accept json
// @Produce json
// @Param id path string true "Trade ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.TradeResponse} "Trade retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid trade ID"
// @Failure 404 {object} dto.ErrorResponse "Trade not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/trades/{id} [get]
func GetTrade(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tradeID := c.Param("id")
		if tradeID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Trade ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// For now, we'll use 0 as userID since we don't have authentication yet
		// TODO: Get userID from JWT token or session
		userID := 0

		repo := repos.NewTradeRepository(db.GetConnection())
		trade, err := repo.GetTradeByID(tradeID, userID)
		if err != nil {
			utils.LogError(err, "Failed to get trade", map[string]interface{}{
				"trade_id": tradeID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Trade not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		response := convertTradeToResponse(trade)

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Trade retrieved successfully",
			Data:    response,
		})
	}
}

// UpdateTrade updates an existing trade
// @Summary Update a trade
// @Description Update an existing trade with new data
// @Tags trades
// @Accept json
// @Produce json
// @Param id path string true "Trade ID"
// @Param trade body dto.UpdateTradeRequest true "Updated trade data"
// @Success 200 {object} dto.SuccessResponse{data=dto.TradeResponse} "Trade updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 404 {object} dto.ErrorResponse "Trade not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/trades/{id} [put]
func UpdateTrade(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tradeID := c.Param("id")
		if tradeID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Trade ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		var req dto.UpdateTradeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind trade update request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Invalid JSON data",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Validate request
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			utils.LogError(err, "Validation failed for trade update request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Parse entry date
		entryDate, err := time.Parse("2006-01-02", req.EntryDate)
		if err != nil {
			utils.LogError(err, "Failed to parse entry date")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Date",
				Message: "Entry date must be in YYYY-MM-DD format",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Convert DTO to model
		trade := &data.Trade{
			ID:             tradeID,
			UserID:         req.UserID,
			Symbol:         req.Symbol,
			MarketType:     req.MarketType,
			EntryDate:      entryDate,
			EntryPrice:     req.EntryPrice,
			Quantity:       req.Quantity,
			TotalAmount:    req.TotalAmount,
			ExitPrice:      req.ExitPrice,
			Direction:      req.Direction,
			StopLoss:       req.StopLoss,
			Target:         req.Target,
			Strategy:       req.Strategy,
			OutcomeSummary: req.OutcomeSummary,
			TradeAnalysis:  req.TradeAnalysis,
			RulesFollowed:  req.RulesFollowed,
			Screenshots:    req.Screenshots,
			UpdatedAt:      time.Now(),
		}

		// Update trade in database
		repo := repos.NewTradeRepository(db.GetConnection())

		// Get existing trade to preserve psychology fields that aren't being updated
		existingTrade, err := repo.GetTradeByID(tradeID, req.UserID)
		if err != nil {
			utils.LogError(err, "Failed to get existing trade for update", map[string]interface{}{
				"trade_id": tradeID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Trade not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		// Handle psychology update - merge with existing if provided
		if req.Psychology != nil {
			// Start with existing psychology or create new
			var psychology *data.TradePsychology
			if existingTrade.Psychology != nil {
				psychology = &data.TradePsychology{
					EntryConfidence:    existingTrade.Psychology.EntryConfidence,
					SatisfactionRating: existingTrade.Psychology.SatisfactionRating,
					EmotionalState:     existingTrade.Psychology.EmotionalState,
					MistakesMade:       existingTrade.Psychology.MistakesMade,
					LessonsLearned:     existingTrade.Psychology.LessonsLearned,
				}
			} else {
				psychology = &data.TradePsychology{}
			}

			// Update only fields that are provided (non-nil)
			if req.Psychology.EntryConfidence != nil {
				psychology.EntryConfidence = *req.Psychology.EntryConfidence
			}
			if req.Psychology.SatisfactionRating != nil {
				psychology.SatisfactionRating = *req.Psychology.SatisfactionRating
			}
			if req.Psychology.EmotionalState != nil {
				psychology.EmotionalState = *req.Psychology.EmotionalState
			}
			if req.Psychology.MistakesMade != nil {
				psychology.MistakesMade = req.Psychology.MistakesMade
			}
			if req.Psychology.LessonsLearned != nil {
				psychology.LessonsLearned = req.Psychology.LessonsLearned
			}

			trade.Psychology = psychology
		} else if existingTrade.Psychology != nil {
			// Keep existing psychology if not provided in update
			trade.Psychology = existingTrade.Psychology
		}
		if err := repo.UpdateTrade(trade); err != nil {
			utils.LogError(err, "Failed to update trade", map[string]interface{}{
				"trade_id": tradeID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to update trade",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Get updated trade
		updatedTrade, err := repo.GetTradeByID(tradeID, req.UserID)
		if err != nil {
			utils.LogError(err, "Failed to get updated trade", map[string]interface{}{
				"trade_id": tradeID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve updated trade",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		response := convertTradeToResponse(updatedTrade)

		utils.LogInfo("Trade updated successfully", map[string]interface{}{
			"trade_id": tradeID,
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Trade updated successfully",
			Data:    response,
		})
	}
}

// DeleteTrade deletes a trade
// @Summary Delete a trade
// @Description Delete a specific trade by ID
// @Tags trades
// @Accept json
// @Produce json
// @Param id path string true "Trade ID"
// @Success 200 {object} dto.SuccessResponse "Trade deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid trade ID"
// @Failure 404 {object} dto.ErrorResponse "Trade not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/trades/{id} [delete]
func DeleteTrade(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tradeID := c.Param("id")
		if tradeID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Trade ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// For now, we'll use 0 as userID since we don't have authentication yet
		// TODO: Get userID from JWT token or session
		userID := 0

		repo := repos.NewTradeRepository(db.GetConnection())
		if err := repo.DeleteTrade(tradeID, userID); err != nil {
			utils.LogError(err, "Failed to delete trade", map[string]interface{}{
				"trade_id": tradeID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Trade not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		utils.LogInfo("Trade deleted successfully", map[string]interface{}{
			"trade_id": tradeID,
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Trade deleted successfully",
		})
	}
}

// ListTrades retrieves trades with pagination
// @Summary List trades
// @Description Retrieve a paginated list of trades
// @Tags trades
// @Accept json
// @Produce json
// @Param limit query int false "Number of trades to return (default: 10, max: 100)"
// @Param offset query int false "Number of trades to skip (default: 0)"
// @Success 200 {object} dto.SuccessResponse{data=dto.GetTradesResponse} "Trades retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/trades [get]
func ListTrades(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		limitStr := c.DefaultQuery("limit", "10")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			limit = 10
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0
		}

		repo := repos.NewTradeRepository(db.GetConnection())
		trades, err := repo.GetTradesByUser(0, limit, offset) // 0 means get all users' trades
		if err != nil {
			utils.LogError(err, "Failed to list trades")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve trades",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// For now, we'll use the length of trades as total since we don't have a separate count method
		total := len(trades)

		// Convert to response DTOs
		var tradeResponses []dto.TradeResponse
		for _, trade := range trades {
			tradeResponses = append(tradeResponses, convertTradeToResponse(trade))
		}

		response := dto.GetTradesResponse{
			Trades: tradeResponses,
			Pagination: dto.PaginationResponse{
				Total:  total,
				Limit:  limit,
				Offset: offset,
				Count:  len(tradeResponses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Trades retrieved successfully",
			Data:    response,
		})
	}
}

// GetTradesByUser retrieves trades for a specific user
// @Summary Get trades by user
// @Description Retrieve a paginated list of trades for a specific user
// @Tags trades
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param limit query int false "Number of trades to return (default: 10, max: 100)"
// @Param offset query int false "Number of trades to skip (default: 0)"
// @Success 200 {object} dto.SuccessResponse{data=dto.GetTradesResponse} "User trades retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid user ID or query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/trades/user/{user_id} [get]
func GetTradesByUser(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Invalid user ID",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Parse query parameters
		limitStr := c.DefaultQuery("limit", "10")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			limit = 10
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0
		}

		repo := repos.NewTradeRepository(db.GetConnection())
		trades, err := repo.GetTradesByUser(userID, limit, offset)
		if err != nil {
			utils.LogError(err, "Failed to get user trades", map[string]interface{}{
				"user_id": userID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve user trades",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// For now, we'll use the length of trades as total since we don't have a separate count method
		total := len(trades)

		// Convert to response DTOs
		var tradeResponses []dto.TradeResponse
		for _, trade := range trades {
			tradeResponses = append(tradeResponses, convertTradeToResponse(trade))
		}

		response := dto.GetTradesResponse{
			Trades: tradeResponses,
			Pagination: dto.PaginationResponse{
				Total:  total,
				Limit:  limit,
				Offset: offset,
				Count:  len(tradeResponses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User trades retrieved successfully",
			Data:    response,
		})
	}
}

// convertTradeToResponse converts a data.Trade to dto.TradeResponse
func convertTradeToResponse(trade *data.Trade) dto.TradeResponse {
	response := dto.TradeResponse{
		ID:             trade.ID,
		UserID:         trade.UserID,
		Symbol:         trade.Symbol,
		MarketType:     trade.MarketType,
		EntryDate:      trade.EntryDate,
		EntryPrice:     trade.EntryPrice,
		Quantity:       trade.Quantity,
		TotalAmount:    trade.TotalAmount,
		ExitPrice:      trade.ExitPrice,
		Direction:      trade.Direction,
		StopLoss:       trade.StopLoss,
		Target:         trade.Target,
		Strategy:       trade.Strategy,
		OutcomeSummary: trade.OutcomeSummary,
		TradeAnalysis:  trade.TradeAnalysis,
		RulesFollowed:  trade.RulesFollowed,
		Screenshots:    trade.Screenshots,
		CreatedAt:      trade.CreatedAt,
		UpdatedAt:      trade.UpdatedAt,
	}

	// Add psychology if present
	if trade.Psychology != nil {
		response.Psychology = &dto.PsychologyResponse{
			EntryConfidence:    trade.Psychology.EntryConfidence,
			SatisfactionRating: trade.Psychology.SatisfactionRating,
			EmotionalState:     trade.Psychology.EmotionalState,
			MistakesMade:       trade.Psychology.MistakesMade,
			LessonsLearned:     trade.Psychology.LessonsLearned,
		}
	}

	return response
}

// SyncDhanTrades syncs trades from Dhan broker for a user
// @Summary Sync Dhan trades
// @Description Fetches all trades from Dhan API since Jan 2022 and saves new trades to the database
// @Tags trades
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.SuccessResponse "Trades synced successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id}/trades/sync-dhan [post]
func SyncDhanTrades(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Invalid user ID",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Get user's Dhan broker configuration
		userRepo := repos.NewUserRepository(db.GetConnection())
		user, err := userRepo.GetUserByID(userIDStr)
		if err != nil {
			utils.LogError(err, "Failed to get user for sync")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to get user",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		config, exists := user.ConfiguredBrokers["dhan"]
		if !exists || config.AccessToken == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "Dhan broker not configured or access token missing. Please configure Dhan first.",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Get the latest trade date for this user from Dhan broker
		tradeRepo := repos.NewTradeRepository(db.GetConnection())
		latestTradeDate, err := tradeRepo.GetLatestTradeDateByBroker(userID, data.TradingBrokerDhan)
		if err != nil {
			utils.LogError(err, "Failed to get latest trade date")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to get latest trade date",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Determine start date: use latest trade date if exists, otherwise start from Jan 1, 1999
		var fromDate time.Time
		if latestTradeDate != nil {
			// Start from the day after the latest trade date to avoid duplicates
			fromDate = latestTradeDate.AddDate(0, 0, 1)
		} else {
			// First sync: start from Jan 1, 1999 to get all historical trades
			fromDate = time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)
		}

		// End date: current date
		toDate := time.Now().UTC()

		// Format dates for API
		fromDateStr := brokers.FormatDateForAPI(fromDate)
		toDateStr := brokers.FormatDateForAPI(toDate)

		// Check if token is expired
		if config.ExpiryTime != nil && config.ExpiryTime.Before(time.Now()) {
			utils.LogError(nil, "Dhan access token expired", map[string]interface{}{
				"user_id":     userID,
				"expiry_time": config.ExpiryTime,
			})
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "Dhan access token has expired. Please renew the token from the Accounts page.",
				Code:    http.StatusBadRequest,
			})
			return
		}

		utils.LogInfo("Starting Dhan trade sync", map[string]interface{}{
			"user_id":           userID,
			"from_date":         fromDateStr,
			"to_date":           toDateStr,
			"from_date_iso":     fromDate.Format(time.RFC3339),
			"to_date_iso":       toDate.Format(time.RFC3339),
			"has_token":         config.AccessToken != "",
			"token_length":      len(config.AccessToken),
			"expiry_time":       config.ExpiryTime,
			"latest_trade_date": latestTradeDate,
		})

		// Fetch trades from Dhan API
		dhanService := brokers.NewDhanService()
		trades, err := dhanService.FetchAndConvertTrades(config.AccessToken, fromDateStr, toDateStr, userID)

		// If we get an "Invalid Token" error, try to auto-renew the token
		if err != nil && strings.Contains(err.Error(), "Invalid Token") {
			utils.LogInfo("Token invalid, attempting auto-renewal", map[string]interface{}{
				"user_id": userID,
			})

			// Try to renew the token automatically
			if config.DhanClientID != nil && *config.DhanClientID != "" {
				renewResponse, renewErr := dhanService.RenewToken(config.AccessToken, *config.DhanClientID)
				if renewErr != nil {
					utils.LogError(renewErr, "Failed to auto-renew token", map[string]interface{}{
						"user_id": userID,
					})
				} else {
					// Parse expiry time
					expiryTime, parseErr := time.Parse("2006-01-02T15:04:05", renewResponse.ExpiryTime)
					if parseErr != nil {
						expiryTime = time.Now().Add(24 * time.Hour)
					}

					// Update config with new token
					config.AccessToken = renewResponse.AccessToken
					config.ExpiryTime = &expiryTime
					config.ConfiguredAt = time.Now()
					user.ConfiguredBrokers["dhan"] = config

					// Save updated config
					if updateErr := userRepo.UpdateUser(user); updateErr != nil {
						utils.LogError(updateErr, "Failed to save renewed token", map[string]interface{}{
							"user_id": userID,
						})
					} else {
						utils.LogInfo("Token auto-renewed successfully", map[string]interface{}{
							"user_id": userID,
						})

						// Retry fetching trades with new token
						trades, err = dhanService.FetchAndConvertTrades(config.AccessToken, fromDateStr, toDateStr, userID)
						if err == nil {
							utils.LogInfo("Successfully fetched trades after token renewal", map[string]interface{}{
								"user_id": userID,
								"count":   len(trades),
							})
						} else {
							utils.LogError(err, "Failed to fetch trades after token renewal", map[string]interface{}{
								"user_id": userID,
							})
						}
					}
				}
			} else {
				utils.LogError(nil, "Cannot auto-renew token: Dhan Client ID missing", map[string]interface{}{
					"user_id": userID,
				})
			}
		}

		if err != nil {
			utils.LogError(err, "Failed to fetch trades from Dhan")
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "Invalid Token") {
				c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error:   "Unauthorized",
					Message: "Dhan access token is invalid or expired. Please renew the token from the Accounts page.",
					Code:    http.StatusUnauthorized,
				})
			} else {
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error:   "Internal Server Error",
					Message: "Failed to fetch trades from Dhan: " + errorMsg,
					Code:    http.StatusInternalServerError,
				})
			}
			return
		}

		utils.LogInfo("Fetched trades from Dhan", map[string]interface{}{
			"user_id": userID,
			"count":   len(trades),
		})

		// Save trades that don't already exist
		savedCount := 0
		skippedCount := 0
		errorCount := 0

		for _, trade := range trades {
			// Check if trade already exists by broker ID
			exchangeOrderID := ""
			orderID := ""
			if trade.ExchangeOrderID != nil {
				exchangeOrderID = *trade.ExchangeOrderID
			}
			if trade.OrderID != nil {
				orderID = *trade.OrderID
			}

			exists, err := tradeRepo.TradeExistsByBrokerID(userID, data.TradingBrokerDhan, exchangeOrderID, orderID)
			if err != nil {
				utils.LogError(err, "Failed to check if trade exists", map[string]interface{}{
					"exchange_order_id": exchangeOrderID,
					"order_id":          orderID,
				})
				errorCount++
				continue
			}

			if exists {
				skippedCount++
				continue
			}

			// Save new trade
			if err := tradeRepo.CreateTrade(trade); err != nil {
				utils.LogError(err, "Failed to save trade", map[string]interface{}{
					"exchange_order_id": exchangeOrderID,
					"order_id":          orderID,
				})
				errorCount++
				continue
			}

			savedCount++
		}

		utils.LogInfo("Completed Dhan trade sync", map[string]interface{}{
			"user_id":       userID,
			"saved_count":   savedCount,
			"skipped_count": skippedCount,
			"error_count":   errorCount,
			"total_fetched": len(trades),
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: fmt.Sprintf("Sync completed. Saved %d new trades, skipped %d existing trades", savedCount, skippedCount),
			Data: map[string]interface{}{
				"saved_count":   savedCount,
				"skipped_count": skippedCount,
				"error_count":   errorCount,
				"total_fetched": len(trades),
				"from_date":     fromDateStr,
				"to_date":       toDateStr,
				"date_range":    fmt.Sprintf("%s to %s", fromDateStr, toDateStr),
			},
		})
	}
}
