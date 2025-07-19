package api

import (
	"net/http"

	"github.com/akshitmadan/100xtrader/backend/internal/api/dto"
	"github.com/akshitmadan/100xtrader/backend/internal/data"
	"github.com/akshitmadan/100xtrader/backend/internal/data/repos"
	"github.com/akshitmadan/100xtrader/backend/internal/engine"
	"github.com/akshitmadan/100xtrader/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

// EnvironmentHandler holds dependencies for environment endpoints

type EnvironmentHandler struct {
	Repo repos.EnvironmentRepository
	Mgr  *engine.EnvironmentManager
}

func NewEnvironmentHandler(repo repos.EnvironmentRepository, mgr *engine.EnvironmentManager) *EnvironmentHandler {
	return &EnvironmentHandler{Repo: repo, Mgr: mgr}
}

// RegisterEnvironmentRoutes registers all environment-related routes
func RegisterEnvironmentRoutes(r *gin.Engine, handler *EnvironmentHandler) {
	r.GET("/environments", handler.GetEnvironments)
	r.POST("/environments", handler.AddEnvironment)
	r.POST("/environments/current", handler.SetCurrentEnvironment)
	r.GET("/environments/current", handler.GetCurrentEnvironment)
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

// AddEnvironment godoc
// @Summary Add a new environment
// @Description Adds a new trading environment
// @Tags environments
// @Accept json
// @Produce json
// @Param request body dto.AddEnvironmentRequest true "Add environment request"
// @Success 200 {object} dto.AddEnvironmentResponse
// @Failure 400 {object} map[string]string
// @Router /environments [post]
func (h *EnvironmentHandler) AddEnvironment(c *gin.Context) {
	var req dto.AddEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	env := &data.Environment{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Volatility:  req.Volatility,
		Trend:       req.Trend,
		Liquidity:   req.Liquidity,
	}
	if err := h.Mgr.AddEnvironment(env); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add environment"})
		return
	}
	c.JSON(http.StatusOK, dto.AddEnvironmentResponse{Status: "added"})
}

// SetCurrentEnvironment godoc
// @Summary Set the current environment
// @Description Sets the current trading environment by ID
// @Tags environments
// @Accept json
// @Produce json
// @Param request body dto.SetCurrentEnvironmentRequest true "Set current environment request"
// @Success 200 {object} dto.SetCurrentEnvironmentResponse
// @Failure 400 {object} map[string]string
// @Router /environments/current [post]
func (h *EnvironmentHandler) SetCurrentEnvironment(c *gin.Context) {
	var req dto.SetCurrentEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	if err := h.Mgr.SetCurrentEnvironment(req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set current environment"})
		return
	}
	c.JSON(http.StatusOK, dto.SetCurrentEnvironmentResponse{Status: "set"})
}

// GetCurrentEnvironment godoc
// @Summary Get the current environment
// @Description Returns the current trading environment
// @Tags environments
// @Produce json
// @Success 200 {object} dto.CurrentEnvironmentResponse
// @Router /environments/current [get]
func (h *EnvironmentHandler) GetCurrentEnvironment(c *gin.Context) {
	env := h.Mgr.GetCurrentEnvironment()
	if env == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no current environment set"})
		return
	}
	c.JSON(http.StatusOK, dto.CurrentEnvironmentResponse{
		ID:          env.ID,
		Name:        env.Name,
		Description: env.Description,
		Volatility:  env.Volatility,
		Trend:       env.Trend,
		Liquidity:   env.Liquidity,
	})
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
