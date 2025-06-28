package errs

import "errors"

var (
	ErrCartItemNotFound = errors.New("cart items not found")
)
