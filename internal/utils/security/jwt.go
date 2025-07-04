package security

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/codepnw/go-cart-system/config"
	"github.com/codepnw/go-cart-system/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	secretKey  string
	refreshKey string
}

func NewTokenConfig(config config.EnvConfig) *TokenConfig {
	return &TokenConfig{
		secretKey:  config.JWTSecretKey,
		refreshKey: config.JWTRefreshKey,
	}
}

func (t *TokenConfig) GenerateAccessToken(user *domain.User) (string, error) {
	duration := time.Hour * 24
	return t.generateToken(user.ID, user.Email, user.Role, t.secretKey, duration)
}

func (t *TokenConfig) GenerateRefreshToken(user *domain.User) (string, error) {
	duration := time.Hour * 24 * 7
	return t.generateToken(user.ID, user.Email, user.Role, t.refreshKey, duration)
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

func (t *TokenConfig) VerifyToken(tokenStr, jwtKey string) (*domain.User, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println(t.Header)
			return nil, errors.New("unknow signing method")
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, errors.New("token is expired")
		}

		user := domain.User{}
		user.ID = int64(claims["user_id"].(float64))
		user.Email = claims["email"].(string)
		user.Role = claims["role"].(string)

		return &user, nil
	}

	return nil, errors.New("token verification failed")
}
