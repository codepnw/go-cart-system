package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/codepnw/go-cart-system/internal/domain"
	"github.com/codepnw/go-cart-system/internal/dto"
	"github.com/codepnw/go-cart-system/internal/repository"
)

const queryTimeOut = time.Second * 5

type CartUsecase interface {
	AddItems(ctx context.Context, req *dto.CreateCartItems) error
	GetCart(ctx context.Context, cartID int64) ([]*dto.CartItemDetailsResponse, error)
	UpdateQuantity(ctx context.Context, input *dto.UpdateCartItems) error
	DeleteItem(ctx context.Context, id int64) error
}

type cartUsecase struct {
	repo repository.CartRepository
}

func NewCartUsecase(repo repository.CartRepository) CartUsecase {
	return &cartUsecase{repo: repo}
}

func (u *cartUsecase) AddItems(ctx context.Context, req *dto.CreateCartItems) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	if len(req.Items) == 0 {
		return errors.New("must provide at least 1 item")
	}

	var items []*domain.CartItems
	for _, i := range req.Items {
		if i.Quantity <= 0 {
			return errors.New("quantity must be greater than 0")
		}

		items = append(items, &domain.CartItems{
			CartID:    req.CartID,
			ProductID: i.ProductID,
			Quantity:  i.Quantity,
		})
	}

	return u.repo.AddItems(ctx, items)
}

func (u *cartUsecase) GetCart(ctx context.Context, cartID int64) ([]*dto.CartItemDetailsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.GetCart(ctx, cartID)
}

func (u *cartUsecase) UpdateQuantity(ctx context.Context, input *dto.UpdateCartItems) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	item := domain.CartItems{
		ID:       input.ID,
		Quantity: input.Quantity,
	}

	return u.repo.UpdateQuantity(ctx, &item)
}

func (u *cartUsecase) DeleteItem(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
