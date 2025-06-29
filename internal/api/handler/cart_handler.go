package handler

import (
	"net/http"

	"github.com/codepnw/go-cart-system/internal/dto"
	"github.com/codepnw/go-cart-system/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type cartHandler struct {
	uc       usecase.CartUsecase
	validate *validator.Validate
}

func NewCartHandler(uc usecase.CartUsecase) *cartHandler {
	return &cartHandler{
		uc:       uc,
		validate: validator.New(),
	}
}

func (h *cartHandler) AddItems(ctx *fiber.Ctx) error {
	var req dto.CreateCartItems

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

	if err := h.uc.AddItems(ctx.Context(), &req); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{
		"message": "items added to cart",
	})
}

func (h *cartHandler) GetCart(ctx *fiber.Ctx) error {
	cartID, err := ctx.ParamsInt("cartID")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	items, err := h.uc.GetCart(ctx.Context(), int64(cartID))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "cart items",
		"data":    items,
	})
}

func (h *cartHandler) UpdateQuantity(ctx *fiber.Ctx) error {
	var req dto.UpdateCartItems

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

	if err := h.uc.UpdateQuantity(ctx.Context(), &req); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "quantity updated",
	})
}

func (h *cartHandler) DeleteItem(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.uc.DeleteItem(ctx.Context(), int64(id)); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusNoContent).JSON(nil)
}
