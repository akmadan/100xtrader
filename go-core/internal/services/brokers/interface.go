package brokers

import (
	"go-core/internal/data"
)

// BrokerTrade represents a trade from a third-party broker
// This is the common format that all brokers should convert their data to
type BrokerTrade struct {
	Symbol          string
	Quantity        int
	Price           float64
	TransactionType string // "buy" | "sell"
	ExchangeOrderID string
	OrderID         string
	ProductType     string // "CNC" | "MIS" | "NRML" | "INTRADAY" | "OTC"
	ExchangeTime    string // ISO format timestamp
}

// BrokerService defines the interface that all broker services must implement
type BrokerService interface {
	// ParseTrades parses raw broker data and converts it to BrokerTrade format
	ParseTrades(rawData []byte) ([]BrokerTrade, error)

	// ConvertToTrade converts BrokerTrade to the internal Trade model
	ConvertToTrade(brokerTrade BrokerTrade, userID int) (*data.Trade, error)

	// GetBrokerName returns the name of the broker (e.g., "zerodha", "dhan")
	GetBrokerName() data.TradingBroker
}
