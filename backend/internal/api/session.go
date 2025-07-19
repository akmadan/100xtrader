package api

import (
	"net/http"
	"time"

	"github.com/akshitmadan/100xtrader/backend/internal/api/dto"
	"github.com/akshitmadan/100xtrader/backend/internal/data"
	"github.com/akshitmadan/100xtrader/backend/internal/data/repos"
	"github.com/gin-gonic/gin"
)

// CreateSession godoc
// @Summary Start a new trading session
// @Description Creates a new trading session
// @Tags sessions
// @Accept json
// @Produce json
// @Param request body dto.CreateSessionRequest true "Create session request"
// @Success 200 {object} dto.CreateSessionResponse
// @Failure 400 {object} map[string]string
// @Router /sessions [post]
func CreateSession(repo repos.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}
		session := &data.Session{
			ID:          req.ID,
			User:        req.User,
			Environment: req.Environment,
			Ticker:      req.Ticker,
			StartedAt:   time.Now().UTC(),
		}
		if err := repo.CreateSession(session); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
			return
		}
		c.JSON(http.StatusOK, dto.CreateSessionResponse{Status: "created"})
	}
}

// EndSession godoc
// @Summary End a trading session
// @Description Ends a trading session by ID
// @Tags sessions
// @Accept json
// @Produce json
// @Param request body dto.EndSessionRequest true "End session request"
// @Success 200 {object} dto.EndSessionResponse
// @Failure 400 {object} map[string]string
// @Router /sessions/end [post]
func EndSession(repo repos.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.EndSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}
		if err := repo.EndSession(req.ID, time.Now().UTC()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to end session"})
			return
		}
		c.JSON(http.StatusOK, dto.EndSessionResponse{Status: "ended"})
	}
}

// ListSessions godoc
// @Summary List all sessions for a user
// @Description Returns all trading sessions for a user
// @Tags sessions
// @Param user query string true "User ID"
// @Produce json
// @Success 200 {object} dto.SessionListResponse
// @Failure 400 {object} map[string]string
// @Router /sessions [get]
func ListSessions(repo repos.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Query("user")
		if user == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user parameter"})
			return
		}
		sessions, err := repo.ListSessions(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sessions"})
			return
		}
		var resp dto.SessionListResponse
		for _, s := range sessions {
			item := dto.SessionItemResponse{
				ID:          s.ID,
				User:        s.User,
				Environment: s.Environment,
				Ticker:      s.Ticker,
				StartedAt:   s.StartedAt,
			}
			if !s.EndedAt.IsZero() {
				item.EndedAt = &s.EndedAt
			}
			resp.Sessions = append(resp.Sessions, item)
		}
		c.JSON(http.StatusOK, resp)
	}
}
