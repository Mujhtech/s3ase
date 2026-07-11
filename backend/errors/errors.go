package errors

import "errors"

var (
	ErrNotAuthorized = errors.New("not authorized")
	ErrInvalidInput  = errors.New("invalid input")
	ErrConflict      = errors.New("conflict")
)
