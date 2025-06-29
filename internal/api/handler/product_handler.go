package handler

import (
	"net/http"

	"github.com/codepnw/go-cart-system/internal/dto"
	"github.com/codepnw/go-cart-system/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

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
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.uc.CreateProduct(ctx.Context(), &req); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{
		"message": "new product created",
	})
}

func (h *productHandler) GetProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	product, err := h.uc.GetProduct(ctx.Context(), int64(id))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"data": product,
	})
}

func (h *productHandler) ListProducts(ctx *fiber.Ctx) error {
	products, err := h.uc.ListProducts(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"data": products,
	})
}

func (h *productHandler) UpdateProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	var req dto.UpdateProduct

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	req.ID = int64(id)
	if err := h.uc.UpdateProduct(ctx.Context(), &req); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "product updated",
	})
}

func (h *productHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if err = h.uc.DeleteProduct(ctx.Context(), int64(id)); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "product deleted",
	})
}
