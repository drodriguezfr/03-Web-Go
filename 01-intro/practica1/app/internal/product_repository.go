package internal

import "errors"

var (
	// ErrProductCodeValueAlreadyInUse is the error when the product code value is already in use
	ErrProductCodeValueAlreadyInUse = errors.New("code_value is already in use")
)

// ProductRepository is the interface for the product repository
type ProductRepository interface {
	// Save saves movies in the repository

	Save(product *Product) (err error)
}
