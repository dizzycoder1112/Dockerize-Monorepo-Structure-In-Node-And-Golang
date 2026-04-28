package order

import "errors"

var (
	ErrNotFound          = errors.New("order not found")
	ErrInvalidInput      = errors.New("invalid order input")
	ErrInvalidAmount     = errors.New("invalid order amount")
	ErrInvalidTransition = errors.New("invalid order status transition")
)
