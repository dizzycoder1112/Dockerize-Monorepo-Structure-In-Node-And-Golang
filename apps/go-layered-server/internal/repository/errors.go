package repository

import "errors"

var (
	ErrNotFound    = errors.New("record not found")
	ErrInvalidArgs = errors.New("invalid repository arguments")
)
