package handler

import (
	"github.com/codepnw/go-cart-system/internal/api/middleware"
	"github.com/codepnw/go-cart-system/internal/api/response"
	"github.com/codepnw/go-cart-system/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type orderHandler struct {
	mid *middleware.Middleware
	uc  usecase.OrderUsecase
}

func NewOrderHandler(mid *middleware.Middleware, orderUc usecase.OrderUsecase) *orderHandler {
	return &orderHandler{
		mid,
		orderUc,
	}
}

func (h *orderHandler) Checkout(ctx *fiber.Ctx) error {
	user, err := h.mid.GetCurrentUser(ctx)
	if err != nil {
		return response.UnauthorizedResponse(ctx, err.Error())
	}

	err = h.uc.Checkout(ctx.Context(), user.ID, "")
	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.SuccessResponse(ctx, "checkout completed", nil)
}
