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

// CreateRule creates a new rule
// @Summary Create a new rule
// @Description Create a new trading rule
// @Tags rules
// @Accept json
// @Produce json
// @Param rule body dto.CreateRuleRequest true "Rule data"
// @Success 201 {object} dto.SuccessResponse{data=dto.RuleResponse} "Rule created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/rules [post]
func CreateRule(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateRuleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind rule request")
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
			utils.LogError(err, "Validation failed for rule request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Convert DTO to model
		rule := &data.Rule{
			ID:          utils.GenerateID(),
			UserID:      req.UserID,
			Name:        req.Name,
			Description: req.Description,
			Category:    req.Category,
			CreatedAt:   utils.GetCurrentTime(),
			UpdatedAt:   utils.GetCurrentTime(),
		}

		// Create rule in database
		repo := repos.NewRuleRepository(db.GetConnection())
		if err := repo.CreateRule(rule); err != nil {
			utils.LogError(err, "Failed to create rule")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to create rule",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Convert to response DTO
		response := convertRuleToResponse(rule)

		utils.LogInfo("Rule created successfully", map[string]interface{}{
			"rule_id": rule.ID,
			"user_id": rule.UserID,
		})

		c.JSON(http.StatusCreated, dto.SuccessResponse{
			Message: "Rule created successfully",
			Data:    response,
		})
	}
}

// GetRule retrieves a rule by ID
// @Summary Get a rule by ID
// @Description Retrieve a specific rule with all its details
// @Tags rules
// @Accept json
// @Produce json
// @Param id path string true "Rule ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.RuleResponse} "Rule retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid rule ID"
// @Failure 404 {object} dto.ErrorResponse "Rule not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/rules/{id} [get]
func GetRule(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ruleID := c.Param("id")
		if ruleID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Rule ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// For now, we'll use 0 as userID since we don't have authentication yet
		// TODO: Get userID from JWT token or session
		userID := 0

		repo := repos.NewRuleRepository(db.GetConnection())
		rule, err := repo.GetRuleByID(ruleID, userID)
		if err != nil {
			utils.LogError(err, "Failed to get rule", map[string]interface{}{
				"rule_id": ruleID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Rule not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		response := convertRuleToResponse(rule)

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Rule retrieved successfully",
			Data:    response,
		})
	}
}

// UpdateRule updates an existing rule
// @Summary Update a rule
// @Description Update an existing rule with new data
// @Tags rules
// @Accept json
// @Produce json
// @Param id path string true "Rule ID"
// @Param rule body dto.UpdateRuleRequest true "Updated rule data"
// @Success 200 {object} dto.SuccessResponse{data=dto.RuleResponse} "Rule updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 404 {object} dto.ErrorResponse "Rule not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/rules/{id} [put]
func UpdateRule(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ruleID := c.Param("id")
		if ruleID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Rule ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		var req dto.UpdateRuleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind rule update request")
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
			utils.LogError(err, "Validation failed for rule update request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Convert DTO to model
		rule := &data.Rule{
			ID:          ruleID,
			UserID:      req.UserID,
			Name:        req.Name,
			Description: req.Description,
			Category:    req.Category,
			UpdatedAt:   utils.GetCurrentTime(),
		}

		// Update rule in database
		repo := repos.NewRuleRepository(db.GetConnection())
		if err := repo.UpdateRule(rule); err != nil {
			utils.LogError(err, "Failed to update rule", map[string]interface{}{
				"rule_id": ruleID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to update rule",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Get updated rule
		updatedRule, err := repo.GetRuleByID(ruleID, req.UserID)
		if err != nil {
			utils.LogError(err, "Failed to get updated rule", map[string]interface{}{
				"rule_id": ruleID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve updated rule",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		response := convertRuleToResponse(updatedRule)

		utils.LogInfo("Rule updated successfully", map[string]interface{}{
			"rule_id": ruleID,
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Rule updated successfully",
			Data:    response,
		})
	}
}

// DeleteRule deletes a rule
// @Summary Delete a rule
// @Description Delete a specific rule by ID
// @Tags rules
// @Accept json
// @Produce json
// @Param id path string true "Rule ID"
// @Success 200 {object} dto.SuccessResponse "Rule deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid rule ID"
// @Failure 404 {object} dto.ErrorResponse "Rule not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/rules/{id} [delete]
func DeleteRule(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ruleID := c.Param("id")
		if ruleID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Rule ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// For now, we'll use 0 as userID since we don't have authentication yet
		// TODO: Get userID from JWT token or session
		userID := 0

		repo := repos.NewRuleRepository(db.GetConnection())
		if err := repo.DeleteRule(ruleID, userID); err != nil {
			utils.LogError(err, "Failed to delete rule", map[string]interface{}{
				"rule_id": ruleID,
			})
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Rule not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		utils.LogInfo("Rule deleted successfully", map[string]interface{}{
			"rule_id": ruleID,
		})

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Rule deleted successfully",
		})
	}
}

// ListRules retrieves rules with pagination
// @Summary List rules
// @Description Retrieve a paginated list of rules
// @Tags rules
// @Accept json
// @Produce json
// @Param limit query int false "Number of rules to return (default: 10, max: 100)"
// @Param offset query int false "Number of rules to skip (default: 0)"
// @Success 200 {object} dto.SuccessResponse{data=dto.GetRulesResponse} "Rules retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/rules [get]
func ListRules(db *data.DB) gin.HandlerFunc {
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

		repo := repos.NewRuleRepository(db.GetConnection())
		rules, err := repo.GetRulesByUser(0, limit, offset) // 0 means get all users' rules
		if err != nil {
			utils.LogError(err, "Failed to list rules")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve rules",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// For now, we'll use the length of rules as total since we don't have a separate count method
		total := len(rules)

		// Convert to response DTOs
		var ruleResponses []dto.RuleResponse
		for _, rule := range rules {
			ruleResponses = append(ruleResponses, convertRuleToResponse(rule))
		}

		response := dto.GetRulesResponse{
			Rules: ruleResponses,
			Pagination: dto.PaginationResponse{
				Total:  total,
				Limit:  limit,
				Offset: offset,
				Count:  len(ruleResponses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Rules retrieved successfully",
			Data:    response,
		})
	}
}

// GetRulesByUser retrieves rules for a specific user
// @Summary Get rules by user
// @Description Retrieve a paginated list of rules for a specific user
// @Tags rules
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param limit query int false "Number of rules to return (default: 10, max: 100)"
// @Param offset query int false "Number of rules to skip (default: 0)"
// @Success 200 {object} dto.SuccessResponse{data=dto.GetRulesResponse} "User rules retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid user ID or query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/rules/user/{user_id} [get]
func GetRulesByUser(db *data.DB) gin.HandlerFunc {
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

		repo := repos.NewRuleRepository(db.GetConnection())
		rules, err := repo.GetRulesByUser(userID, limit, offset)
		if err != nil {
			utils.LogError(err, "Failed to get user rules", map[string]interface{}{
				"user_id": userID,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve user rules",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// For now, we'll use the length of rules as total since we don't have a separate count method
		total := len(rules)

		// Convert to response DTOs
		var ruleResponses []dto.RuleResponse
		for _, rule := range rules {
			ruleResponses = append(ruleResponses, convertRuleToResponse(rule))
		}

		response := dto.GetRulesResponse{
			Rules: ruleResponses,
			Pagination: dto.PaginationResponse{
				Total:  total,
				Limit:  limit,
				Offset: offset,
				Count:  len(ruleResponses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User rules retrieved successfully",
			Data:    response,
		})
	}
}

// GetRulesByCategory retrieves rules by category for a specific user
// @Summary Get rules by category
// @Description Retrieve a paginated list of rules filtered by category for a specific user
// @Tags rules
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param category path string true "Rule category"
// @Param limit query int false "Number of rules to return (default: 10, max: 100)"
// @Param offset query int false "Number of rules to skip (default: 0)"
// @Success 200 {object} dto.SuccessResponse{data=dto.GetRulesResponse} "User rules by category retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid user ID, category, or query parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/rules/user/{user_id}/category/{category} [get]
func GetRulesByCategory(db *data.DB) gin.HandlerFunc {
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

		repo := repos.NewRuleRepository(db.GetConnection())
		rules, err := repo.GetRulesByCategory(userID, data.RuleCategory(category), limit, offset)
		if err != nil {
			utils.LogError(err, "Failed to get user rules by category", map[string]interface{}{
				"user_id":  userID,
				"category": category,
			})
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Database Error",
				Message: "Failed to retrieve user rules by category",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// For now, we'll use the length of rules as total since we don't have a separate count method
		total := len(rules)

		// Convert to response DTOs
		var ruleResponses []dto.RuleResponse
		for _, rule := range rules {
			ruleResponses = append(ruleResponses, convertRuleToResponse(rule))
		}

		response := dto.GetRulesResponse{
			Rules: ruleResponses,
			Pagination: dto.PaginationResponse{
				Total:  total,
				Limit:  limit,
				Offset: offset,
				Count:  len(ruleResponses),
			},
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "User rules by category retrieved successfully",
			Data:    response,
		})
	}
}

// convertRuleToResponse converts a data.Rule to dto.RuleResponse
func convertRuleToResponse(rule *data.Rule) dto.RuleResponse {
	return dto.RuleResponse{
		ID:          rule.ID,
		UserID:      rule.UserID,
		Name:        rule.Name,
		Description: rule.Description,
		Category:    rule.Category,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}
}
