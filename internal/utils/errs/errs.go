package errs

import "errors"

var (
	ErrCartItemNotFound = errors.New("cart items not found")
	ErrProductNotFound  = errors.New("product not found")

	// User
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user email already exists")
)
