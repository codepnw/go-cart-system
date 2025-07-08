package handler

import (
	"github.com/codepnw/go-cart-system/internal/api/middleware"
	"github.com/codepnw/go-cart-system/internal/api/response"
	"github.com/codepnw/go-cart-system/internal/dto"
	"github.com/codepnw/go-cart-system/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type cartHandler struct {
	mid      *middleware.Middleware
	uc       usecase.CartUsecase
	validate *validator.Validate
}

func NewCartHandler(uc usecase.CartUsecase, mid *middleware.Middleware) *cartHandler {
	return &cartHandler{
		mid:      mid,
		uc:       uc,
		validate: validator.New(),
	}
}

func (h *cartHandler) AddItems(ctx *fiber.Ctx) error {
	var req dto.CreateCartItems

	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validate.Struct(req); err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	if err := h.uc.AddItems(ctx.Context(), &req); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.CreastedResponse(ctx, "", nil)
}

func (h *cartHandler) GetCart(ctx *fiber.Ctx) error {
	user, err := h.mid.GetCurrentUser(ctx)
	if err != nil {
		return response.UnauthorizedResponse(ctx, err.Error())
	}

	items, err := h.uc.GetCart(ctx.Context(), int64(user.ID))
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.SuccessResponse(ctx, "cart items", items)
}

func (h *cartHandler) UpdateQuantity(ctx *fiber.Ctx) error {
	user, err := h.mid.GetCurrentUser(ctx)
	if err != nil {
		return response.UnauthorizedResponse(ctx, err.Error())
	}

	var req []dto.UpdateCartItem

	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	for _, i := range req {
		if err := h.validate.Struct(i); err != nil {
			return response.BadRequestResponse(ctx, err.Error())
		}
	}

	if err := h.uc.UpdateQuantity(ctx.Context(), user.ID, req); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.SuccessResponse(ctx, "quantity updated", nil)
}

func (h *cartHandler) DeleteItem(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	if err := h.uc.DeleteItem(ctx.Context(), int64(id)); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.SuccessResponse(ctx, "item deleted", nil)
}
