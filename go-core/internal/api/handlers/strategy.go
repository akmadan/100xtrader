package handlers

import (
	"net/http"
	"strconv"

	"go-core/internal/api/dto"
	"go-core/internal/data"
	"go-core/internal/data/repos"
	"go-core/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateStrategy creates a new strategy
// @Summary Create a new strategy
// @Description Create a new trading strategy
// @Tags strategies
// @Accept json
// @Produce json
// @Param strategy body dto.CreateStrategyRequest true "Strategy data"
// @Success 201 {object} dto.SuccessResponse{data=dto.StrategyResponse} "Strategy created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/strategies [post]
func CreateStrategy(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateStrategyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind strategy request")
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
			utils.LogError(err, "Validation failed for strategy request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Convert DTO to model
		strategy := &data.Strategy{
			ID:          utils.GenerateID(),
			UserID:      req.UserID,
			Name:        req.Name,
			Description: req.Description,
			CreatedAt:   utils.GetCurrentTime(),
			UpdatedAt:   utils.GetCurrentTime(),
		}

		// Create strategy in database
		repo := repos.NewStrategyRepository(db.GetConnection())
		if err := repo.CreateStrategy(strategy); err != nil {
			utils.LogError(err, "Failed to create strategy")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to create strategy",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Convert to response DTO
		response := convertStrategyToResponse(strategy)

		utils.LogInfo("Strategy created successfully", map[string]interface{}{
			"strategy_id": strategy.ID,
			"user_id":     strategy.UserID,
		})

		c.JSON(http.StatusCreated, dto.SuccessResponse{
			Message: "Strategy created successfully",
			Data:    response,
		})
	}
}

// GetStrategy retrieves a strategy by ID
// @Summary Get a strategy by ID
// @Description Retrieve a specific strategy with all its details
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "Strategy ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.StrategyResponse} "Strategy retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid strategy ID"
// @Failure 404 {object} dto.ErrorResponse "Strategy not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/strategies/{id} [get]
func GetStrategy(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strategyID := c.Param("id")
		if strategyID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Strategy ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// For now, we'll use 0 as userID since we don't have authentication yet
		// TODO: Get userID from JWT token or session
		userID := 0

		repo := repos.NewStrategyRepository(db.GetConnection())
		strategy, err := repo.GetStrategyByID(strategyID, userID)
		if err != nil {
			utils.LogError(err, "Failed to get strategy", map[string]interface{}{
				"strategy_id": strategyID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Strategy not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		response := convertStrategyToResponse(strategy)

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Strategy retrieved successfully",
			Data:    response,
		})
	}
}

// UpdateStrategy updates an existing strategy
// @Summary Update a strategy
// @Description Update an existing strategy with new data
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "Strategy ID"
// @Param strategy body dto.UpdateStrategyRequest true "Updated strategy data"
// @Success 200 {object} dto.SuccessResponse{data=dto.StrategyResponse} "Strategy updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 404 {object} dto.ErrorResponse "Strategy not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/strategies/{id} [put]
func UpdateStrategy(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strategyID := c.Param("id")
		if strategyID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Strategy ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		var req dto.UpdateStrategyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind strategy update request")
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
			utils.LogError(err, "Validation failed for strategy update request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Convert DTO to model
		strategy := &data.Strategy{
			ID:          strategyID,
			UserID:      req.UserID,
			Name:        req.Name,
			Description: req.Description,
			UpdatedAt:   utils.GetCurrentTime(),
		}

		// Update strategy in database
		repo := repos.NewStrategyRepository(db.GetConnection())
		if err := repo.UpdateStrategy(strategy); err != nil {
			utils.LogError(err, "Failed to update strategy", map[string]interface{}{
				"strategy_id": strategyID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to update strategy",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Get updated strategy
		updatedStrategy, err := repo.GetStrategyByID(strategyID, req.UserID)
		if err != nil {
			utils.LogError(err, "Failed to get updated strategy", map[string]interface{}{
				"strategy_id": strategyID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve updated strategy",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		response := convertStrategyToResponse(updatedStrategy)

		utils.LogInfo("Strategy updated successfully", map[string]interface{}{
			"strategy_id": strategyID,
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Strategy updated successfully",
			Data:    response,
		})
	}
}

// DeleteStrategy deletes a strategy
// @Summary Delete a strategy
// @Description Delete a specific strategy by ID
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "Strategy ID"
// @Success 200 {object} dto.SuccessResponse "Strategy deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid strategy ID"
// @Failure 404 {object} dto.ErrorResponse "Strategy not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/strategies/{id} [delete]
func DeleteStrategy(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strategyID := c.Param("id")
		if strategyID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Strategy ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// For now, we'll use 0 as userID since we don't have authentication yet
		// TODO: Get userID from JWT token or session
		userID := 0

		repo := repos.NewStrategyRepository(db.GetConnection())
		if err := repo.DeleteStrategy(strategyID, userID); err != nil {
			utils.LogError(err, "Failed to delete strategy", map[string]interface{}{
				"strategy_id": strategyID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Strategy not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		utils.LogInfo("Strategy deleted successfully", map[string]interface{}{
			"strategy_id": strategyID,
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Strategy deleted successfully",
		})
	}
}

// ListStrategies retrieves strategies with pagination
// @Summary List strategies
// @Description Retrieve a paginated list of strategies
// @Tags strategies
// @Accept json
// @Produce json
// @Param limit query int false "Number of strategies to return (default: 10, max: 100)"
// @Param offset query int false "Number of strategies to skip (default: 0)"
// @Success 200 {object} dto.SuccessResponse{data=dto.GetStrategiesResponse} "Strategies retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/strategies [get]
func ListStrategies(db *data.DB) gin.HandlerFunc {
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

		repo := repos.NewStrategyRepository(db.GetConnection())
		strategies, err := repo.GetStrategiesByUser(0, limit, offset) // 0 means get all users' strategies
		if err != nil {
			utils.LogError(err, "Failed to list strategies")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve strategies",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// For now, we'll use the length of strategies as total since we don't have a separate count method
		total := len(strategies)

		// Convert to response DTOs
		var strategyResponses []dto.StrategyResponse
		for _, strategy := range strategies {
			strategyResponses = append(strategyResponses, convertStrategyToResponse(strategy))
		}

		response := dto.GetStrategiesResponse{
			Strategies: strategyResponses,
			Pagination: dto.PaginationResponse{
				Total:  total,
				Limit:  limit,
				Offset: offset,
				Count:  len(strategyResponses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Strategies retrieved successfully",
			Data:    response,
		})
	}
}

// GetStrategiesByUser retrieves strategies for a specific user
// @Summary Get strategies by user
// @Description Retrieve a paginated list of strategies for a specific user
// @Tags strategies
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param limit query int false "Number of strategies to return (default: 10, max: 100)"
// @Param offset query int false "Number of strategies to skip (default: 0)"
// @Success 200 {object} dto.SuccessResponse{data=dto.GetStrategiesResponse} "User strategies retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid user ID or query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/strategies/user/{user_id} [get]
func GetStrategiesByUser(db *data.DB) gin.HandlerFunc {
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

		repo := repos.NewStrategyRepository(db.GetConnection())
		strategies, err := repo.GetStrategiesByUser(userID, limit, offset)
		if err != nil {
			utils.LogError(err, "Failed to get user strategies", map[string]interface{}{
				"user_id": userID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve user strategies",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// For now, we'll use the length of strategies as total since we don't have a separate count method
		total := len(strategies)

		// Convert to response DTOs
		var strategyResponses []dto.StrategyResponse
		for _, strategy := range strategies {
			strategyResponses = append(strategyResponses, convertStrategyToResponse(strategy))
		}

		response := dto.GetStrategiesResponse{
			Strategies: strategyResponses,
			Pagination: dto.PaginationResponse{
				Total:  total,
				Limit:  limit,
				Offset: offset,
				Count:  len(strategyResponses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User strategies retrieved successfully",
			Data:    response,
		})
	}
}

// convertStrategyToResponse converts a data.Strategy to dto.StrategyResponse
func convertStrategyToResponse(strategy *data.Strategy) dto.StrategyResponse {
	return dto.StrategyResponse{
		ID:          strategy.ID,
		UserID:      strategy.UserID,
		Name:        strategy.Name,
		Description: strategy.Description,
		CreatedAt:   strategy.CreatedAt,
		UpdatedAt:   strategy.UpdatedAt,
	}
}
