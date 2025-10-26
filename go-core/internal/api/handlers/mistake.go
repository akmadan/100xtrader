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

// CreateMistake creates a new mistake
// @Summary Create a new mistake
// @Description Create a new trading mistake
// @Tags mistakes
// @Accept json
// @Produce json
// @Param mistake body dto.CreateMistakeRequest true "Mistake data"
// @Success 201 {object} dto.SuccessResponse{data=dto.MistakeResponse} "Mistake created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/mistakes [post]
func CreateMistake(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateMistakeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind mistake request")
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
			utils.LogError(err, "Validation failed for mistake request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Convert DTO to model
		mistake := &data.Mistake{
			ID:        utils.GenerateID(),
			UserID:    req.UserID,
			Name:      req.Name,
			Category:  req.Category,
			CreatedAt: utils.GetCurrentTime(),
			UpdatedAt: utils.GetCurrentTime(),
		}

		// Create mistake in database
		repo := repos.NewMistakeRepository(db.GetConnection())
		if err := repo.CreateMistake(mistake); err != nil {
			utils.LogError(err, "Failed to create mistake")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to create mistake",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Convert to response DTO
		response := convertMistakeToResponse(mistake)

		utils.LogInfo("Mistake created successfully", map[string]interface{}{
			"mistake_id": mistake.ID,
			"user_id":    mistake.UserID,
		})

		c.JSON(http.StatusCreated, dto.SuccessResponse{
			Message: "Mistake created successfully",
			Data:    response,
		})
	}
}

// GetMistake retrieves a mistake by ID
// @Summary Get a mistake by ID
// @Description Retrieve a specific mistake with all its details
// @Tags mistakes
// @Accept json
// @Produce json
// @Param id path string true "Mistake ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.MistakeResponse} "Mistake retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid mistake ID"
// @Failure 404 {object} dto.ErrorResponse "Mistake not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/mistakes/{id} [get]
func GetMistake(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		mistakeID := c.Param("id")
		if mistakeID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Mistake ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// For now, we'll use 0 as userID since we don't have authentication yet
		// TODO: Get userID from JWT token or session
		userID := 0

		repo := repos.NewMistakeRepository(db.GetConnection())
		mistake, err := repo.GetMistakeByID(mistakeID, userID)
		if err != nil {
			utils.LogError(err, "Failed to get mistake", map[string]interface{}{
				"mistake_id": mistakeID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Mistake not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		response := convertMistakeToResponse(mistake)

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Mistake retrieved successfully",
			Data:    response,
		})
	}
}

// UpdateMistake updates an existing mistake
// @Summary Update a mistake
// @Description Update an existing mistake with new data
// @Tags mistakes
// @Accept json
// @Produce json
// @Param id path string true "Mistake ID"
// @Param mistake body dto.UpdateMistakeRequest true "Updated mistake data"
// @Success 200 {object} dto.SuccessResponse{data=dto.MistakeResponse} "Mistake updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 404 {object} dto.ErrorResponse "Mistake not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/mistakes/{id} [put]
func UpdateMistake(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		mistakeID := c.Param("id")
		if mistakeID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Mistake ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		var req dto.UpdateMistakeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind mistake update request")
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
			utils.LogError(err, "Validation failed for mistake update request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Convert DTO to model
		mistake := &data.Mistake{
			ID:        mistakeID,
			UserID:    req.UserID,
			Name:      req.Name,
			Category:  req.Category,
			UpdatedAt: utils.GetCurrentTime(),
		}

		// Update mistake in database
		repo := repos.NewMistakeRepository(db.GetConnection())
		if err := repo.UpdateMistake(mistake); err != nil {
			utils.LogError(err, "Failed to update mistake", map[string]interface{}{
				"mistake_id": mistakeID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to update mistake",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Get updated mistake
		updatedMistake, err := repo.GetMistakeByID(mistakeID, req.UserID)
		if err != nil {
			utils.LogError(err, "Failed to get updated mistake", map[string]interface{}{
				"mistake_id": mistakeID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve updated mistake",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		response := convertMistakeToResponse(updatedMistake)

		utils.LogInfo("Mistake updated successfully", map[string]interface{}{
			"mistake_id": mistakeID,
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Mistake updated successfully",
			Data:    response,
		})
	}
}

// DeleteMistake deletes a mistake
// @Summary Delete a mistake
// @Description Delete a specific mistake by ID
// @Tags mistakes
// @Accept json
// @Produce json
// @Param id path string true "Mistake ID"
// @Success 200 {object} dto.SuccessResponse "Mistake deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid mistake ID"
// @Failure 404 {object} dto.ErrorResponse "Mistake not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/mistakes/{id} [delete]
func DeleteMistake(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		mistakeID := c.Param("id")
		if mistakeID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Mistake ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// For now, we'll use 0 as userID since we don't have authentication yet
		// TODO: Get userID from JWT token or session
		userID := 0

		repo := repos.NewMistakeRepository(db.GetConnection())
		if err := repo.DeleteMistake(mistakeID, userID); err != nil {
			utils.LogError(err, "Failed to delete mistake", map[string]interface{}{
				"mistake_id": mistakeID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Mistake not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		utils.LogInfo("Mistake deleted successfully", map[string]interface{}{
			"mistake_id": mistakeID,
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Mistake deleted successfully",
		})
	}
}

// ListMistakes retrieves mistakes with pagination
// @Summary List mistakes
// @Description Retrieve a paginated list of mistakes
// @Tags mistakes
// @Accept json
// @Produce json
// @Param limit query int false "Number of mistakes to return (default: 10, max: 100)"
// @Param offset query int false "Number of mistakes to skip (default: 0)"
// @Success 200 {object} dto.SuccessResponse{data=dto.GetMistakesResponse} "Mistakes retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/mistakes [get]
func ListMistakes(db *data.DB) gin.HandlerFunc {
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

		repo := repos.NewMistakeRepository(db.GetConnection())
		mistakes, err := repo.GetMistakesByUser(0, limit, offset) // 0 means get all users' mistakes
		if err != nil {
			utils.LogError(err, "Failed to list mistakes")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve mistakes",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// For now, we'll use the length of mistakes as total since we don't have a separate count method
		total := len(mistakes)

		// Convert to response DTOs
		var mistakeResponses []dto.MistakeResponse
		for _, mistake := range mistakes {
			mistakeResponses = append(mistakeResponses, convertMistakeToResponse(mistake))
		}

		response := dto.GetMistakesResponse{
			Mistakes: mistakeResponses,
			Pagination: dto.PaginationResponse{
				Total:  total,
				Limit:  limit,
				Offset: offset,
				Count:  len(mistakeResponses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Mistakes retrieved successfully",
			Data:    response,
		})
	}
}

// GetMistakesByUser retrieves mistakes for a specific user
// @Summary Get mistakes by user
// @Description Retrieve a paginated list of mistakes for a specific user
// @Tags mistakes
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param limit query int false "Number of mistakes to return (default: 10, max: 100)"
// @Param offset query int false "Number of mistakes to skip (default: 0)"
// @Success 200 {object} dto.SuccessResponse{data=dto.GetMistakesResponse} "User mistakes retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid user ID or query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/mistakes/user/{user_id} [get]
func GetMistakesByUser(db *data.DB) gin.HandlerFunc {
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

		repo := repos.NewMistakeRepository(db.GetConnection())
		mistakes, err := repo.GetMistakesByUser(userID, limit, offset)
		if err != nil {
			utils.LogError(err, "Failed to get user mistakes", map[string]interface{}{
				"user_id": userID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve user mistakes",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// For now, we'll use the length of mistakes as total since we don't have a separate count method
		total := len(mistakes)

		// Convert to response DTOs
		var mistakeResponses []dto.MistakeResponse
		for _, mistake := range mistakes {
			mistakeResponses = append(mistakeResponses, convertMistakeToResponse(mistake))
		}

		response := dto.GetMistakesResponse{
			Mistakes: mistakeResponses,
			Pagination: dto.PaginationResponse{
				Total:  total,
				Limit:  limit,
				Offset: offset,
				Count:  len(mistakeResponses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User mistakes retrieved successfully",
			Data:    response,
		})
	}
}

// GetMistakesByCategory retrieves mistakes by category for a specific user
// @Summary Get mistakes by category
// @Description Retrieve a paginated list of mistakes filtered by category for a specific user
// @Tags mistakes
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param category path string true "Mistake category"
// @Param limit query int false "Number of mistakes to return (default: 10, max: 100)"
// @Param offset query int false "Number of mistakes to skip (default: 0)"
// @Success 200 {object} dto.SuccessResponse{data=dto.GetMistakesResponse} "User mistakes by category retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid user ID, category, or query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/mistakes/user/{user_id}/category/{category} [get]
func GetMistakesByCategory(db *data.DB) gin.HandlerFunc {
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

		category := c.Param("category")
		if category == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Category is required",
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

		repo := repos.NewMistakeRepository(db.GetConnection())
		mistakes, err := repo.GetMistakesByCategory(userID, data.MistakeCategory(category), limit, offset)
		if err != nil {
			utils.LogError(err, "Failed to get user mistakes by category", map[string]interface{}{
				"user_id":  userID,
				"category": category,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve user mistakes by category",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// For now, we'll use the length of mistakes as total since we don't have a separate count method
		total := len(mistakes)

		// Convert to response DTOs
		var mistakeResponses []dto.MistakeResponse
		for _, mistake := range mistakes {
			mistakeResponses = append(mistakeResponses, convertMistakeToResponse(mistake))
		}

		response := dto.GetMistakesResponse{
			Mistakes: mistakeResponses,
			Pagination: dto.PaginationResponse{
				Total:  total,
				Limit:  limit,
				Offset: offset,
				Count:  len(mistakeResponses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User mistakes by category retrieved successfully",
			Data:    response,
		})
	}
}

// convertMistakeToResponse converts a data.Mistake to dto.MistakeResponse
func convertMistakeToResponse(mistake *data.Mistake) dto.MistakeResponse {
	return dto.MistakeResponse{
		ID:        mistake.ID,
		UserID:    mistake.UserID,
		Name:      mistake.Name,
		Category:  mistake.Category,
		CreatedAt: mistake.CreatedAt,
		UpdatedAt: mistake.UpdatedAt,
	}
}
