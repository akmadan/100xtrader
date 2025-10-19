package middleware

import (
	"time"

	"100xtrader/go-core/internal/utils"

	"github.com/gin-gonic/gin"
)

// Logger returns a gin.HandlerFunc for logging
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Use our custom logger instead of gin's default
		utils.LogRequest(
			param.Method,
			param.Path,
			param.Request.UserAgent(),
			param.StatusCode,
			param.Latency,
			map[string]interface{}{
				"client_ip": param.ClientIP,
				"bytes_in":  param.BodySize,
			},
		)
		return ""
	})
}

// Recovery returns a gin.HandlerFunc for recovery
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		utils.LogError(nil, "Panic recovered", map[string]interface{}{
			"error":  recovered,
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})
		c.AbortWithStatus(500)
	})
}

// CORS returns a gin.HandlerFunc for CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestID adds a request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
