package data

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")

	ErrCollectionNameEmpty = errors.New("collection name is empty")
	ErrDuplicateCollection = errors.New("duplicate collection")
	ErrCollectionNotEmpty  = errors.New("collection is not empty")
)
