package repos

import (
	"database/sql"

	"github.com/akshitmadan/100xtrader/backend/internal/data"
)

type TickerRepository interface {
	SaveTicker(ticker *data.Ticker) error
	ListTickers() ([]*data.Ticker, error)
	DeleteTicker(symbol string) error
}

type tickerRepository struct {
	db *sql.DB
}

func NewTickerRepository(db *sql.DB) TickerRepository {
	return &tickerRepository{db: db}
}

func (r *tickerRepository) SaveTicker(ticker *data.Ticker) error {
	_, err := r.db.Exec(`INSERT OR REPLACE INTO tickers (symbol, name) VALUES (?, ?)`, ticker.Symbol, ticker.Name)
	return err
}

func (r *tickerRepository) ListTickers() ([]*data.Ticker, error) {
	rows, err := r.db.Query(`SELECT symbol, name FROM tickers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tickers []*data.Ticker
	for rows.Next() {
		var t data.Ticker
		if err := rows.Scan(&t.Symbol, &t.Name); err != nil {
			return nil, err
		}
		tickers = append(tickers, &t)
	}
	return tickers, nil
}

func (r *tickerRepository) DeleteTicker(symbol string) error {
	_, err := r.db.Exec(`DELETE FROM tickers WHERE symbol = ?`, symbol)
	return err
}
