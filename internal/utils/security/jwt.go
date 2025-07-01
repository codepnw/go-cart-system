package security

import (
	"fmt"
	"time"

	"github.com/codepnw/go-cart-system/config"
	"github.com/codepnw/go-cart-system/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type tokenConfig struct {
	id         int64
	email      string
	role       string
	secretKey  string
	refreshKey string
}

func NewTokenConfig(config config.EnvConfig, user *domain.User) *tokenConfig {
	return &tokenConfig{
		id:         user.ID,
		email:      user.Email,
		role:       user.Role,
		secretKey:  config.JWTSecretKey,
		refreshKey: config.JWTRefreshKey,
	}
}

func (t *tokenConfig) GenerateAccessToken() (string, error) {
	duration := time.Hour * 24
	return t.generateToken(t.id, t.email, t.role, t.secretKey, duration)
}

func (t *tokenConfig) GenerateRefreshToken() (string, error) {
	duration := time.Hour * 24 * 7
	return t.generateToken(t.id, t.email, t.role, t.refreshKey, duration)
}

func (t *tokenConfig) generateToken(id int64, email, role, key string, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(exp).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("signed token failed: %w", err)
	}

	return tokenStr, nil
}
