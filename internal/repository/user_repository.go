package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/codepnw/go-cart-system/internal/domain"
	"github.com/codepnw/go-cart-system/internal/utils/errs"
)

type UserRepository interface {
	Create(ctx context.Context, input *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, input *domain.User) error {
	query := `
		INSERT INTO users (email, password, role) VALUES ($1, $2, $3)
	`
	res, err := r.db.ExecContext(ctx, query, input.Email, input.Password, input.Role)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) selectQuery(where string) string {
	return fmt.Sprintf(`
		SELECT id, email, password, role, created_at, updated_at 
		FROM users WHERE %s = $1`, where)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := r.selectQuery("email")
	user := domain.User{}

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := r.selectQuery("id")
	user := domain.User{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
