package brokers

import (
	"time"

	"go-core/internal/data"
	"go-core/internal/utils"
)

// ConvertBrokerTradeToTrade converts a BrokerTrade to the internal Trade model
// This is a helper function that can be used by all broker implementations
func ConvertBrokerTradeToTrade(
	brokerTrade BrokerTrade,
	userID int,
	broker data.TradingBroker,
) (*data.Trade, error) {
	// Parse exchange time
	exchangeTime, err := time.Parse("2006-01-02 15:04:05", brokerTrade.ExchangeTime)
	if err != nil {
		// Try alternative formats
		exchangeTime, err = time.Parse(time.RFC3339, brokerTrade.ExchangeTime)
		if err != nil {
			exchangeTime = time.Now() // Fallback to current time
		}
	}

	// Determine direction from transaction type
	var direction data.TradeDirection
	if brokerTrade.TransactionType == "buy" || brokerTrade.TransactionType == "BUY" {
		direction = data.TradeDirectionLong
	} else {
		direction = data.TradeDirectionShort
	}

	// Convert product type
	var productType *data.ProductType
	if brokerTrade.ProductType != "" {
		pt := convertProductType(brokerTrade.ProductType)
		productType = &pt
	}

	// Convert transaction type to lowercase
	transactionType := brokerTrade.TransactionType
	if transactionType != "" {
		if transactionType == "BUY" {
			transactionType = "buy"
		} else if transactionType == "SELL" {
			transactionType = "sell"
		}
	}

	trade := &data.Trade{
		ID:             utils.GenerateID(),
		UserID:         userID,
		Symbol:         brokerTrade.Symbol,
		MarketType:     data.MarketTypeIndian, // Default, can be enhanced based on exchange
		EntryDate:      exchangeTime,
		EntryPrice:     brokerTrade.Price,
		Quantity:       brokerTrade.Quantity,
		TotalAmount:    brokerTrade.Price * float64(brokerTrade.Quantity),
		ExitPrice:      nil, // Not available from broker data
		Direction:      direction,
		StopLoss:       nil,
		Target:         nil,
		Strategy:       "",                           // User needs to set this
		OutcomeSummary: data.OutcomeSummaryBreakeven, // Default, user can update
		TradeAnalysis:  nil,
		RulesFollowed:  []string{},
		Screenshots:    []string{},
		Psychology:     nil,
		// Broker-specific fields
		TradingBroker:   &broker,
		TraderBrokerID:  nil, // Can be set if available
		ExchangeOrderID: &brokerTrade.ExchangeOrderID,
		OrderID:         &brokerTrade.OrderID,
		ProductType:     productType,
		TransactionType: &transactionType,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return trade, nil
}

// convertProductType converts string product type to ProductType enum
func convertProductType(productType string) data.ProductType {
	switch productType {
	case "CNC":
		return data.ProductTypeCNC
	case "MIS":
		return data.ProductTypeMIS
	case "NRML":
		return data.ProductTypeNRML
	case "INTRADAY":
		return data.ProductTypeIntraday
	case "OTC":
		return data.ProductTypeOTC
	default:
		return data.ProductTypeCNC // Default
	}
}
