package domain

import "errors"

var (
	ErrNotFound   = errors.New("url not found")
	ErrConflict   = errors.New("alias already exists")
	ErrValidation = errors.New("validation failed")
)
