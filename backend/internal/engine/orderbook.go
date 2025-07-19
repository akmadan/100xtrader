package engine

import (
	"container/heap"
	"sync"
	"time"

	"github.com/akshitmadan/100xtrader/backend/internal/data"
)

// OrderBook represents the order book for a single symbol
type OrderBook struct {
	Symbol     string
	bids       *OrderQueue // Max heap for bids (highest price first)
	asks       *OrderQueue // Min heap for asks (lowest price first)
	mu         sync.RWMutex
	stopOrders []*data.Order // Stop orders for this symbol
}

// OrderQueue is a priority queue for orders
type OrderQueue struct {
	orders []*data.Order
	heap.Interface
}

// MarketDepth represents the market depth at a price level
type MarketDepth struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
	Orders   int     `json:"orders"`
}

// OrderBookSnapshot represents a snapshot of the order book
type OrderBookSnapshot struct {
	Symbol string        `json:"symbol"`
	Bids   []MarketDepth `json:"bids"`
	Asks   []MarketDepth `json:"asks"`
	Time   time.Time     `json:"time"`
}

// NewOrderBook creates a new order book for a symbol
func NewOrderBook(symbol string) *OrderBook {
	ob := &OrderBook{
		Symbol: symbol,
		bids:   &OrderQueue{orders: make([]*data.Order, 0)},
		asks:   &OrderQueue{orders: make([]*data.Order, 0)},
	}
	heap.Init(ob.bids)
	heap.Init(ob.asks)
	return ob
}

// AddOrder adds an order to the order book
func (ob *OrderBook) AddOrder(order *data.Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if order.Side == "buy" {
		heap.Push(ob.bids, order)
	} else {
		heap.Push(ob.asks, order)
	}
}

// RemoveOrder removes an order from the order book
func (ob *OrderBook) RemoveOrder(order *data.Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	// Note: This is a simplified implementation
	// In production, you'd want to maintain an index for O(1) removal
	if order.Side == "buy" {
		ob.removeFromQueue(ob.bids, order.ID)
	} else {
		ob.removeFromQueue(ob.asks, order.ID)
	}
}

// AddStopOrder adds a stop order to the order book's stop order list
func (ob *OrderBook) AddStopOrder(order *data.Order) {
	ob.stopOrders = append(ob.stopOrders, order)
}

// CheckAndTriggerStopOrders checks if any stop orders should be triggered based on the last trade price
// If triggered, converts them to market orders and returns the triggered orders
func (ob *OrderBook) CheckAndTriggerStopOrders(lastPrice float64) []*data.Order {
	var triggered []*data.Order
	remaining := ob.stopOrders[:0]
	for _, order := range ob.stopOrders {
		trigger := false
		if order.Side == "buy" && lastPrice >= order.Price {
			trigger = true
		}
		if order.Side == "sell" && lastPrice <= order.Price {
			trigger = true
		}
		if trigger {
			order.Type = "market" // Convert to market order
			triggered = append(triggered, order)
		} else {
			remaining = append(remaining, order)
		}
	}
	ob.stopOrders = remaining
	return triggered
}

// GetBestBid returns the best bid (highest price)
func (ob *OrderBook) GetBestBid() *data.Order {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.bids.Len() == 0 {
		return nil
	}
	return ob.bids.orders[0]
}

// GetBestAsk returns the best ask (lowest price)
func (ob *OrderBook) GetBestAsk() *data.Order {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.asks.Len() == 0 {
		return nil
	}
	return ob.asks.orders[0]
}

// GetMarketDepth returns the market depth for bids and asks
func (ob *OrderBook) GetMarketDepth(levels int) OrderBookSnapshot {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	snapshot := OrderBookSnapshot{
		Symbol: ob.Symbol,
		Time:   time.Now().UTC(),
		Bids:   make([]MarketDepth, 0),
		Asks:   make([]MarketDepth, 0),
	}

	// Aggregate bids by price level
	bidLevels := make(map[float64]*MarketDepth)
	for i := 0; i < ob.bids.Len() && len(bidLevels) < levels; i++ {
		order := ob.bids.orders[i]
		if depth, exists := bidLevels[order.Price]; exists {
			depth.Quantity += order.Quantity
			depth.Orders++
		} else {
			bidLevels[order.Price] = &MarketDepth{
				Price:    order.Price,
				Quantity: order.Quantity,
				Orders:   1,
			}
		}
	}

	// Aggregate asks by price level
	askLevels := make(map[float64]*MarketDepth)
	for i := 0; i < ob.asks.Len() && len(askLevels) < levels; i++ {
		order := ob.asks.orders[i]
		if depth, exists := askLevels[order.Price]; exists {
			depth.Quantity += order.Quantity
			depth.Orders++
		} else {
			askLevels[order.Price] = &MarketDepth{
				Price:    order.Price,
				Quantity: order.Quantity,
				Orders:   1,
			}
		}
	}

	// Convert maps to slices and sort
	for _, depth := range bidLevels {
		snapshot.Bids = append(snapshot.Bids, *depth)
	}
	for _, depth := range askLevels {
		snapshot.Asks = append(snapshot.Asks, *depth)
	}

	return snapshot
}

// Helper method to remove order from queue (simplified)
func (ob *OrderBook) removeFromQueue(queue *OrderQueue, orderID int64) {
	for i, order := range queue.orders {
		if order.ID == orderID {
			heap.Remove(queue, i)
			break
		}
	}
}

// OrderQueue heap implementation
func (oq *OrderQueue) Len() int { return len(oq.orders) }

func (oq *OrderQueue) Less(i, j int) bool {
	// For bids: higher price first, then earlier time
	// For asks: lower price first, then earlier time
	if oq.orders[i].Price != oq.orders[j].Price {
		return oq.orders[i].Price > oq.orders[j].Price // Max heap for bids
	}
	return oq.orders[i].CreatedAt.Before(oq.orders[j].CreatedAt)
}

func (oq *OrderQueue) Swap(i, j int) {
	oq.orders[i], oq.orders[j] = oq.orders[j], oq.orders[i]
}

func (oq *OrderQueue) Push(x interface{}) {
	oq.orders = append(oq.orders, x.(*data.Order))
}

func (oq *OrderQueue) Pop() interface{} {
	old := oq.orders
	n := len(old)
	item := old[n-1]
	oq.orders = old[0 : n-1]
	return item
}
