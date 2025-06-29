package usecase

import (
	"context"
	"time"

	"github.com/codepnw/go-cart-system/internal/domain"
	"github.com/codepnw/go-cart-system/internal/dto"
	"github.com/codepnw/go-cart-system/internal/repository"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *dto.CreateProduct) error
	GetProduct(ctx context.Context, id int64) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]*domain.Product, error)
	UpdateProduct(ctx context.Context, req *dto.UpdateProduct) error
	DeleteProduct(ctx context.Context, id int64) error
}

type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}

func (uc *productUsecase) CreateProduct(ctx context.Context, req *dto.CreateProduct) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	product := domain.Product{
		Name:  req.Name,
		Price: req.Price,
	}

	return uc.repo.Create(ctx, &product)
}

func (uc *productUsecase) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return uc.repo.GetByID(ctx, id)
}

func (uc *productUsecase) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return uc.repo.List(ctx)
}

func (uc *productUsecase) UpdateProduct(ctx context.Context, req *dto.UpdateProduct) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	product, err := uc.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if req.Name != nil {
		product.Name = *req.Name
	}

	if req.Price != nil {
		product.Price = *req.Price
	}

	now := time.Now()
	product.UpdatedAt = &now

	return uc.repo.Update(ctx, product)
}

func (uc *productUsecase) DeleteProduct(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return uc.repo.Delete(ctx, id)
}
