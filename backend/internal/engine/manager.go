package engine

import (
	"sync"
	"time"

	"github.com/akshitmadan/100xtrader/backend/internal/data"
)

// OrderBookManager manages order books for all symbols (singleton)
type OrderBookManager struct {
	books map[string]*OrderBook
	mu    sync.RWMutex
}

var (
	orderBookManager     *OrderBookManager
	orderBookManagerOnce sync.Once
)

// GetOrderBookManager returns the singleton instance
func GetOrderBookManager() *OrderBookManager {
	orderBookManagerOnce.Do(func() {
		orderBookManager = &OrderBookManager{
			books: make(map[string]*OrderBook),
		}
	})
	return orderBookManager
}

// GetOrCreateOrderBook returns the order book for a symbol, creating it if needed
func (m *OrderBookManager) GetOrCreateOrderBook(symbol string) *OrderBook {
	m.mu.Lock()
	defer m.mu.Unlock()
	ob, exists := m.books[symbol]
	if !exists {
		ob = NewOrderBook(symbol)
		m.books[symbol] = ob
	}
	return ob
}

// AddOrderAndMatch adds an order to the book and runs the matcher, returning executed trades
func (m *OrderBookManager) AddOrderAndMatch(order *data.Order) []*data.Trade {
	ob := m.GetOrCreateOrderBook(order.Symbol)
	ob.AddOrder(order)
	return ob.MatchOrders()
}

// AddMarketOrderAndMatch matches a market order immediately against the book and returns executed trades
func (m *OrderBookManager) AddMarketOrderAndMatch(order *data.Order) []*data.Trade {
	ob := m.GetOrCreateOrderBook(order.Symbol)
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var trades []*data.Trade
	remaining := order.Quantity

	if order.Side == "buy" {
		// Match against asks (lowest price first)
		for ob.asks.Len() > 0 && remaining > 0 {
			ask := ob.asks.orders[0]
			qty := min(remaining, ask.Quantity)
			trade := &data.Trade{
				BuyOrderID:  order.ID,
				SellOrderID: ask.ID,
				Symbol:      order.Symbol,
				Quantity:    qty,
				Price:       ask.Price,
				Timestamp:   time.Now().UTC(),
			}
			trades = append(trades, trade)
			remaining -= qty
			ask.Quantity -= qty
			if ask.Quantity == 0 {
				heapRemove(ob.asks, 0)
			}
		}
	} else {
		// Match against bids (highest price first)
		for ob.bids.Len() > 0 && remaining > 0 {
			bid := ob.bids.orders[0]
			qty := min(remaining, bid.Quantity)
			trade := &data.Trade{
				BuyOrderID:  bid.ID,
				SellOrderID: order.ID,
				Symbol:      order.Symbol,
				Quantity:    qty,
				Price:       bid.Price,
				Timestamp:   time.Now().UTC(),
			}
			trades = append(trades, trade)
			remaining -= qty
			bid.Quantity -= qty
			if bid.Quantity == 0 {
				heapRemove(ob.bids, 0)
			}
		}
	}

	return trades
}

// AddStopOrderAndCheck triggers stop orders for a symbol after a trade
func (m *OrderBookManager) AddStopOrderAndCheck(order *data.Order, lastPrice float64) []*data.Order {
	ob := m.GetOrCreateOrderBook(order.Symbol)
	ob.AddStopOrder(order)
	return ob.CheckAndTriggerStopOrders(lastPrice)
}
