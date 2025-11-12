package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-core/internal/api/dto"
	"go-core/internal/data"
	"go-core/internal/data/repos"
	"go-core/internal/services/brokers"
	"go-core/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// RenewDhanToken renews the Dhan access token
// @Summary Renew Dhan access token
// @Description Renews the Dhan access token for a user
// @Tags dhan
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param request body dto.DhanRenewTokenRequest true "Renew token request"
// @Success 200 {object} dto.SuccessResponse{data=dto.DhanRenewTokenResponse} "Token renewed successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/{user_id}/dhan/renew-token [post]
func RenewDhanToken(db *data.DB) gin.HandlerFunc {
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

		var req dto.DhanRenewTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind renew token request")
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
			utils.LogError(err, "Validation failed for renew token request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Get user's broker config to get stored access token and client ID
		repo := repos.NewUserRepository(db.GetConnection())
		user, err := repo.GetUserByID(strconv.Itoa(userID))
		if err != nil {
			utils.LogError(err, "Failed to get user")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to get user",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Initialize map if nil
		if user.ConfiguredBrokers == nil {
			user.ConfiguredBrokers = make(map[string]data.BrokerConfig)
		}

		config := user.ConfiguredBrokers["dhan"]

		// Use stored access token and client ID if available, otherwise use request
		accessToken := req.AccessToken
		clientID := req.DhanClientID
		if config.AccessToken != "" {
			accessToken = config.AccessToken
		}
		if config.DhanClientID != nil {
			clientID = *config.DhanClientID
		}

		if accessToken == "" || clientID == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "Access token and client ID are required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Call Dhan service to renew token
		dhanService := brokers.NewDhanService()
		response, err := dhanService.RenewToken(accessToken, clientID)
		if err != nil {
			utils.LogError(err, "Failed to renew Dhan token")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to renew token: " + err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Parse expiry time
		expiryTime, err := time.Parse("2006-01-02T15:04:05", response.ExpiryTime)
		if err != nil {
			utils.LogError(err, "Failed to parse expiry time")
			expiryTime = time.Now().Add(24 * time.Hour) // Default to 24 hours
		}

		// Update Dhan config (preserve API key/secret if they exist)
		config.AccessToken = response.AccessToken
		config.ExpiryTime = &expiryTime
		config.ConfiguredAt = time.Now()
		if config.DhanClientID == nil {
			config.DhanClientID = &clientID
		}
		user.ConfiguredBrokers["dhan"] = config

		// Save updated user
		if err := repo.UpdateUser(user); err != nil {
			utils.LogError(err, "Failed to update user broker config")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to update broker configuration",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Token renewed successfully",
			Data: dto.DhanRenewTokenResponse{
				Status:      response.Status,
				AccessToken: response.AccessToken,
				ExpiryTime:  response.ExpiryTime,
			},
		})
	}
}

// GenerateDhanConsent generates a consent for Dhan OAuth flow
// @Summary Generate Dhan consent
// @Description Generates a consent for Dhan OAuth flow (no client_id needed, uses stored API credentials)
// @Tags dhan
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.DhanGenerateConsentResponse} "Consent generated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id}/dhan/generate-consent [post]
func GenerateDhanConsent(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")
		_, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Invalid user ID",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Get user's API credentials from broker config
		repo := repos.NewUserRepository(db.GetConnection())
		user, err := repo.GetUserByID(userIDStr)
		if err != nil {
			utils.LogError(err, "Failed to get user")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to get user",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		config, exists := user.ConfiguredBrokers["dhan"]
		if !exists || config.APIKey == nil || config.APISecret == nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "API key and secret not configured. Please save credentials first.",
				Code:    http.StatusBadRequest,
			})
			return
		}

		appID := *config.APIKey
		appSecret := *config.APISecret

		// Get client_id from config - it's required for generate-consent
		clientID := ""
		if config.DhanClientID != nil && *config.DhanClientID != "" {
			clientID = *config.DhanClientID
		} else {
			// client_id is required for generate-consent API
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "Dhan Client ID is required. Please provide it when saving credentials, or it will be set after first authentication.",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Call Dhan service to generate consent
		dhanService := brokers.NewDhanService()
		response, err := dhanService.GenerateConsent(clientID, appID, appSecret)
		if err != nil {
			utils.LogError(err, "Failed to generate Dhan consent", map[string]interface{}{
				"error":          err.Error(),
				"dhan_client_id": clientID,
			})
			// Check error type and return appropriate status code
			errMsg := err.Error()
			if strings.Contains(errMsg, "dhan client ID is required") {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "Bad Request",
					Message: "Dhan Client ID is required",
					Code:    http.StatusBadRequest,
				})
			} else if strings.Contains(errMsg, "status 401") {
				c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error:   "Unauthorized",
					Message: "Invalid API credentials or Client ID. Please verify your API Key, API Secret, and Dhan Client ID are correct and match your Dhan account.",
					Code:    http.StatusUnauthorized,
				})
			} else if strings.Contains(errMsg, "status 400") {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "Bad Request",
					Message: "Invalid request to Dhan API: " + errMsg,
					Code:    http.StatusBadRequest,
				})
			} else {
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error:   "Internal Server Error",
					Message: "Failed to generate consent: " + errMsg,
					Code:    http.StatusInternalServerError,
				})
			}
			return
		}

		// Build login URL
		loginURL := "https://auth.dhan.co/login/consentApp-login?consentAppId=" + response.ConsentAppID

		// Build callback URL for Dhan redirect configuration
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		host := c.Request.Host
		if host == "" {
			host = "localhost:8080" // Default for local development
		}
		callbackURL := fmt.Sprintf("%s://%s/api/v1/dhan/consent-callback?userId=%s", scheme, host, userIDStr)

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Consent generated successfully",
			Data: dto.DhanGenerateConsentResponse{
				ConsentAppID:     response.ConsentAppID,
				ConsentAppStatus: response.ConsentAppStatus,
				Status:           response.Status,
				LoginURL:         loginURL,
				CallbackURL:      callbackURL,
			},
		})
	}
}

// ConsumeDhanConsent consumes the consent token and saves the access token
// @Summary Consume Dhan consent
// @Description Consumes the consent token and saves the access token for the user
// @Tags dhan
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param request body dto.DhanConsumeConsentRequest true "Consume consent request"
// @Success 200 {object} dto.SuccessResponse{data=dto.DhanConsumeConsentResponse} "Consent consumed successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/{user_id}/dhan/consume-consent [post]
func ConsumeDhanConsent(db *data.DB) gin.HandlerFunc {
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

		var req dto.DhanConsumeConsentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind consume consent request")
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
			utils.LogError(err, "Validation failed for consume consent request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Get user's API credentials from broker config
		repo := repos.NewUserRepository(db.GetConnection())
		user, err := repo.GetUserByID(c.Param("id"))
		if err != nil {
			utils.LogError(err, "Failed to get user")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to get user",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		config, exists := user.ConfiguredBrokers["dhan"]
		if !exists || config.APIKey == nil || config.APISecret == nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "API key and secret not configured. Please configure Dhan first.",
				Code:    http.StatusBadRequest,
			})
			return
		}

		appID := *config.APIKey
		appSecret := *config.APISecret

		// Call Dhan service to consume consent
		dhanService := brokers.NewDhanService()
		response, err := dhanService.ConsumeConsent(req.TokenID, appID, appSecret)
		if err != nil {
			utils.LogError(err, "Failed to consume Dhan consent")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to consume consent",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Parse expiry time
		expiryTime, err := time.Parse("2006-01-02T15:04:05", response.ExpiryTime)
		if err != nil {
			utils.LogError(err, "Failed to parse expiry time")
			expiryTime = time.Now().Add(24 * time.Hour) // Default to 24 hours
		}

		// Update user's broker config (preserve existing API key and secret)
		config.AccessToken = response.AccessToken
		config.DhanClientID = &response.DhanClientID
		config.DhanClientName = &response.DhanClientName
		config.DhanClientUcc = &response.DhanClientUcc
		config.ExpiryTime = &expiryTime
		config.ConfiguredAt = time.Now()
		user.ConfiguredBrokers["dhan"] = config

		if err := repo.UpdateUserBrokerConfig(userID, "dhan", config); err != nil {
			utils.LogError(err, "Failed to update user broker config")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to update broker configuration",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Consent consumed successfully",
			Data: dto.DhanConsumeConsentResponse{
				DhanClientID:         response.DhanClientID,
				DhanClientName:       response.DhanClientName,
				DhanClientUcc:        response.DhanClientUcc,
				GivenPowerOfAttorney: response.GivenPowerOfAttorney,
				AccessToken:          response.AccessToken,
				ExpiryTime:           response.ExpiryTime,
			},
		})
	}
}

// ConsumeDhanConsentCallback is a webhook endpoint that receives tokenId from Dhan OAuth redirect
// This endpoint accepts GET requests with tokenId and userId as query parameters
// It automatically consumes the consent and returns an HTML page with success message
// @Summary Dhan OAuth callback webhook
// @Description Webhook endpoint that receives tokenId from Dhan OAuth redirect and automatically consumes consent
// @Tags dhan
// @Accept html
// @Produce html
// @Param tokenId query string true "Token ID from Dhan OAuth redirect"
// @Param userId query int true "User ID"
// @Success 200 {string} string "HTML success page"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/dhan/consent-callback [get]
func ConsumeDhanConsentCallback(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Helper function to render error HTML
		renderError := func(statusCode int, message string) {
			errorHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dhan Authentication Error</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
            background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
        }
        .container {
            background: white;
            padding: 2rem;
            border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            text-align: center;
            max-width: 500px;
        }
        .error-icon {
            width: 64px;
            height: 64px;
            background: #ef4444;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 0 auto 1rem;
        }
        .error-icon::after {
            content: '✕';
            color: white;
            font-size: 32px;
            font-weight: bold;
        }
        h1 {
            color: #1f2937;
            margin: 0 0 0.5rem 0;
        }
        p {
            color: #6b7280;
            margin: 0.5rem 0;
        }
        .close-btn {
            background: #ef4444;
            color: white;
            border: none;
            padding: 0.75rem 2rem;
            border-radius: 6px;
            font-size: 1rem;
            cursor: pointer;
            margin-top: 1rem;
        }
        .close-btn:hover {
            background: #dc2626;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="error-icon"></div>
        <h1>Authentication Failed</h1>
        <p>` + message + `</p>
        <button class="close-btn" onclick="window.close()">Close Window</button>
    </div>
</body>
</html>`
			c.Data(statusCode, "text/html; charset=utf-8", []byte(errorHTML))
		}

		// Get tokenId from query parameters
		tokenID := c.Query("tokenId")
		if tokenID == "" {
			renderError(http.StatusBadRequest, "Missing tokenId parameter")
			return
		}

		// Get userId from query parameters
		userIDStr := c.Query("userId")
		if userIDStr == "" {
			renderError(http.StatusBadRequest, "Missing userId parameter")
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			renderError(http.StatusBadRequest, "Invalid userId parameter")
			return
		}

		// Get user's API credentials from broker config
		repo := repos.NewUserRepository(db.GetConnection())
		user, err := repo.GetUserByID(userIDStr)
		if err != nil {
			utils.LogError(err, "Failed to get user in consent callback")
			renderError(http.StatusInternalServerError, "Failed to get user information")
			return
		}

		config, exists := user.ConfiguredBrokers["dhan"]
		if !exists || config.APIKey == nil || config.APISecret == nil {
			renderError(http.StatusBadRequest, "API key and secret not configured. Please configure Dhan first.")
			return
		}

		appID := *config.APIKey
		appSecret := *config.APISecret

		// Call Dhan service to consume consent
		dhanService := brokers.NewDhanService()
		response, err := dhanService.ConsumeConsent(tokenID, appID, appSecret)
		if err != nil {
			utils.LogError(err, "Failed to consume Dhan consent in callback")
			renderError(http.StatusInternalServerError, "Failed to consume consent: "+err.Error())
			return
		}

		// Parse expiry time
		expiryTime, err := time.Parse("2006-01-02T15:04:05", response.ExpiryTime)
		if err != nil {
			utils.LogError(err, "Failed to parse expiry time in callback")
			expiryTime = time.Now().Add(24 * time.Hour) // Default to 24 hours
		}

		// Update user's broker config (preserve existing API key and secret)
		config.AccessToken = response.AccessToken
		config.DhanClientID = &response.DhanClientID
		config.DhanClientName = &response.DhanClientName
		config.DhanClientUcc = &response.DhanClientUcc
		config.ExpiryTime = &expiryTime
		config.ConfiguredAt = time.Now()
		user.ConfiguredBrokers["dhan"] = config

		if err := repo.UpdateUserBrokerConfig(userID, "dhan", config); err != nil {
			utils.LogError(err, "Failed to update user broker config in callback")
			renderError(http.StatusInternalServerError, "Failed to save configuration")
			return
		}

		// Return success HTML page
		successHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dhan Authentication Success</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        .container {
            background: white;
            padding: 2rem;
            border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            text-align: center;
            max-width: 500px;
        }
        .success-icon {
            width: 64px;
            height: 64px;
            background: #10b981;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 0 auto 1rem;
        }
        .success-icon::after {
            content: '✓';
            color: white;
            font-size: 32px;
            font-weight: bold;
        }
        h1 {
            color: #1f2937;
            margin: 0 0 0.5rem 0;
        }
        p {
            color: #6b7280;
            margin: 0.5rem 0;
        }
        .info {
            background: #f3f4f6;
            padding: 1rem;
            border-radius: 8px;
            margin: 1rem 0;
            text-align: left;
        }
        .info-item {
            margin: 0.5rem 0;
        }
        .info-label {
            font-weight: 600;
            color: #374151;
        }
        .close-btn {
            background: #667eea;
            color: white;
            border: none;
            padding: 0.75rem 2rem;
            border-radius: 6px;
            font-size: 1rem;
            cursor: pointer;
            margin-top: 1rem;
        }
        .close-btn:hover {
            background: #5568d3;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="success-icon"></div>
        <h1>Authentication Successful!</h1>
        <p>Your Dhan account has been successfully connected.</p>
        <div class="info">
            <div class="info-item">
                <span class="info-label">Client ID:</span> ` + response.DhanClientID + `
            </div>
            <div class="info-item">
                <span class="info-label">Client Name:</span> ` + response.DhanClientName + `
            </div>
            <div class="info-item">
                <span class="info-label">UCC:</span> ` + response.DhanClientUcc + `
            </div>
            <div class="info-item">
                <span class="info-label">Token Expires:</span> ` + response.ExpiryTime + `
            </div>
        </div>
        <p style="font-size: 0.9rem; color: #9ca3af;">You can close this window and return to the application.</p>
        <button class="close-btn" onclick="window.close()">Close Window</button>
    </div>
</body>
</html>`

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(successHTML))
	}
}

// SaveDhanCredentials saves the API key and secret for OAuth flow
// @Summary Save Dhan API credentials
// @Description Saves the Dhan API key and secret for a user (OAuth flow)
// @Tags dhan
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.DhanSaveCredentialsRequest true "Save credentials request"
// @Success 200 {object} dto.SuccessResponse{data=dto.DhanBrokerConfigResponse} "Credentials saved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id}/dhan/save-credentials [post]
func SaveDhanCredentials(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")
		_, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Invalid user ID",
				Code:    http.StatusBadRequest,
			})
			return
		}

		var req dto.DhanSaveCredentialsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind save credentials request")
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
			utils.LogError(err, "Validation failed for save credentials request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Update user's broker config
		repo := repos.NewUserRepository(db.GetConnection())
		user, err := repo.GetUserByID(userIDStr)
		if err != nil {
			utils.LogError(err, "Failed to get user")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to get user",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Initialize map if nil
		if user.ConfiguredBrokers == nil {
			user.ConfiguredBrokers = make(map[string]data.BrokerConfig)
		}

		// Update or create Dhan config
		config := user.ConfiguredBrokers["dhan"]
		config.APIKey = &req.APIKey
		config.APISecret = &req.APISecret
		config.DhanClientID = &req.DhanClientID
		config.ConfiguredAt = time.Now()
		user.ConfiguredBrokers["dhan"] = config

		// Save updated user
		if err := repo.UpdateUser(user); err != nil {
			utils.LogError(err, "Failed to update user broker config")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to save broker configuration",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Credentials saved successfully",
			Data: dto.DhanBrokerConfigResponse{
				Configured: true,
			},
		})
	}
}

// SaveDhanToken saves the access token directly (for direct token method - deprecated, use OAuth flow)
// @Summary Save Dhan access token
// @Description Saves the Dhan access token directly for a user (direct token method)
// @Tags dhan
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.DhanSaveTokenRequest true "Save token request"
// @Success 200 {object} dto.SuccessResponse{data=dto.DhanBrokerConfigResponse} "Token saved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id}/dhan/save-token [post]
func SaveDhanToken(db *data.DB) gin.HandlerFunc {
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

		var req dto.DhanSaveTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.LogError(err, "Failed to bind save token request")
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
			utils.LogError(err, "Validation failed for save token request")
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validation Error",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Update user's broker config
		repo := repos.NewUserRepository(db.GetConnection())
		config := data.BrokerConfig{
			AccessToken:  req.AccessToken,
			DhanClientID: &req.DhanClientID,
			// Set expiry to 24 hours from now (default, can be updated when renewing)
			ExpiryTime:   func() *time.Time { t := time.Now().Add(24 * time.Hour); return &t }(),
			ConfiguredAt: time.Now(),
		}

		if err := repo.UpdateUserBrokerConfig(userID, "dhan", config); err != nil {
			utils.LogError(err, "Failed to update user broker config")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to save broker configuration",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		expiryTimeStr := ""
		if config.ExpiryTime != nil {
			expiryTimeStr = config.ExpiryTime.Format("2006-01-02T15:04:05")
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Token saved successfully",
			Data: dto.DhanBrokerConfigResponse{
				Configured:   true,
				DhanClientID: req.DhanClientID,
				ExpiryTime:   expiryTimeStr,
			},
		})
	}
}

// GetDhanBrokerConfig gets the Dhan broker configuration for a user
// @Summary Get Dhan broker configuration
// @Description Gets the Dhan broker configuration for a user
// @Tags dhan
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.DhanBrokerConfigResponse} "Broker configuration retrieved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/users/{user_id}/dhan/config [get]
func GetDhanBrokerConfig(db *data.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")
		_, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid Request",
				Message: "Invalid user ID",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Get user
		repo := repos.NewUserRepository(db.GetConnection())
		user, err := repo.GetUserByID(userIDStr)
		if err != nil {
			utils.LogError(err, "Failed to get user")
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "Failed to get user",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		// Get Dhan config
		config, exists := user.ConfiguredBrokers["dhan"]
		if !exists || config.AccessToken == "" {
			c.JSON(http.StatusOK, dto.SuccessResponse{
				Message: "Broker configuration retrieved successfully",
				Data: dto.DhanBrokerConfigResponse{
					Configured: false,
				},
			})
			return
		}

		expiryTimeStr := ""
		if config.ExpiryTime != nil {
			expiryTimeStr = config.ExpiryTime.Format("2006-01-02T15:04:05")
		}

		c.JSON(http.StatusOK, dto.SuccessResponse{
			Message: "Broker configuration retrieved successfully",
			Data: dto.DhanBrokerConfigResponse{
				Configured:     true,
				DhanClientID:   getStringValue(config.DhanClientID),
				DhanClientName: getStringValue(config.DhanClientName),
				ExpiryTime:     expiryTimeStr,
			},
		})
	}
}

// Helper function to get string value from pointer
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
