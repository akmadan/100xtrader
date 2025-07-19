package engine

import (
	"sync"
	"time"

	"github.com/akshitmadan/100xtrader/backend/internal/data"
	"github.com/akshitmadan/100xtrader/backend/internal/data/repos"
)

type TickerManager struct {
	mu      sync.RWMutex
	tickers map[string]*data.Ticker
	repo    repos.TickerRepository
}

var (
	tickerManager     *TickerManager
	tickerManagerOnce sync.Once
)

// GetTickerManager returns the singleton instance
func GetTickerManager(repo repos.TickerRepository) *TickerManager {
	tickerManagerOnce.Do(func() {
		tickerManager = &TickerManager{
			tickers: make(map[string]*data.Ticker),
			repo:    repo,
		}
		tickerManager.LoadTickersFromRepo()
	})
	return tickerManager
}

// AddTicker registers a new ticker, initializes market data and order book, and persists it
func (m *TickerManager) AddTicker(ticker *data.Ticker) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tickers[ticker.Symbol] = ticker
	if err := m.repo.SaveTicker(ticker); err != nil {
		return err
	}
	// Initialize market data and order book
	if _, ok := marketDataCache[ticker.Symbol]; !ok {
		marketDataCache[ticker.Symbol] = &MarketDataCache{
			RecentTrades: make([]TradeTick, 0, defaultMaxTrades),
			OHLC: OHLC{
				Open:   defaultOpen,
				High:   defaultOpen,
				Low:    defaultOpen,
				Close:  defaultOpen,
				Volume: 0,
				Start:  nowMinute(),
			},
		}
	}
	GetOrderBookManager().GetOrCreateOrderBook(ticker.Symbol)
	return nil
}

// RemoveTicker removes a ticker and cleans up resources and DB
func (m *TickerManager) RemoveTicker(symbol string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tickers, symbol)
	delete(marketDataCache, symbol)
	obm := GetOrderBookManager()
	obm.mu.Lock()
	delete(obm.books, symbol)
	obm.mu.Unlock()
	return m.repo.DeleteTicker(symbol)
}

// ListTickers returns all active tickers
func (m *TickerManager) ListTickers() []*data.Ticker {
	m.mu.RLock()
	defer m.mu.RUnlock()
	syms := make([]*data.Ticker, 0, len(m.tickers))
	for _, t := range m.tickers {
		syms = append(syms, t)
	}
	return syms
}

// LoadTickersFromRepo loads tickers from the repository
func (m *TickerManager) LoadTickersFromRepo() error {
	tickers, err := m.repo.ListTickers()
	if err != nil {
		return err
	}
	for _, t := range tickers {
		m.tickers[t.Symbol] = t
	}
	return nil
}

// nowMinute returns the current time truncated to the minute
func nowMinute() time.Time {
	return time.Now().Truncate(time.Minute)
}
