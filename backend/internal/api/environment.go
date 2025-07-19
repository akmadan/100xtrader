package api

import (
	"net/http"

	"github.com/akshitmadan/100xtrader/backend/internal/api/dto"
	"github.com/akshitmadan/100xtrader/backend/internal/data/repos"
	"github.com/akshitmadan/100xtrader/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

// EnvironmentHandler holds dependencies for environment endpoints

type EnvironmentHandler struct {
	Repo repos.EnvironmentRepository
}

func NewEnvironmentHandler(repo repos.EnvironmentRepository) *EnvironmentHandler {
	return &EnvironmentHandler{Repo: repo}
}

// RegisterEnvironmentRoutes registers environment-related routes with DI
func RegisterEnvironmentRoutes(r *gin.Engine, handler *EnvironmentHandler) {
	r.GET("/environments", handler.GetEnvironments)
	r.POST("/environments/:id/start", handler.StartEnvironment)
}

// GetEnvironments godoc
// @Summary List all trading environments
// @Description Returns all available trading environments
// @Tags environments
// @Produce json
// @Success 200 {object} dto.EnvironmentListResponse
// @Failure 500 {object} map[string]string
// @Router /environments [get]
func (h *EnvironmentHandler) GetEnvironments(c *gin.Context) {
	envs, err := h.Repo.ListEnvironments()
	if err != nil {
		utils.Logger.WithError(err).Error("failed to list environments")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list environments"})
		return
	}
	var resp dto.EnvironmentListResponse
	for _, env := range envs {
		resp.Environments = append(resp.Environments, dto.EnvironmentListItem{
			ID:          env.ID,
			Name:        env.Name,
			Description: env.Description,
			Volatility:  env.Volatility,
			Trend:       env.Trend,
			Liquidity:   env.Liquidity,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// StartEnvironment godoc
// @Summary Start a trading environment
// @Description Starts a trading environment by ID
// @Tags environments
// @Param id path string true "Environment ID"
// @Success 200 {object} dto.StartEnvironmentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /environments/{id}/start [post]
func (h *EnvironmentHandler) StartEnvironment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing environment id"})
		return
	}
	env, err := h.Repo.GetEnvironmentByID(id)
	if err != nil {
		utils.Logger.WithError(err).Error("failed to get environment")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get environment"})
		return
	}
	if env == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "environment not found"})
		return
	}
	// TODO: Start environment logic (spawn agents, etc.)
	c.JSON(http.StatusOK, dto.StartEnvironmentResponse{Status: "started"})
}
