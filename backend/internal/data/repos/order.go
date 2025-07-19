package repos

import (
	"context"
	"database/sql"
	"time"

	"github.com/akshitmadan/100xtrader/backend/internal/data"
)

// OrderRepository defines the interface for order data access
// This makes it easy to mock for testing

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *data.Order) (int64, error)
	ListOrders(ctx context.Context) ([]*data.Order, error)
	DeleteOrder(ctx context.Context, id int64) error
	GetOrderByID(ctx context.Context, id int64) (*data.Order, error)
}

// orderRepository is the concrete implementation

type orderRepository struct {
	db *sql.DB
}

// NewOrderRepository returns a new OrderRepository
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *data.Order) (int64, error) {
	query := `INSERT INTO orders (user, symbol, side, type, quantity, price, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		order.User, order.Symbol, order.Side, order.Type, order.Quantity, order.Price, order.Status, order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *orderRepository) ListOrders(ctx context.Context) ([]*data.Order, error) {
	query := `SELECT id, user, symbol, side, type, quantity, price, status, created_at, updated_at FROM orders ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*data.Order
	for rows.Next() {
		var o data.Order
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&o.ID, &o.User, &o.Symbol, &o.Side, &o.Type, &o.Quantity, &o.Price, &o.Status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		o.CreatedAt = createdAt
		o.UpdatedAt = updatedAt
		orders = append(orders, &o)
	}
	return orders, nil
}

func (r *orderRepository) DeleteOrder(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM orders WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, id int64) (*data.Order, error) {
	query := `SELECT id, user, symbol, side, type, quantity, price, status, created_at, updated_at FROM orders WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	var o data.Order
	var createdAt, updatedAt time.Time
	if err := row.Scan(&o.ID, &o.User, &o.Symbol, &o.Side, &o.Type, &o.Quantity, &o.Price, &o.Status, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	o.CreatedAt = createdAt
	o.UpdatedAt = updatedAt
	return &o, nil
}
