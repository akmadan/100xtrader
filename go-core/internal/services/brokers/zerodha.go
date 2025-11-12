package brokers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// ZerodhaService implements BrokerService for Zerodha broker
type ZerodhaService struct{}

// NewZerodhaService creates a new Zerodha broker service
func NewZerodhaService() *ZerodhaService {
	return &ZerodhaService{}
}

// GetBrokerName returns the broker name
func (z *ZerodhaService) GetBrokerName() data.TradingBroker {
	return data.TradingBrokerZerodha
}

// ZerodhaTradeResponse represents the Zerodha API response structure
type ZerodhaTradeResponse struct {
	Status string `json:"status"`
	Data   []struct {
		TradeID           string  `json:"trade_id"`
		OrderID           string  `json:"order_id"`
		Exchange          string  `json:"exchange"`
		Tradingsymbol     string  `json:"tradingsymbol"`
		InstrumentToken   int     `json:"instrument_token"`
		Product           string  `json:"product"`
		AveragePrice      float64 `json:"average_price"`
		Quantity          int     `json:"quantity"`
		ExchangeOrderID   string  `json:"exchange_order_id"`
		TransactionType   string  `json:"transaction_type"`
		FillTimestamp     string  `json:"fill_timestamp"`
		OrderTimestamp    string  `json:"order_timestamp"`
		ExchangeTimestamp string  `json:"exchange_timestamp"`
	} `json:"data"`
}

// ParseTrades parses Zerodha trade data and converts it to BrokerTrade format
func (z *ZerodhaService) ParseTrades(rawData []byte) ([]BrokerTrade, error) {
	var response ZerodhaTradeResponse
	if err := json.Unmarshal(rawData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse zerodha data: %w", err)
	}

	if response.Status != "success" {
		return nil, fmt.Errorf("zerodha API returned non-success status: %s", response.Status)
	}

	brokerTrades := make([]BrokerTrade, 0, len(response.Data))
	for _, trade := range response.Data {
		brokerTrade := BrokerTrade{
			Symbol:          trade.Tradingsymbol,
			Quantity:        trade.Quantity,
			Price:           trade.AveragePrice,
			TransactionType: trade.TransactionType,
			ExchangeOrderID: trade.ExchangeOrderID,
			OrderID:         trade.OrderID,
			ProductType:     trade.Product,
			ExchangeTime:    trade.FillTimestamp,
		}
		brokerTrades = append(brokerTrades, brokerTrade)
	}

	return brokerTrades, nil
}

// ConvertToTrade converts BrokerTrade to the internal Trade model
func (z *ZerodhaService) ConvertToTrade(brokerTrade BrokerTrade, userID int) (*data.Trade, error) {
	return ConvertBrokerTradeToTrade(brokerTrade, userID, data.TradingBrokerZerodha)
}

// FetchTrades fetches today's trades from Zerodha API
// Zerodha API only allows fetching trades for the current day
// apiKey is the Zerodha API key
// accessToken is the Zerodha access token
func (z *ZerodhaService) FetchTrades(apiKey, accessToken string) ([]BrokerTrade, error) {
	// Build the API URL
	url := "https://api.kite.trade/trades"

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Kite-Version", "3")
	req.Header.Set("Authorization", fmt.Sprintf("token %s:%s", apiKey, accessToken))

	// Make HTTP request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	utils.LogInfo("Fetching trades from Zerodha API", map[string]interface{}{
		"url": url,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Zerodha API: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("zerodha API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response using existing ParseTrades method
	brokerTrades, err := z.ParseTrades(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Zerodha API response: %w", err)
	}

	utils.LogInfo("Successfully fetched trades from Zerodha API", map[string]interface{}{
		"count": len(brokerTrades),
	})

	return brokerTrades, nil
}

// FetchAndConvertTrades fetches today's trades from Zerodha API and converts them to internal Trade model
// This is a convenience method that combines fetching and conversion
// apiKey is the Zerodha API key
// accessToken is the Zerodha access token
// userID is the user ID to associate the trades with
func (z *ZerodhaService) FetchAndConvertTrades(apiKey, accessToken string, userID int) ([]*data.Trade, error) {
	// Fetch trades from API
	brokerTrades, err := z.FetchTrades(apiKey, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trades: %w", err)
	}

	// Convert to internal Trade model
	trades := make([]*data.Trade, 0, len(brokerTrades))
	for _, brokerTrade := range brokerTrades {
		trade, err := z.ConvertToTrade(brokerTrade, userID)
		if err != nil {
			utils.LogError(err, "Failed to convert broker trade", map[string]interface{}{
				"symbol":   brokerTrade.Symbol,
				"order_id": brokerTrade.OrderID,
			})
			continue // Skip this trade but continue with others
		}
		trades = append(trades, trade)
	}

	return trades, nil
}
