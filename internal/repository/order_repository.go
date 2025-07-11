package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/codepnw/go-cart-system/internal/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, input *domain.Order) (*domain.Order, error)
	CreateOrderItems(ctx context.Context, orderID int64, items []*domain.OrderItem) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	query := `
		INSERT INTO orders (user_id, total_price, discount, final_price, coupon_code, status)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at;
	`
	var id int64
	var createdAt time.Time

	err := r.db.QueryRowContext(
		ctx,
		query,
		order.UserID,
		order.TotalPrice,
		order.Discount,
		order.FinalPrice,
		order.CouponCode,
		order.Status,
	).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	order.ID = id
	order.CreatedAt = createdAt

	return order, nil
}

func (r *orderRepository) CreateOrderItems(ctx context.Context, orderID int64, items []*domain.OrderItem) error {
	query := `
		INSERT INTO order_items 
			(order_id, product_id, product_name, quantity, unit_price, total_price) 
		VALUES
	`
	var args []any
	var values []string

	for i, item := range items {
		item.OrderID = orderID
		idx := i*6 + 1

		values = append(
			values,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", idx, idx+1, idx+2, idx+3, idx+4, idx+5),
		)
		args = append(
			args,
			item.OrderID,
			item.ProductID,
			item.ProductName,
			item.Quantity,
			item.UnitPrice,
			item.TotalPrice,
		)
	}
	query += strings.Join(values, ", ")

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
