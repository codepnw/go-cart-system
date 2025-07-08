package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/go-cart-system/internal/domain"
	"github.com/codepnw/go-cart-system/internal/dto"
	"github.com/codepnw/go-cart-system/internal/repository"
	"github.com/codepnw/go-cart-system/internal/utils/errs"
)

const queryTimeOut = time.Second * 5

type CartUsecase interface {
	AddItems(ctx context.Context, req *dto.CreateCartItems) error
	GetCart(ctx context.Context, userID int64) (*dto.CartResponse, error)
	UpdateQuantity(ctx context.Context, userID int64, items []dto.UpdateCartItem) error
	DeleteItem(ctx context.Context, itemID int64) error
}

type cartUsecase struct {
	cartRepo repository.CartRepository
	prodRepo repository.ProductRepository
}

func NewCartUsecase(cartRepo repository.CartRepository, prodRepo repository.ProductRepository) CartUsecase {
	return &cartUsecase{
		cartRepo: cartRepo,
		prodRepo: prodRepo,
	}
}

func (uc *cartUsecase) AddItems(ctx context.Context, req *dto.CreateCartItems) error {
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

	return uc.cartRepo.AddItems(ctx, items)
}

func (uc *cartUsecase) GetCart(ctx context.Context, userID int64) (*dto.CartResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	items, err := uc.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalItems int
	var totalPrice float64

	for _, item := range items {
		totalItems += item.Quantity
		totalPrice += float64(item.Quantity) * item.Price
	}

	return &dto.CartResponse{
		Items:      items,
		TotalItems: totalItems,
		TotalPrice: totalPrice,
	}, nil
}

func (uc *cartUsecase) UpdateQuantity(ctx context.Context, userID int64, items []dto.UpdateCartItem) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	for _, item := range items {
		// check cartItem
		cartItem, err := uc.cartRepo.GetCartItem(ctx, userID, item.ProductID)
		if err != nil {
			return fmt.Errorf("get cart item failed: %w", err)
		}
		if cartItem == nil {
			continue
		}

		// find product
		product, err := uc.prodRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("get product failed: %w", err)
		}
		if product == nil {
			return errs.ErrProductNotFound
		}

		// check stock
		if item.Quantity > product.Stock {
			return errors.New("product not enough")
		}

		if item.Quantity == 0 {
			if err := uc.cartRepo.DeleteCartItem(ctx, cartItem.ID); err != nil {
				return fmt.Errorf("delete cart item failed: %w", err)
			}
			continue
		}

		// update quantity
		err = uc.cartRepo.UpdateQuantity(ctx, &domain.CartItems{
			ID:        cartItem.ID,
			Quantity:  item.Quantity,
			ProductID: product.ID,
		})
		if err != nil {
			return fmt.Errorf("update quantity failed: %w", err)
		}
	}

	return nil
}

func (uc *cartUsecase) DeleteItem(ctx context.Context, itemID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return uc.cartRepo.DeleteCartItem(ctx, itemID)
}
