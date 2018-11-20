package errors

import "errors"

var (
	ErrInvalidID = errors.New("Invalid ID")
	ErrNil       = errors.New("Null variable")
)
