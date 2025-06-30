package handler

import (
	"github.com/codepnw/go-cart-system/internal/api/response"
	"github.com/codepnw/go-cart-system/internal/dto"
	"github.com/codepnw/go-cart-system/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const paramIDKey = "id"

type productHandler struct {
	uc       usecase.ProductUsecase
	validate *validator.Validate
}

func NewProductHandler(uc usecase.ProductUsecase) *productHandler {
	return &productHandler{
		uc:       uc,
		validate: validator.New(),
	}
}

func (h *productHandler) CreateProduct(ctx *fiber.Ctx) error {
	var req dto.CreateProduct

	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validate.Struct(req); err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	if err := h.uc.CreateProduct(ctx.Context(), &req); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.CreastedResponse(ctx, "product created", nil)
}

func (h *productHandler) GetProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt(paramIDKey)
	if err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	product, err := h.uc.GetProduct(ctx.Context(), int64(id))
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.SuccessResponse(ctx, "", product)
}

func (h *productHandler) ListProducts(ctx *fiber.Ctx) error {
	products, err := h.uc.ListProducts(ctx.Context())
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.SuccessResponse(ctx, "", products)
}

func (h *productHandler) UpdateProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt(paramIDKey)
	if err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	var req dto.UpdateProduct

	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validate.Struct(req); err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	req.ID = int64(id)
	if err := h.uc.UpdateProduct(ctx.Context(), &req); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.SuccessResponse(ctx, "product updated", nil)
}

func (h *productHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt(paramIDKey)
	if err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	if err = h.uc.DeleteProduct(ctx.Context(), int64(id)); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.SuccessResponse(ctx, "", nil)
}
