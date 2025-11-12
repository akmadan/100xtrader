package brokers

import (
	"fmt"

	"go-core/internal/data"
)

// GetBrokerService returns the appropriate broker service based on broker name
func GetBrokerService(brokerName data.TradingBroker) (BrokerService, error) {
	switch brokerName {
	case data.TradingBrokerZerodha:
		return NewZerodhaService(), nil
	case data.TradingBrokerDhan:
		return NewDhanService(), nil
	default:
		return nil, fmt.Errorf("unsupported broker: %s", brokerName)
	}
}

// ImportTrades imports trades from broker data
// This is a convenience function that handles the full import process
func ImportTrades(
	brokerName data.TradingBroker,
	rawData []byte,
	userID int,
) ([]*data.Trade, error) {
	// Get the appropriate broker service
	brokerService, err := GetBrokerService(brokerName)
	if err != nil {
		return nil, err
	}

	// Parse trades from broker data
	brokerTrades, err := brokerService.ParseTrades(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse broker trades: %w", err)
	}

	// Convert to internal Trade model
	trades := make([]*data.Trade, 0, len(brokerTrades))
	for _, brokerTrade := range brokerTrades {
		trade, err := brokerService.ConvertToTrade(brokerTrade, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert trade: %w", err)
		}
		trades = append(trades, trade)
	}

	return trades, nil
}
