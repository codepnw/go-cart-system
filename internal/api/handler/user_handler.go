package handler

import (
	"errors"

	"github.com/codepnw/go-cart-system/internal/api/response"
	"github.com/codepnw/go-cart-system/internal/dto"
	"github.com/codepnw/go-cart-system/internal/usecase"
	"github.com/codepnw/go-cart-system/internal/utils/errs"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	uc       usecase.UserUsecase
	validate *validator.Validate
}

func NewUserHandler(uc usecase.UserUsecase) *userHandler {
	return &userHandler{
		uc:       uc,
		validate: validator.New(),
	}
}

func (h *userHandler) Register(ctx *fiber.Ctx) error {
	var req dto.CreateUser

	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validate.Struct(req); err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	if err := h.uc.Register(ctx.Context(), &req); err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			return response.BadRequestResponse(ctx, err.Error())
		}
		return response.InternalServerError(ctx, err)
	}

	return response.CreastedResponse(ctx, "register successfully", req.Email)
}

func (h *userHandler) Login(ctx *fiber.Ctx) error {
	var req dto.UserCredential

	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validate.Struct(req); err != nil {
		return response.BadRequestResponse(ctx, err.Error())
	}

	data, err := h.uc.Login(ctx.Context(), &req)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			return response.BadRequestResponse(ctx, err.Error())
		}
		return response.InternalServerError(ctx, err)
	}

	return response.SuccessResponse(ctx, "", data)
}
