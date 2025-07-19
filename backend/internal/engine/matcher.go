package engine

import (
	"time"

	"github.com/akshitmadan/100xtrader/backend/internal/data"
)

// MatchOrders matches buy and sell orders in the order book and returns executed trades
func (ob *OrderBook) MatchOrders() []*data.Trade {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var trades []*data.Trade

	for ob.bids.Len() > 0 && ob.asks.Len() > 0 {
		buy := ob.bids.orders[0]
		sell := ob.asks.orders[0]

		// Check if orders cross (buy price >= sell price)
		if buy.Price < sell.Price {
			break
		}

		// Determine trade quantity
		qty := min(buy.Quantity, sell.Quantity)
		tradePrice := sell.Price // Price is determined by the resting order (sell)

		trade := &data.Trade{
			BuyOrderID:  buy.ID,
			SellOrderID: sell.ID,
			Symbol:      ob.Symbol,
			Quantity:    qty,
			Price:       tradePrice,
			Timestamp:   time.Now().UTC(),
		}
		trades = append(trades, trade)

		// Update order quantities
		buy.Quantity -= qty
		sell.Quantity -= qty

		// Remove filled orders
		if buy.Quantity == 0 {
			heapRemove(ob.bids, 0)
		}
		if sell.Quantity == 0 {
			heapRemove(ob.asks, 0)
		}
	}

	return trades
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// heapRemove removes the order at index i from the queue
func heapRemove(q *OrderQueue, i int) {
	// Use the heap.Remove interface
	q.Interface = q
	heapRemoveImpl(q, i)
}

// heapRemoveImpl is a helper for heap.Remove
func heapRemoveImpl(q *OrderQueue, i int) {
	// Standard heap.Remove implementation
	last := len(q.orders) - 1
	q.Swap(i, last)
	q.orders = q.orders[:last]
	if i < len(q.orders) {
		heapify(q, i)
	}
}

// heapify restores the heap property
func heapify(q *OrderQueue, i int) {
	for {
		left := 2*i + 1
		right := 2*i + 2
		smallest := i
		if left < len(q.orders) && q.Less(left, smallest) {
			smallest = left
		}
		if right < len(q.orders) && q.Less(right, smallest) {
			smallest = right
		}
		if smallest == i {
			break
		}
		q.Swap(i, smallest)
		i = smallest
	}
}
