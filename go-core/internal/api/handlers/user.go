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

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user account with name, email, and optional phone
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserCreateRequest true "User creation data"
// @Success 201 {object} dto.SuccessResponse{data=dto.UserResponse} "User created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users [post]
func CreateUser(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.UserCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind user creation request")
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
			utils.LogError(err, "User creation validation failed")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation failed",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Create user
		user := &data.User{
			Name:  req.Name,
			Email: req.Email,
			Phone: req.Phone,
		}

		userRepo := repos.NewUserRepository(db.GetConnection())
		if err := userRepo.Create(user); err != nil {
			utils.LogError(err, "Failed to create user", map[string]interface{}{
				"email": req.Email,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Failed to create user",
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		utils.LogInfo("User created successfully", map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
		})

		// Return response
		response := dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
		}

		c.JSON(http.StatusCreated, dto.SuccessResponse{
			Message: "User created successfully",
			Data:    response,
		})
	}
}

// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Description Retrieve a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.UserResponse} "User retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid user ID"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [get]
func GetUser(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid user ID",
				Message: "User ID must be a valid integer",
				Code:    http.StatusBadRequest,
			})
			return
		}

		userRepo := repos.NewUserRepository(db.GetConnection())
		user, err := userRepo.GetByID(id)
		if err != nil {
			utils.LogError(err, "Failed to get user", map[string]interface{}{
				"user_id": id,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "User not found",
				Message: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}

		response := dto.UserResponse{
			ID:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			Phone:        user.Phone,
			LastSignedIn: user.LastSignedIn,
			CreatedAt:    user.CreatedAt,
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User retrieved successfully",
			Data:    response,
		})
	}
}

// ListUsers retrieves all users
// @Summary List all users
// @Description Retrieve a list of all users in the system
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserListResponse "Users retrieved successfully"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users [get]
func ListUsers(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepo := repos.NewUserRepository(db.GetConnection())
		users, err := userRepo.List()
		if err != nil {
			utils.LogError(err, "Failed to list users")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Failed to retrieve users",
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		var responses []dto.UserResponse
		for _, user := range users {
			response := dto.UserResponse{
				ID:           user.ID,
				Name:         user.Name,
				Email:        user.Email,
				Phone:        user.Phone,
				LastSignedIn: user.LastSignedIn,
				CreatedAt:    user.CreatedAt,
			}
			responses = append(responses, response)
		}

		c.JSON(http.StatusOK, dto.UserListResponse{
			Users: responses,
			Total: len(responses),
		})
	}
}

// UpdateUser updates a user
func UpdateUser(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid user ID",
				Message: "User ID must be a valid integer",
				Code:    http.StatusBadRequest,
			})
			return
		}

		var req dto.UserUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind user update request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid request",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		userRepo := repos.NewUserRepository(db.GetConnection())
		user, err := userRepo.GetByID(id)
		if err != nil {
			utils.LogError(err, "Failed to get user for update", map[string]interface{}{
				"user_id": id,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "User not found",
				Message: err.Error(),
				Code:    http.StatusNotFound,
			})
			return
		}

		// Update fields if provided
		if req.Name != nil {
			user.Name = *req.Name
		}
		if req.Email != nil {
			user.Email = *req.Email
		}
		if req.Phone != nil {
			user.Phone = req.Phone
		}

		if err := userRepo.Update(user); err != nil {
			utils.LogError(err, "Failed to update user", map[string]interface{}{
				"user_id": id,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Failed to update user",
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		response := dto.UserResponse{
			ID:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			Phone:        user.Phone,
			LastSignedIn: user.LastSignedIn,
			CreatedAt:    user.CreatedAt,
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User updated successfully",
			Data:    response,
		})
	}
}

// DeleteUser deletes a user
func DeleteUser(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid user ID",
				Message: "User ID must be a valid integer",
				Code:    http.StatusBadRequest,
			})
			return
		}

		userRepo := repos.NewUserRepository(db.GetConnection())
		if err := userRepo.Delete(id); err != nil {
			utils.LogError(err, "Failed to delete user", map[string]interface{}{
				"user_id": id,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Failed to delete user",
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		utils.LogInfo("User deleted successfully", map[string]interface{}{
			"user_id": id,
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User deleted successfully",
		})
	}
}

// SignInUser handles user sign in
func SignInUser(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.UserSignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind sign in request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid request",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		userRepo := repos.NewUserRepository(db.GetConnection())
		user, err := userRepo.GetByEmail(req.Email)
		if err != nil {
			utils.LogError(err, "Failed to sign in user", map[string]interface{}{
				"email": req.Email,
			})
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Invalid credentials",
				Message: "User not found",
				Code:    http.StatusUnauthorized,
			})
			return
		}

		// Update last sign in time
		now := time.Now()
		user.LastSignedIn = &now
		if err := userRepo.Update(user); err != nil {
			utils.LogError(err, "Failed to update last sign in", map[string]interface{}{
				"user_id": user.ID,
			})
		}

		response := dto.UserSignInResponse{
			User: dto.UserResponse{
				ID:           user.ID,
				Name:         user.Name,
				Email:        user.Email,
				Phone:        user.Phone,
				LastSignedIn: user.LastSignedIn,
				CreatedAt:    user.CreatedAt,
			},
			LastSignIn: now,
			IsNewUser:  false,
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Sign in successful",
			Data:    response,
		})
	}
}
