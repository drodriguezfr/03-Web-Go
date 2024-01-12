package internal

import "errors"

var (
	// ErrProductFieldRequired is the error when a field is required
	ErrProductFieldRequired = errors.New("missing required field")

	// ErrProductQuality is the error when the product quality is invalid
	ErrProductQuality = errors.New("product quality")

	// ErrProductCodeValueExist is the error when the product code value is already in use
	ErrProductCodeValueExist = errors.New("code_value already exist")
)

// ProductService is the interface for the product service
// business logic
// validations
// external services
type ProductService interface {
	// Save saves a product
	Save(product *Product) (err error)

	// Get gets all products
	//Get() (products []Product, err error)

	// GetByID gets a product
	GetByID(id int) (product Product, err error)

	// Update updates a product
	Update(product *Product) (err error)

	// Delete deletes a product
	Delete(id int) (err error)
}
