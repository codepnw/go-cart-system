package security

import (
	"fmt"

	"github.com/codepnw/go-cart-system/internal/utils/errs"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate hash password failed: %w", err)
	}
	return string(hashed), nil
}

func CheckHashPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errs.ErrInvalidCredentials
	}
	return nil
}
