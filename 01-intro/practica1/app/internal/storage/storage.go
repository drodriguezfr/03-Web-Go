package storage

import (
	"app/internal"
	"errors"
)

var (
	// ErrStorageTimeLayout is the error when the time layout is invalid
	ErrStorageTimeLayout = errors.New("time layout is invalid")
)

type Storage interface {

	// ReadAll reads a json file and returns a map of products
	ReadAll() (p []internal.Product, err error)

	// WriteAll writes a map of products to a json file
	WriteAll(p []internal.Product) (err error)
}
