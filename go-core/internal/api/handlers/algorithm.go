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

// CreateAlgorithm creates a new algorithm
func CreateAlgorithm(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateAlgorithmRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind algorithm request")
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
			utils.LogError(err, "Validation failed for algorithm request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Create algorithm
		algo := &data.Algorithm{
			ID:            utils.GenerateID(),
			UserID:        req.UserID,
			Name:          req.Name,
			Description:   req.Description,
			Code:          req.Code,
			Status:        req.Status,
			Symbol:        req.Symbol,
			Timeframe:     req.Timeframe,
			ExecutionMode: req.ExecutionMode,
			Broker:        req.Broker,
			Enabled:       req.Enabled,
			Config:        req.Config,
			State:         make(map[string]interface{}),
			TotalTrades:   0,
			WinRate:       0.0,
			TotalPnL:      0.0,
			Version:       1,
			Tags:          req.Tags,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if algo.Config == nil {
			algo.Config = make(map[string]interface{})
		}
		if algo.Tags == nil {
			algo.Tags = []string{}
		}

		repo := repos.NewAlgorithmRepository(db.GetConnection())
		if err := repo.CreateAlgorithm(algo); err != nil {
			utils.LogError(err, "Failed to create algorithm")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to create algorithm",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		response := convertAlgorithmToResponse(algo)
		c.JSON(http.StatusCreated, dto.SuccessResponse{
			Message: "Algorithm created successfully",
			Data:    response,
		})
	}
}

// GetAlgorithm retrieves an algorithm by ID
func GetAlgorithm(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Algorithm ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "User ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Invalid user ID",
				Code:    http.StatusBadRequest,
			})
			return
		}

		repo := repos.NewAlgorithmRepository(db.GetConnection())
		algo, err := repo.GetAlgorithmByID(id, userID)
		if err != nil {
			utils.LogError(err, "Failed to get algorithm")
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Algorithm not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		response := convertAlgorithmToResponse(algo)
		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Algorithm retrieved successfully",
			Data:    response,
		})
	}
}

// GetAlgorithmsByUser retrieves all algorithms for a user
func GetAlgorithmsByUser(db *data.DB) gin.HandlerFunc {
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

		repo := repos.NewAlgorithmRepository(db.GetConnection())
		algorithms, err := repo.GetAlgorithmsByUser(userID, limit, offset)
		if err != nil {
			utils.LogError(err, "Failed to get algorithms")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve algorithms",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		var responses []dto.AlgorithmResponse
		for _, algo := range algorithms {
			responses = append(responses, convertAlgorithmToResponse(algo))
		}

		response := dto.GetAlgorithmsResponse{
			Algorithms: responses,
			Pagination: dto.PaginationResponse{
				Total:  len(responses),
				Limit:  limit,
				Offset: offset,
				Count:  len(responses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Algorithms retrieved successfully",
			Data:    response,
		})
	}
}

// UpdateAlgorithm updates an existing algorithm
func UpdateAlgorithm(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Algorithm ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		var req dto.UpdateAlgorithmRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind update request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Invalid JSON data",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Get user_id from query or body
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			// Try to get from body if not in query
			var body map[string]interface{}
			if err := c.ShouldBindJSON(&body); err == nil {
				if uid, ok := body["user_id"].(float64); ok {
					userIDStr = strconv.Itoa(int(uid))
				}
			}
		}

		if userIDStr == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "User ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Invalid user ID",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Get existing algorithm
		repo := repos.NewAlgorithmRepository(db.GetConnection())
		existing, err := repo.GetAlgorithmByID(id, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Algorithm not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		// Update fields if provided
		if req.Name != nil {
			existing.Name = *req.Name
		}
		if req.Description != nil {
			existing.Description = req.Description
		}
		if req.Code != nil {
			existing.Code = *req.Code
			existing.Version++ // Increment version on code change
		}
		if req.Status != nil {
			existing.Status = *req.Status
		}
		if req.Symbol != nil {
			existing.Symbol = *req.Symbol
		}
		if req.Timeframe != nil {
			existing.Timeframe = *req.Timeframe
		}
		if req.ExecutionMode != nil {
			existing.ExecutionMode = *req.ExecutionMode
		}
		if req.Broker != nil {
			existing.Broker = req.Broker
		}
		if req.Enabled != nil {
			existing.Enabled = *req.Enabled
		}
		if req.Config != nil {
			existing.Config = req.Config
		}
		if req.Tags != nil {
			existing.Tags = req.Tags
		}

		existing.UpdatedAt = time.Now()

		if err := repo.UpdateAlgorithm(existing); err != nil {
			utils.LogError(err, "Failed to update algorithm")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to update algorithm",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		response := convertAlgorithmToResponse(existing)
		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Algorithm updated successfully",
			Data:    response,
		})
	}
}

// DeleteAlgorithm deletes an algorithm
func DeleteAlgorithm(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Algorithm ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "User ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Invalid user ID",
				Code:    http.StatusBadRequest,
			})
			return
		}

		repo := repos.NewAlgorithmRepository(db.GetConnection())
		if err := repo.DeleteAlgorithm(id, userID); err != nil {
			utils.LogError(err, "Failed to delete algorithm")
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Algorithm not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Algorithm deleted successfully",
			Data:    nil,
		})
	}
}

// convertAlgorithmToResponse converts a data.Algorithm to dto.AlgorithmResponse
func convertAlgorithmToResponse(algo *data.Algorithm) dto.AlgorithmResponse {
	response := dto.AlgorithmResponse{
		ID:            algo.ID,
		UserID:        algo.UserID,
		Name:          algo.Name,
		Description:   algo.Description,
		Code:          algo.Code,
		Status:        algo.Status,
		Symbol:        algo.Symbol,
		Timeframe:     algo.Timeframe,
		ExecutionMode: algo.ExecutionMode,
		Broker:        algo.Broker,
		Enabled:       algo.Enabled,
		Config:        algo.Config,
		State:         algo.State,
		LastSignal:    algo.LastSignal,
		TotalTrades:   algo.TotalTrades,
		WinRate:       algo.WinRate,
		TotalPnL:      algo.TotalPnL,
		Version:       algo.Version,
		Tags:          algo.Tags,
		CreatedAt:     algo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     algo.UpdatedAt.Format(time.RFC3339),
	}

	if algo.LastRunAt != nil {
		lastRunAtStr := algo.LastRunAt.Format(time.RFC3339)
		response.LastRunAt = &lastRunAtStr
	}

	return response
}
