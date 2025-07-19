package repos

import (
	"context"
	"database/sql"

	"github.com/akshitmadan/100xtrader/backend/internal/data"
)

type PositionRepository interface {
	UpdatePosition(ctx context.Context, user, symbol string, quantity, price float64) error
	GetPositionsByUser(ctx context.Context, user string) ([]*data.Position, error)
}

type positionRepository struct {
	db *sql.DB
}

func NewPositionRepository(db *sql.DB) PositionRepository {
	return &positionRepository{db: db}
}

func (r *positionRepository) UpdatePosition(ctx context.Context, user, symbol string, quantity, price float64) error {
	// Try to update existing position
	res, err := r.db.ExecContext(ctx, `UPDATE positions SET quantity = quantity + ?, average_price = ((average_price * quantity) + (? * ?)) / (quantity + ?) WHERE user = ? AND symbol = ?`,
		quantity, price, quantity, quantity, user, symbol,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		// Insert new position
		_, err = r.db.ExecContext(ctx, `INSERT INTO positions (user, symbol, quantity, average_price) VALUES (?, ?, ?, ?)`,
			user, symbol, quantity, price,
		)
		return err
	}
	return nil
}

func (r *positionRepository) GetPositionsByUser(ctx context.Context, user string) ([]*data.Position, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user, symbol, quantity, average_price FROM positions WHERE user = ?`, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var positions []*data.Position
	for rows.Next() {
		var p data.Position
		if err := rows.Scan(&p.ID, &p.User, &p.Symbol, &p.Quantity, &p.AveragePrice); err != nil {
			return nil, err
		}
		positions = append(positions, &p)
	}
	return positions, nil
}
