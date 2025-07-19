package repos

import (
	"context"
	"database/sql"

	"github.com/akshitmadan/100xtrader/backend/internal/data"
)

type TradeRepository interface {
	RecordTrade(ctx context.Context, trade *data.Trade) error
	GetTradesByUser(ctx context.Context, user string) ([]*data.Trade, error)
}

type tradeRepository struct {
	db *sql.DB
}

func NewTradeRepository(db *sql.DB) TradeRepository {
	return &tradeRepository{db: db}
}

func (r *tradeRepository) RecordTrade(ctx context.Context, trade *data.Trade) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO trades (buy_order_id, sell_order_id, symbol, quantity, price, timestamp) VALUES (?, ?, ?, ?, ?, ?)`,
		trade.BuyOrderID, trade.SellOrderID, trade.Symbol, trade.Quantity, trade.Price, trade.Timestamp,
	)
	return err
}

func (r *tradeRepository) GetTradesByUser(ctx context.Context, user string) ([]*data.Trade, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, buy_order_id, sell_order_id, symbol, quantity, price, timestamp FROM trades WHERE buy_order_id IN (SELECT id FROM orders WHERE user = ?) OR sell_order_id IN (SELECT id FROM orders WHERE user = ?) ORDER BY timestamp DESC`, user, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var trades []*data.Trade
	for rows.Next() {
		var t data.Trade
		if err := rows.Scan(&t.ID, &t.BuyOrderID, &t.SellOrderID, &t.Symbol, &t.Quantity, &t.Price, &t.Timestamp); err != nil {
			return nil, err
		}
		trades = append(trades, &t)
	}
	return trades, nil
}
