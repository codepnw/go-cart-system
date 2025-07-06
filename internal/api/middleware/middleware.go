package middleware

import (
	"errors"
	"strings"

	"github.com/codepnw/go-cart-system/config"
	"github.com/codepnw/go-cart-system/internal/api/response"
	"github.com/codepnw/go-cart-system/internal/domain"
	"github.com/codepnw/go-cart-system/internal/utils/security"
	"github.com/gofiber/fiber/v2"
)

const userContextKey = "user"

type Middleware struct {
	cfg   *config.EnvConfig
	token *security.TokenConfig
}

func NewMiddleware(cfg *config.EnvConfig) *Middleware {
	return &Middleware{
		cfg:   cfg,
		token: security.NewTokenConfig(*cfg),
	}
}

func (m *Middleware) Authorize() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")

		if authHeader == "" {
			return response.UnauthorizedResponse(ctx, "auth header is missing")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.UnauthorizedResponse(ctx, "invalid token format")
		}

		user, err := m.token.VerifyToken(parts[1], m.cfg.JWTSecretKey)
		if err != nil || user.ID < 0 {
			return response.UnauthorizedResponse(ctx, err.Error())
		}

		ctx.Locals(userContextKey, user)
		return ctx.Next()
	}
}

func (m *Middleware) GetCurrentUser(ctx *fiber.Ctx) (*domain.User, error) {
	user, ok := ctx.Locals(userContextKey).(*domain.User)
	if !ok || user == nil {
		return nil, errors.New("unauthorized")
	}
	return user, nil
}
