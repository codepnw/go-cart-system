package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/codepnw/go-cart-system/internal/domain"
	"github.com/codepnw/go-cart-system/internal/dto"
	"github.com/codepnw/go-cart-system/internal/utils/errs"
)

type CartRepository interface {
	AddItems(ctx context.Context, items []*domain.CartItems) error
	GetCart(ctx context.Context, userID int64) ([]*dto.CartItem, error)
	UpdateQuantity(ctx context.Context, input *domain.CartItems) error
	Delete(ctx context.Context, id int64) error
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddItems(ctx context.Context, items []*domain.CartItems) error {
	var (
		values []string
		args   []any
	)
	for i, item := range items {
		n := i * 3
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", n+1, n+2, n+3))
		args = append(args, item.CartID, item.ProductID, item.Quantity)
	}

	query := fmt.Sprintf(`
		INSERT INTO cart_items (cart_id, product_id, quantity) VALUES %s`,
		strings.Join(values, ","),
	)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrCartItemNotFound
	}

	return nil
}

func (r *cartRepository) GetCart(ctx context.Context, userID int64) ([]*dto.CartItem, error) {
	query := `
		SELECT 
			ci.id AS cart_item_id, 
			ci.product_id,
			p.name AS product_name,
			p.price,
			ci.quantity
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		JOIN carts c ON ci.cart_id = c.id
		WHERE c.user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*dto.CartItem

	for rows.Next() {
		var item dto.CartItem
		err = rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ProductName,
			&item.Price,
			&item.Quantity,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}

func (r *cartRepository) UpdateQuantity(ctx context.Context, input *domain.CartItems) error {
	query := `UPDATE cart_items SET quantity = $1 WHERE id = $2 AND product_id = $3`

	res, err := r.db.ExecContext(ctx, query, input.Quantity, input.ID, input.ProductID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrCartItemNotFound
	}

	return nil
}

func (r *cartRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM cart_items WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrCartItemNotFound
	}

	return nil
}
