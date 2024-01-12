package repository

import (
	"app/internal"
	"encoding/json"
	"fmt"
	"os"
)

type ProductMap struct {
	// db is a map that represents the database
	// key: product ID
	// value: product

	db map[int]internal.Product

	// lastID is the last ID used for a product
	lastID int
}

// ProductJSON is the response for create a new product
type ProductJSON struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

// NewProductMap returns a new MovieMap instance
func NewProductMap(db map[int]internal.Product, lastID int) *ProductMap {
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

// GetByID returns a product by id
func (p *ProductMap) GetByID(id int) (product internal.Product, err error) {
	product, ok := p.db[id]
	if !ok {
		err = internal.ErrProductNotFound
		return
	}
	return
}

// Update updates a product
func (p *ProductMap) Update(product *internal.Product) (err error) {
	// validate existance
	_, ok := p.db[(*product).ID]
	if !ok {
		err = internal.ErrProductNotFound
		return
	}
	p.db[(*product).ID] = *product
	return
}

// Delete deletes a product
func (p *ProductMap) Delete(id int) (err error) {
	// validate existance
	_, ok := p.db[id]
	if !ok {
		err = internal.ErrProductNotFound
		return
	}
	delete(p.db, id)
	return
}

// LoadJson loads the products from the json file
func (p *ProductMap) LoadJson() {
	// slice auxiliar
	var products []internal.Product
	// abro el json
	file, err := os.Open("./products.json")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	// paso el json al slice
	readerProducts := json.NewDecoder(file)
	err = readerProducts.Decode(&products)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	//instancio el mapa de mi repository con el tamaño del slice auxiliar
	p.db = make(map[int]internal.Product, len(products))

	//recorro el slice auxiliar y almaceno los datos en el mapa
	for _, product := range products {
		p.db[product.ID] = product
	}
}

// SaveJson saves a product in the repository
func (p *ProductMap) SaveJson() {

	// slice auxiliar
	var products []internal.Product
	// recorro el mapa y almaceno los datos en el slice auxiliar
	for _, product := range p.db {
		products = append(products, product)
	}

	// creamos el archivo json
	file, err := os.Create("data.json")
	if err != nil {
		fmt.Println("Error al crear el archivo JSON:", err)
		return
	}
	defer file.Close()

	// Crear un encoder para el archivo
	encoder := json.NewEncoder(file)

	// Codificar el slice y escribirlo en el archivo
	if err := encoder.Encode(products); err != nil {
		fmt.Println("Error al codificar y escribir el JSON:", err)
		return
	}

	println(err)
}
