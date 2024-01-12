package service

import (
	"app/internal"
	"fmt"
	"time"
)

// NewDefaultProduct creates a new instance of ProductDefault
func NewDefaultProduct(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		rp: rp,
	}
}

// ProductDefault is a struct that represents the product default service
type ProductDefault struct {

	// rp is a product repository
	rp internal.ProductRepository

	// external services
	// ... (weather api, etc.)
}

// Save saves a product
func (p *ProductDefault) Save(product *internal.Product) (err error) {
	// external services
	// ...

	//business logic
	// - validate required fields
	if product.Name == "" {
		return fmt.Errorf("%w: name", internal.ErrProductFieldRequired)
	}

	if product.Code_value == "" {
		return fmt.Errorf("%w: product_code", internal.ErrProductFieldRequired)
	}
	if product.Price == 0 {
		return fmt.Errorf("%w: price", internal.ErrProductFieldRequired)
	}
	if product.Quantity == 0 {
		return fmt.Errorf("%w: quantity", internal.ErrProductFieldRequired)
	}
	if product.Expiration == "" {
		return fmt.Errorf("%w: expiration", internal.ErrProductFieldRequired)
	}

	// - validate quality
	dateFormat := "02/01/2006"
	_, error := time.Parse(dateFormat, product.Expiration)

	if error != nil {
		return fmt.Errorf("%w: expiration", internal.ErrProductQuality)
	}

	// save product
	err = p.rp.Save(product)
	if err != nil {
		switch err {
		case internal.ErrProductCodeValueExist:
			return fmt.Errorf("%w: code_value", internal.ErrProductCodeValueExist)
		}
		return
	}
	return
}

// GetByID gets a product
func (p *ProductDefault) GetByID(id int) (product internal.Product, err error) {
	// get product
	product, err = p.rp.GetByID(id)
	if err != nil {
		switch err {
		case internal.ErrProductNotFound:
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
		return
	}
	return
}

// Update updates a product
func (p *ProductDefault) Update(product *internal.Product) (err error) {
	// update product
	err = p.rp.Update(product)
	if err != nil {
		switch err {
		case internal.ErrProductNotFound:
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
		return
	}
	return
}

// Delete deletes a product
func (p *ProductDefault) Delete(id int) (err error) {
	// delete product
	err = p.rp.Delete(id)
	if err != nil {
		switch err {
		case internal.ErrProductNotFound:
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
		return
	}
	return
}
