package internal

import "errors"

var (
	// ErrProductCodeValueAlreadyInUse is the error when the product code value is already in use
	ErrProductCodeValueAlreadyInUse = errors.New("code_value is already in use")

	// ErrProductNotFound is the error when the product is not found
	ErrProductNotFound = errors.New("product not found")
)

// ProductRepository is the interface for the product repository
type ProductRepository interface {
	// Save saves movies in the repository

	Save(product *Product) (err error)

	// GetByID gets a product by ID
	GetByID(id int) (product Product, err error)

	// Get gets all products
	//Get() (products []Product, err error)

	// Update updates a product
	Update(product *Product) (err error)

	// Delete deletes a product
	Delete(id int) (err error)
}
