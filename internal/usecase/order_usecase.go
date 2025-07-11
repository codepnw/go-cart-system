package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/codepnw/go-cart-system/internal/domain"
	"github.com/codepnw/go-cart-system/internal/repository"
)

type OrderUsecase interface {
	Checkout(ctx context.Context, userID int64, couponCode string) error
}

type orderUsecase struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewOrderUsecase(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository) OrderUsecase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (uc *orderUsecase) Checkout(ctx context.Context, userID int64, couponCode string) error {
	// Get Cart Items
	cartItems, err := uc.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return err
	}
	if len(cartItems) == 0 {
		return errors.New("cart is empty")
	}

	var totalPrice float64
	var orderItems []*domain.OrderItem

	for _, item := range cartItems {
		product, err := uc.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("get product failed: %w", err)
		}
		if product.Stock < item.Quantity {
			return fmt.Errorf("product %s out of stock", product.Name)
		}

		itemPrice := product.Price * float64(item.Quantity)
		totalPrice += itemPrice

		orderItems = append(orderItems, &domain.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   product.Price,
			TotalPrice:  itemPrice,
		})
	}

	// TODO Discount

	discount := 0
	finalPrice := totalPrice - float64(discount)

	// Create Order
	order := &domain.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Discount:   float64(discount),
		FinalPrice: finalPrice,
		CouponCode: couponCode,
		Status:     domain.Pending,
		Items:      orderItems,
	}
	order, err = uc.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	// Create Order Items
	err = uc.orderRepo.CreateOrderItems(ctx, order.ID, orderItems)
	if err != nil {
		return err
	}

	// TODO Decrease Stock

	// TODO Clear Cart

	return nil
}
