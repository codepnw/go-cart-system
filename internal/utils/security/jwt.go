package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	ID    int64
	Email string
	Role  string

	secretKey  string
	refreshKey string
}

func SetupJWT(secretKey, refreshKey string) *TokenConfig {
	return &TokenConfig{
		secretKey:  secretKey,
		refreshKey: refreshKey,
	}
}

func (t *TokenConfig) GenerateAccessToken() (string, error) {
	duration := time.Hour * 24
	return t.generateToken(t.ID, t.Email, t.Role, t.secretKey, duration)
}

func (t *TokenConfig) GenerateRefreshToken() (string, error) {
	duration := time.Hour * 24 * 7
	return t.generateToken(t.ID, t.Email, t.Role, t.refreshKey, duration)
}

func (t *TokenConfig) generateToken(id int64, email, role, key string, exp time.Duration) (string, error) {
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
