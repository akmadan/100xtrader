package brokers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DhanService implements BrokerService for Dhan broker
type DhanService struct{}

// NewDhanService creates a new Dhan broker service
func NewDhanService() *DhanService {
	return &DhanService{}
}

// GetBrokerName returns the broker name
func (d *DhanService) GetBrokerName() data.TradingBroker {
	return data.TradingBrokerDhan
}

// DhanTrade represents a single trade from Dhan API
type DhanTrade struct {
	DhanClientID               string  `json:"dhanClientId"`
	OrderID                    string  `json:"orderId"`
	ExchangeOrderID            string  `json:"exchangeOrderId"`
	ExchangeTradeID            string  `json:"exchangeTradeId"`
	TransactionType            string  `json:"transactionType"`
	ExchangeSegment            string  `json:"exchangeSegment"`
	ProductType                string  `json:"productType"`
	OrderType                  string  `json:"orderType"`
	CustomSymbol               string  `json:"customSymbol"`
	SecurityID                 string  `json:"securityId"`
	TradedQuantity             int     `json:"tradedQuantity"`
	TradedPrice                float64 `json:"tradedPrice"`
	ISIN                       string  `json:"isin"`
	Instrument                 string  `json:"instrument"`
	SebiTax                    float64 `json:"sebiTax"`
	STT                        float64 `json:"stt"`
	BrokerageCharges           float64 `json:"brokerageCharges"`
	ServiceTax                 float64 `json:"serviceTax"`
	ExchangeTransactionCharges float64 `json:"exchangeTransactionCharges"`
	StampDuty                  float64 `json:"stampDuty"`
	CreateTime                 string  `json:"createTime"`
	UpdateTime                 string  `json:"updateTime"`
	ExchangeTime               string  `json:"exchangeTime"`
	DrvExpiryDate              string  `json:"drvExpiryDate"`
	DrvOptionType              string  `json:"drvOptionType"`
	DrvStrikePrice             float64 `json:"drvStrikePrice"`
}

// ParseTrades parses Dhan trade data and converts it to BrokerTrade format
func (d *DhanService) ParseTrades(rawData []byte) ([]BrokerTrade, error) {
	var trades []DhanTrade
	if err := json.Unmarshal(rawData, &trades); err != nil {
		return nil, fmt.Errorf("failed to parse Dhan data: %w", err)
	}

	brokerTrades := make([]BrokerTrade, 0, len(trades))
	for _, trade := range trades {
		// Use customSymbol or construct symbol from other fields
		symbol := trade.CustomSymbol
		if symbol == "" {
			symbol = trade.SecurityID // Fallback to security ID
		}

		brokerTrade := BrokerTrade{
			Symbol:          symbol,
			Quantity:        trade.TradedQuantity,
			Price:           trade.TradedPrice,
			TransactionType: trade.TransactionType,
			ExchangeOrderID: trade.ExchangeOrderID,
			OrderID:         trade.OrderID,
			ProductType:     trade.ProductType,
			ExchangeTime:    trade.ExchangeTime,
		}
		brokerTrades = append(brokerTrades, brokerTrade)
	}

	return brokerTrades, nil
}

// ConvertToTrade converts BrokerTrade to the internal Trade model
func (d *DhanService) ConvertToTrade(brokerTrade BrokerTrade, userID int) (*data.Trade, error) {
	return ConvertBrokerTradeToTrade(brokerTrade, userID, data.TradingBrokerDhan)
}

// FetchTrades fetches trades from Dhan API for the given date range
// fromDate and toDate should be in YYYY-MM-DD format
// accessToken is the Dhan API access token
// pageNumber is the page number for pagination (starts from 0)
func (d *DhanService) FetchTrades(accessToken, fromDate, toDate string, pageNumber int) ([]DhanTrade, error) {
	// Get base URL from environment - check both possible env var names
	baseURL := os.Getenv("DHAN_PROD_API")
	if baseURL == "" {
		baseURL = os.Getenv("DHAN_PROD_API_ENDPOINT")
	}
	if baseURL == "" {
		// Default to production API (api.dhan.co) as per curl example
		baseURL = "https://api.dhan.co/v2"
	}

	// Build the API URL
	url := fmt.Sprintf("%s/trades/%s/%s/%d", baseURL, fromDate, toDate, pageNumber)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("access-token", accessToken)

	// Make HTTP request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	utils.LogInfo("Fetching trades from Dhan API", map[string]interface{}{
		"url":       url,
		"from_date": fromDate,
		"to_date":   toDate,
		"page":      pageNumber,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Dhan API: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("dhan API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log raw response for debugging
	responseBodyStr := string(bodyBytes)
	utils.LogInfo("Dhan API raw response", map[string]interface{}{
		"url":          url,
		"status":       resp.StatusCode,
		"body_len":     len(bodyBytes),
		"body_preview": responseBodyStr[:min(500, len(bodyBytes))], // First 500 chars
		"full_body":    responseBodyStr,                            // Full response for debugging
	})

	// Try to parse as array first
	var trades []DhanTrade
	if err := json.Unmarshal(bodyBytes, &trades); err != nil {
		// If direct array parsing fails, try parsing as wrapped object
		var wrappedResponse struct {
			Data   []DhanTrade `json:"data"`
			Trades []DhanTrade `json:"trades"`
			Result []DhanTrade `json:"result"`
		}
		if err2 := json.Unmarshal(bodyBytes, &wrappedResponse); err2 != nil {
			utils.LogError(err, "Failed to parse Dhan API response as array", map[string]interface{}{
				"response_body": responseBodyStr,
			})
			utils.LogError(err2, "Failed to parse Dhan API response as wrapped object", map[string]interface{}{
				"response_body": responseBodyStr,
			})
			return nil, fmt.Errorf("failed to parse Dhan API response: %w (also tried wrapped format: %v)", err, err2)
		}
		// Use whichever field has data
		if len(wrappedResponse.Data) > 0 {
			trades = wrappedResponse.Data
		} else if len(wrappedResponse.Trades) > 0 {
			trades = wrappedResponse.Trades
		} else if len(wrappedResponse.Result) > 0 {
			trades = wrappedResponse.Result
		}
	}

	utils.LogInfo("Successfully fetched trades from Dhan API", map[string]interface{}{
		"count": len(trades),
		"url":   url,
	})

	return trades, nil
}

// FetchTradesForDateRange fetches all trades for a date range, handling pagination
// fromDate and toDate should be in YYYY-MM-DD format
// accessToken is the Dhan API access token
func (d *DhanService) FetchTradesForDateRange(accessToken, fromDate, toDate string) ([]DhanTrade, error) {
	allTrades := make([]DhanTrade, 0)
	pageNumber := 0 // Start from page 0 as per API documentation (curl example uses 0)

	utils.LogInfo("Starting paginated fetch for date range", map[string]interface{}{
		"from_date":  fromDate,
		"to_date":    toDate,
		"start_page": pageNumber,
	})

	for {
		trades, err := d.FetchTrades(accessToken, fromDate, toDate, pageNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch trades for page %d: %w", pageNumber, err)
		}

		utils.LogInfo("Fetched page", map[string]interface{}{
			"page":  pageNumber,
			"count": len(trades),
		})

		// If no trades returned, we've reached the end
		if len(trades) == 0 {
			utils.LogInfo("No more trades found, stopping pagination", map[string]interface{}{
				"last_page": pageNumber,
				"total":     len(allTrades),
			})
			break
		}

		allTrades = append(allTrades, trades...)

		// If we got fewer trades than expected, we might be on the last page
		// Note: Dhan API doesn't specify page size, so we'll fetch until empty
		// You may want to adjust this logic based on actual API behavior
		if len(trades) < 100 { // Assuming page size is around 100, adjust as needed
			utils.LogInfo("Last page detected (less than 100 trades)", map[string]interface{}{
				"page":  pageNumber,
				"count": len(trades),
				"total": len(allTrades),
			})
			break
		}

		pageNumber++
	}

	utils.LogInfo("Completed paginated fetch", map[string]interface{}{
		"total_trades": len(allTrades),
		"total_pages":  pageNumber + 1,
	})

	return allTrades, nil
}

// FetchAndConvertTrades fetches trades from Dhan API and converts them to internal Trade model
// This is a convenience method that combines fetching and conversion
// fromDate and toDate should be in YYYY-MM-DD format
// accessToken is the Dhan API access token
// userID is the user ID to associate the trades with
func (d *DhanService) FetchAndConvertTrades(accessToken, fromDate, toDate string, userID int) ([]*data.Trade, error) {
	// Fetch trades from API
	dhanTrades, err := d.FetchTradesForDateRange(accessToken, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trades: %w", err)
	}

	// Convert to BrokerTrade format with full Dhan trade data
	enhancedTrades := make([]EnhancedBrokerTrade, 0, len(dhanTrades))
	for _, trade := range dhanTrades {
		symbol := trade.CustomSymbol
		if symbol == "" {
			symbol = trade.SecurityID // Fallback to security ID
		}

		// Parse exchange time
		exchangeTime, err := time.Parse("2006-01-02 15:04:05", trade.ExchangeTime)
		if err != nil {
			exchangeTime, err = time.Parse(time.RFC3339, trade.ExchangeTime)
			if err != nil {
				exchangeTime = time.Now()
			}
		}

		enhancedTrade := EnhancedBrokerTrade{
			BrokerTrade: BrokerTrade{
				Symbol:          symbol,
				Quantity:        trade.TradedQuantity,
				Price:           trade.TradedPrice,
				TransactionType: trade.TransactionType,
				ExchangeOrderID: trade.ExchangeOrderID,
				OrderID:         trade.OrderID,
				ProductType:     trade.ProductType,
				ExchangeTime:    trade.ExchangeTime,
			},
			ExchangeTime: exchangeTime,
			ProductType:  trade.ProductType,
		}
		enhancedTrades = append(enhancedTrades, enhancedTrade)
	}

	// Match BUY and SELL trades to calculate P&L
	matchedTrades := d.matchBuySellTrades(enhancedTrades, userID)

	return matchedTrades, nil
}

// EnhancedBrokerTrade represents a broker trade with parsed time
type EnhancedBrokerTrade struct {
	BrokerTrade
	ExchangeTime time.Time
	ProductType  string
}

// matchBuySellTrades matches BUY and SELL trades to create complete trade entries with P&L
// Creates a separate trade entry for each transaction, then tries to find matching exit prices
func (d *DhanService) matchBuySellTrades(enhancedTrades []EnhancedBrokerTrade, userID int) []*data.Trade {
	// First, create a trade entry for each transaction
	trades := make([]*data.Trade, 0, len(enhancedTrades))
	usedSells := make(map[int]bool) // Track which SELL trades have been used for matching

	// Sort trades by time to process in chronological order
	sortedTrades := make([]EnhancedBrokerTrade, len(enhancedTrades))
	copy(sortedTrades, enhancedTrades)
	for i := 0; i < len(sortedTrades)-1; i++ {
		for j := i + 1; j < len(sortedTrades); j++ {
			if sortedTrades[i].ExchangeTime.After(sortedTrades[j].ExchangeTime) {
				sortedTrades[i], sortedTrades[j] = sortedTrades[j], sortedTrades[i]
			}
		}
	}

	// Create trade entries and try to match exit prices
	for i, trade := range sortedTrades {
		// Determine direction
		direction := data.TradeDirectionLong
		if trade.TransactionType == "SELL" {
			direction = data.TradeDirectionShort
		}

		// Convert product type
		var productType *data.ProductType
		if trade.ProductType != "" {
			pt := convertProductType(trade.ProductType)
			productType = &pt
		}

		// Convert transaction type
		transactionType := "buy"
		if trade.TransactionType == "SELL" {
			transactionType = "sell"
		}

		// Try to find matching exit price
		var exitPrice *float64
		var pnl float64
		var outcome data.OutcomeSummary

		if direction == data.TradeDirectionLong && trade.TransactionType == "BUY" {
			// For LONG: find a SELL on same or later date with matching quantity
			for j := i + 1; j < len(sortedTrades); j++ {
				if usedSells[j] {
					continue
				}
				sell := sortedTrades[j]
				if sell.Symbol == trade.Symbol &&
					sell.TransactionType == "SELL" &&
					sell.ProductType == trade.ProductType &&
					sell.Quantity == trade.Quantity &&
					!sell.ExchangeTime.Before(trade.ExchangeTime) {
					exitPrice = &sell.Price
					pnl = (sell.Price - trade.Price) * float64(trade.Quantity)
					usedSells[j] = true
					break
				}
			}
		} else if direction == data.TradeDirectionShort && trade.TransactionType == "SELL" {
			// For SHORT: find a BUY on same or later date with matching quantity
			for j := i + 1; j < len(sortedTrades); j++ {
				if usedSells[j] {
					continue
				}
				buy := sortedTrades[j]
				if buy.Symbol == trade.Symbol &&
					buy.TransactionType == "BUY" &&
					buy.ProductType == trade.ProductType &&
					buy.Quantity == trade.Quantity &&
					!buy.ExchangeTime.Before(trade.ExchangeTime) {
					exitPrice = &buy.Price
					pnl = (trade.Price - buy.Price) * float64(trade.Quantity)
					usedSells[j] = true
					break
				}
			}
		}

		// Determine outcome based on P&L
		if exitPrice != nil {
			if pnl > 0 {
				outcome = data.OutcomeSummaryProfitable
			} else if pnl < 0 {
				outcome = data.OutcomeSummaryLoss
			} else {
				outcome = data.OutcomeSummaryBreakeven
			}
		} else {
			outcome = data.OutcomeSummaryBreakeven
		}

		tradeEntry := &data.Trade{
			ID:             utils.GenerateID(),
			UserID:         userID,
			Symbol:         trade.Symbol,
			MarketType:     data.MarketTypeIndian,
			EntryDate:      trade.ExchangeTime,
			EntryPrice:     trade.Price,
			Quantity:       trade.Quantity,
			TotalAmount:    trade.Price * float64(trade.Quantity),
			ExitPrice:      exitPrice,
			Direction:      direction,
			StopLoss:       nil,
			Target:         nil,
			Strategy:       "",
			OutcomeSummary: outcome,
			TradeAnalysis:  nil,
			RulesFollowed:  []string{},
			Screenshots:    []string{},
			Psychology:     nil,
			// Broker-specific fields
			TradingBroker:   func() *data.TradingBroker { b := data.TradingBrokerDhan; return &b }(),
			TraderBrokerID:  nil,
			ExchangeOrderID: &trade.ExchangeOrderID,
			OrderID:         &trade.OrderID,
			ProductType:     productType,
			TransactionType: &transactionType,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		trades = append(trades, tradeEntry)
	}

	return trades
}

// FormatDateForAPI formats a time.Time to YYYY-MM-DD format for Dhan API
func FormatDateForAPI(date time.Time) string {
	return date.Format("2006-01-02")
}

// RenewTokenResponse represents the response from Dhan RenewToken API
type RenewTokenResponse struct {
	Status      string `json:"status"`
	AccessToken string `json:"accessToken"`
	ExpiryTime  string `json:"expiryTime"`
}

// RenewToken renews the Dhan access token
// accessToken is the current access token
// dhanClientID is the Dhan client ID
func (d *DhanService) RenewToken(accessToken, dhanClientID string) (*RenewTokenResponse, error) {
	// Get base URL from environment - check both possible env var names, use same base as FetchTrades for consistency
	baseURL := os.Getenv("DHAN_PROD_API")
	if baseURL == "" {
		baseURL = os.Getenv("DHAN_PROD_API_ENDPOINT")
	}
	if baseURL == "" {
		baseURL = "https://api.dhan.co/v2" // Default to api (same as FetchTrades)
	}

	// Build the API URL
	url := fmt.Sprintf("%s/RenewToken", baseURL)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("access-token", accessToken)
	req.Header.Set("dhanClientId", dhanClientID)

	// Make HTTP request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	utils.LogInfo("Renewing Dhan access token", map[string]interface{}{
		"url":            url,
		"dhan_client_id": dhanClientID,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Dhan API: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		errorMsg := string(bodyBytes)
		if errorMsg == "" {
			errorMsg = "No error message provided by Dhan API"
		}
		if readErr != nil {
			errorMsg = fmt.Sprintf("Failed to read error response: %v", readErr)
		}
		return nil, fmt.Errorf("dhan API returned status %d: %s", resp.StatusCode, errorMsg)
	}

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var response RenewTokenResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to parse Dhan API response: %w", err)
	}

	utils.LogInfo("Successfully renewed Dhan access token", map[string]interface{}{
		"status": response.Status,
	})

	return &response, nil
}

// GenerateConsentResponse represents the response from Dhan GenerateConsent API
type GenerateConsentResponse struct {
	ConsentAppID     string `json:"consentAppId"`
	ConsentAppStatus string `json:"consentAppStatus"`
	Status           string `json:"status"`
}

// GenerateConsent generates a consent for Dhan OAuth flow
// dhanClientID is the Dhan client ID (required)
// appID is the Dhan API key
// appSecret is the Dhan API secret
func (d *DhanService) GenerateConsent(dhanClientID, appID, appSecret string) (*GenerateConsentResponse, error) {
	// Build the API URL - client_id is required according to API docs
	if dhanClientID == "" {
		return nil, fmt.Errorf("dhan client ID is required")
	}
	url := fmt.Sprintf("https://auth.dhan.co/app/generate-consent?client_id=%s", dhanClientID)

	// Create HTTP request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("app_id", appID)
	req.Header.Set("app_secret", appSecret)

	// Make HTTP request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	utils.LogInfo("Generating Dhan consent", map[string]interface{}{
		"url":            url,
		"dhan_client_id": dhanClientID,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Dhan API: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		errorMsg := string(bodyBytes)
		if errorMsg == "" {
			errorMsg = "No error message provided by Dhan API"
		}
		if readErr != nil {
			errorMsg = fmt.Sprintf("Failed to read error response: %v", readErr)
		}
		return nil, fmt.Errorf("dhan API returned status %d: %s", resp.StatusCode, errorMsg)
	}

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var response GenerateConsentResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to parse Dhan API response: %w", err)
	}

	utils.LogInfo("Successfully generated Dhan consent", map[string]interface{}{
		"consent_app_id": response.ConsentAppID,
		"status":         response.Status,
	})

	return &response, nil
}

// ConsumeConsentResponse represents the response from Dhan ConsumeConsent API
type ConsumeConsentResponse struct {
	DhanClientID         string `json:"dhanClientId"`
	DhanClientName       string `json:"dhanClientName"`
	DhanClientUcc        string `json:"dhanClientUcc"`
	GivenPowerOfAttorney bool   `json:"givenPowerOfAttorney"`
	AccessToken          string `json:"accessToken"`
	ExpiryTime           string `json:"expiryTime"`
}

// ConsumeConsent consumes the consent token and gets the access token
// tokenID is the token ID received from the browser login redirect
// appID is the Dhan API key
// appSecret is the Dhan API secret
func (d *DhanService) ConsumeConsent(tokenID, appID, appSecret string) (*ConsumeConsentResponse, error) {
	// Build the API URL
	url := fmt.Sprintf("https://auth.dhan.co/app/consumeApp-consent?tokenId=%s", tokenID)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("app_id", appID)
	req.Header.Set("app_secret", appSecret)

	// Make HTTP request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	utils.LogInfo("Consuming Dhan consent", map[string]interface{}{
		"url": url,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Dhan API: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		errorMsg := string(bodyBytes)
		if errorMsg == "" {
			errorMsg = "No error message provided by Dhan API"
		}
		if readErr != nil {
			errorMsg = fmt.Sprintf("Failed to read error response: %v", readErr)
		}
		return nil, fmt.Errorf("dhan API returned status %d: %s", resp.StatusCode, errorMsg)
	}

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var response ConsumeConsentResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to parse Dhan API response: %w", err)
	}

	utils.LogInfo("Successfully consumed Dhan consent", map[string]interface{}{
		"dhan_client_id": response.DhanClientID,
		"status":         "success",
	})

	return &response, nil
}
