package repository

import "app/internal"

type ProductMap struct {
	// db is a map that represents the database
	// key: product ID
	// value: product

	db map[int]internal.Product

	// lastID is the last ID used for a product
	lastID int
}

// NewProductMap returns a new MovieMap instance
func NewMovieMap(db map[int]internal.Product, lastID int) *ProductMap {
	// default config / values
	// ...

	return &ProductMap{
		db:     db,
		lastID: lastID,
	}
}

func (p *ProductMap) Save(product *internal.Product) (err error) {
	// validate product
	// code_value is unique

	for _, p := range p.db {
		if p.Code_value == product.Code_value {

			return internal.ErrProductCodeValueAlreadyInUse
		}
	}
	// autoincrement id
	(*p).lastID++
	// set id
	(*product).ID = (*p).lastID

	// save product
	(*p).db[(*product).ID] = *product

	return
}
