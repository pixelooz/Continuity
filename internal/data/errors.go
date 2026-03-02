package data

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)
