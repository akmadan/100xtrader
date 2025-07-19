package engine

import (
	"sync"
	"time"
)

const (
	defaultOpen      = 2800.0
	defaultMaxTrades = 100
)

type TradeTick struct {
	Price     float64
	Quantity  float64
	Side      string
	Timestamp time.Time
}

type OHLC struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
	Start  time.Time
}

type MarketDataCache struct {
	mu           sync.RWMutex
	RecentTrades []TradeTick
	OHLC         OHLC
}

var (
	marketDataCache = make(map[string]*MarketDataCache)
	marketDataOnce  sync.Once
)

// StartMarketSimulation starts a background goroutine to simulate trades for a symbol
func StartMarketSimulation(symbol string) {
	go func() {
		price := defaultOpen
		for {
			// Simulate random walk price
			move := (0.5 - (float64(time.Now().UnixNano()%1000) / 1000.0)) * 2.0 // -1 to +1
			price += move
			if price < 1000 {
				price = 1000
			}
			if price > 4000 {
				price = 4000
			}
			qty := float64(1 + time.Now().UnixNano()%10)
			side := "buy"
			if time.Now().UnixNano()%2 == 0 {
				side = "sell"
			}
			AddTradeToMarketData(symbol, price, qty, side, time.Now())
			time.Sleep(1 * time.Second)
		}
	}()
}

// AddTradeToMarketData updates the cache with a new trade
func AddTradeToMarketData(symbol string, price, qty float64, side string, ts time.Time) {
	cache, ok := marketDataCache[symbol]
	if !ok {
		return
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()
	// Update recent trades
	tt := TradeTick{Price: price, Quantity: qty, Side: side, Timestamp: ts}
	cache.RecentTrades = append([]TradeTick{tt}, cache.RecentTrades...)
	if len(cache.RecentTrades) > defaultMaxTrades {
		cache.RecentTrades = cache.RecentTrades[:defaultMaxTrades]
	}
	// Update OHLC for current 1m candle
	minute := ts.Truncate(time.Minute)
	if !minute.Equal(cache.OHLC.Start) {
		// New candle
		cache.OHLC = OHLC{Open: price, High: price, Low: price, Close: price, Volume: qty, Start: minute}
	} else {
		if price > cache.OHLC.High {
			cache.OHLC.High = price
		}
		if price < cache.OHLC.Low {
			cache.OHLC.Low = price
		}
		cache.OHLC.Close = price
		cache.OHLC.Volume += qty
	}
}

// GetMarketDataSnapshot returns the current market data for a symbol
func GetMarketDataSnapshot(symbol string) (lastPrice float64, lastTradeTime time.Time, ohlc OHLC, trades []TradeTick) {
	cache, ok := marketDataCache[symbol]
	if !ok {
		return 0, time.Time{}, OHLC{}, nil
	}
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	if len(cache.RecentTrades) > 0 {
		lastPrice = cache.RecentTrades[0].Price
		lastTradeTime = cache.RecentTrades[0].Timestamp
	}
	trades = append([]TradeTick(nil), cache.RecentTrades...)
	return lastPrice, lastTradeTime, cache.OHLC, trades
}
