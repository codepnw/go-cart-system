package repository

import (
	"context"
	"database/sql"

	"github.com/codepnw/go-cart-system/internal/domain"
	"github.com/codepnw/go-cart-system/internal/errs"
)

type ProductRepository interface {
	Create(ctx context.Context, input *domain.Product) error
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	List(ctx context.Context) ([]*domain.Product, error)
	Update(ctx context.Context, input *domain.Product) error
	Delete(ctx context.Context, id int64) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, input *domain.Product) error {
	query := `
		INSERT INTO products (name, price) VALUES ($1, $2)
	`
	res, err := r.db.ExecContext(ctx, query, input.Name, input.Price)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT id, name, price, created_at, updated_at
		FROM products WHERE id = $1
	`
	var product domain.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	return &product, err
}

func (r *productRepository) List(ctx context.Context) ([]*domain.Product, error) {
	query := `
		SELECT id, name, price, created_at, updated_at
		FROM products
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var products []*domain.Product

	for rows.Next() {
		var p domain.Product
		err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) Update(ctx context.Context, input *domain.Product) error {
	query := `
		UPDATE products SET name = $1, price = $2, updated_at = NOW()
		WHERE id = $3
	`
	res, err := r.db.ExecContext(ctx, query, input.Name, input.Price, input.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}
