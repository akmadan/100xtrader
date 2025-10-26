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

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User data"
// @Success 201 {object} dto.SuccessResponse{data=dto.UserResponse} "User created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users [post]
func CreateUser(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind user request")
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
			utils.LogError(err, "Validation failed for user request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Convert DTO to model
		user := &data.User{
			ID:        0, // Will be set by database
			Name:      req.Username,
			Email:     req.Email,
			CreatedAt: utils.GetCurrentTime(),
		}

		// Create user in database
		repo := repos.NewUserRepository(db.GetConnection())
		if err := repo.CreateUser(user); err != nil {
			utils.LogError(err, "Failed to create user")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to create user",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Convert to response DTO
		response := convertUserToResponse(user)

		utils.LogInfo("User created successfully", map[string]interface{}{
			"user_id": user.ID,
		})

		c.JSON(http.StatusCreated, dto.SuccessResponse{
			Message: "User created successfully",
			Data:    response,
		})
	}
}

// GetUser retrieves a user by ID
// @Summary Get a user by ID
// @Description Retrieve a specific user with all its details
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.UserResponse} "User retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid user ID"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [get]
func GetUser(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "User ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		repo := repos.NewUserRepository(db.GetConnection())
		user, err := repo.GetUserByID(userID)
		if err != nil {
			utils.LogError(err, "Failed to get user", map[string]interface{}{
				"user_id": userID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "User not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		response := convertUserToResponse(user)

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User retrieved successfully",
			Data:    response,
		})
	}
}

// UpdateUser updates an existing user
// @Summary Update a user
// @Description Update an existing user with new data
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UpdateUserRequest true "Updated user data"
// @Success 200 {object} dto.SuccessResponse{data=dto.UserResponse} "User updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [put]
func UpdateUser(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "User ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		var req dto.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind user update request")
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
			utils.LogError(err, "Validation failed for user update request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Convert DTO to model
		user := &data.User{
			ID:    0, // Will be set by database
			Name:  req.Username,
			Email: req.Email,
		}

		// Update user in database
		repo := repos.NewUserRepository(db.GetConnection())
		if err := repo.UpdateUser(user); err != nil {
			utils.LogError(err, "Failed to update user", map[string]interface{}{
				"user_id": userID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to update user",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Get updated user
		updatedUser, err := repo.GetUserByID(userID)
		if err != nil {
			utils.LogError(err, "Failed to get updated user", map[string]interface{}{
				"user_id": userID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve updated user",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		response := convertUserToResponse(updatedUser)

		utils.LogInfo("User updated successfully", map[string]interface{}{
			"user_id": userID,
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User updated successfully",
			Data:    response,
		})
	}
}

// DeleteUser deletes a user
// @Summary Delete a user
// @Description Delete a specific user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.SuccessResponse "User deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid user ID"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [delete]
func DeleteUser(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "User ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		repo := repos.NewUserRepository(db.GetConnection())
		if err := repo.DeleteUser(userID); err != nil {
			utils.LogError(err, "Failed to delete user", map[string]interface{}{
				"user_id": userID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "User not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		utils.LogInfo("User deleted successfully", map[string]interface{}{
			"user_id": userID,
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User deleted successfully",
		})
	}
}

// ListUsers retrieves users with pagination
// @Summary List users
// @Description Retrieve a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Number of users to return (default: 10, max: 100)"
// @Param offset query int false "Number of users to skip (default: 0)"
// @Success 200 {object} dto.SuccessResponse{data=dto.GetUsersResponse} "Users retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users [get]
func ListUsers(db *data.DB) gin.HandlerFunc {
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

		repo := repos.NewUserRepository(db.GetConnection())
		users, total, err := repo.GetUsers(limit, offset)
		if err != nil {
			utils.LogError(err, "Failed to list users")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve users",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Convert to response DTOs
		var userResponses []dto.UserResponse
		for _, user := range users {
			userResponses = append(userResponses, convertUserToResponse(user))
		}

		response := dto.GetUsersResponse{
			Users: userResponses,
			Pagination: dto.PaginationResponse{
				Total:  total,
				Limit:  limit,
				Offset: offset,
				Count:  len(userResponses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Users retrieved successfully",
			Data:    response,
		})
	}
}

// SignInUser handles user sign in
// @Summary Sign in user
// @Description Authenticate user and return user data
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body dto.SignInRequest true "User credentials"
// @Success 200 {object} dto.SuccessResponse{data=dto.UserResponse} "User signed in successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid credentials"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/signin [post]
func SignInUser(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind sign in request")
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
			utils.LogError(err, "Validation failed for sign in request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// TODO: Implement actual authentication logic
		// For now, just return a placeholder response
		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Sign in functionality not implemented yet",
			Data: map[string]interface{}{
				"username": req.Username,
			},
		})
	}
}

// convertUserToResponse converts a data.User to dto.UserResponse
func convertUserToResponse(user *data.User) dto.UserResponse {
	return dto.UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Phone:        user.Phone,
		LastSignedIn: user.LastSignedIn,
		CreatedAt:    user.CreatedAt,
	}
}
