package handlers

import (
	"app/internal"
	"app/internal/platform/web/request"
	"app/internal/platform/web/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// NewDefaultProduct returns a new DefaultProduct
func NewDefaultProduct(sv internal.ProductService) *DefaultProducts {
	return &DefaultProducts{
		sv: sv,
	}
}

// DefaultProducts is a simple in-memory collection of products
type DefaultProducts struct {
	// sv is a product service
	sv internal.ProductService
}

// BodyRequestProductJSON is the body request for create a new product
type BodyRequestProductJSON struct {
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
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

// Creates a new product
func (d *DefaultProducts) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// token
		//token := r.Header.Get("Authorization")
		//if token != "123456" {
		//	w.Header().Set("Content-Type", "text/plain")
		//	w.WriteHeader(http.StatusUnauthorized)
		//	w.Write([]byte("invalid token"))
		//	return
		//}

		// request

		// validate token
		if ValidateToken(r) != nil {
			response.Text(w, http.StatusUnauthorized, "invalid token")
			return
		}

		var body BodyRequestProductJSON
		if err := request.JSON(r, &body); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")

			return
		}

		// process

		// serialize product.Product
		product := internal.Product{
			Name:         body.Name,
			Quantity:     body.Quantity,
			Code_value:   body.Code_value,
			Is_published: body.Is_published,
			Expiration:   body.Expiration,
			Price:        body.Price,
		}

		// save product
		if err := d.sv.Save(&product); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductFieldRequired), errors.Is(err, internal.ErrProductQuality):
				response.Text(w, http.StatusBadRequest, "invalid body")
			case errors.Is(err, internal.ErrProductCodeValueExist):
				response.Text(w, http.StatusConflict, "code_value already exists")
			case errors.Is(err, internal.ErrProductCodeValueAlreadyInUse):
				response.Text(w, http.StatusConflict, "code_value already in use")

			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// serialize product.Product to ProductJSON
		data := ProductJSON{
			ID:           product.ID,
			Name:         product.Name,
			Quantity:     product.Quantity,
			Code_value:   product.Code_value,
			Is_published: product.Is_published,
			Expiration:   product.Expiration,
			Price:        product.Price,
		}
		response.JSON(w, http.StatusCreated, map[string]any{"message": "product created", "data": data})
	}
}

// GetByID returns a product by id
func (d *DefaultProducts) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// validate token
		if ValidateToken(r) != nil {
			response.Text(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// get id from url
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		// get product

		product, err := d.sv.GetByID(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product not found")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// serialize ProductJSON
		data := ProductJSON{
			ID:           product.ID,
			Name:         product.Name,
			Quantity:     product.Quantity,
			Code_value:   product.Code_value,
			Is_published: product.Is_published,
			Expiration:   product.Expiration,
			Price:        product.Price,
		}
		response.JSON(w, http.StatusOK, map[string]any{"message": "product found", "data": data})
	}

}

// Update updates a product
func (d *DefaultProducts) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// request

		// validate token
		if ValidateToken(r) != nil {
			response.Text(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// get id from url
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// get product to []byte (porque solo se puede leer una vez el body)
		// entonces se guarda en un []byte para poder leerlo varias veces
		// en este caso una vez para validar el body y otra para serializarlo
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		// body to map[string]any
		var bodyMap map[string]any
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		// validate body
		if err := ValidateExists(bodyMap, "name", "quantity", "code_value", "is_published", "expiration", "price"); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		// get body
		var body BodyRequestProductJSON
		if err = json.Unmarshal(bytes, &body); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		// process
		// serialize internal.Product
		product := internal.Product{
			ID:           id,
			Name:         body.Name,
			Quantity:     body.Quantity,
			Code_value:   body.Code_value,
			Is_published: body.Is_published,
			Expiration:   body.Expiration,
			Price:        body.Price,
		}

		// update product
		if err := d.sv.Update(&product); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductFieldRequired), errors.Is(err, internal.ErrProductQuality):
				response.Text(w, http.StatusBadRequest, "invalid body")
			case errors.Is(err, internal.ErrProductCodeValueExist):
				response.Text(w, http.StatusConflict, "code_value already exists")
			case errors.Is(err, internal.ErrProductCodeValueAlreadyInUse):
				response.Text(w, http.StatusConflict, "code_value already in use")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// deserialize data
		data := ProductJSON{
			ID:           product.ID,
			Name:         product.Name,
			Quantity:     product.Quantity,
			Code_value:   product.Code_value,
			Is_published: product.Is_published,
			Expiration:   product.Expiration,
			Price:        product.Price,
		}
		response.JSON(w, http.StatusOK, map[string]any{"message": "product updated", "data": data})
	}
}

// UpdatePartial updates a product
func (d *DefaultProducts) UpdatePartial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// validate token
		if ValidateToken(r) != nil {
			response.Text(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// get id from url
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// get product from service
		product, err := d.sv.GetByID(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product not found")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}
		// process
		// serialize internal.Product to BodyRequestProductJSON
		reqBody := BodyRequestProductJSON{
			Name:         product.Name,
			Quantity:     product.Quantity,
			Code_value:   product.Code_value,
			Is_published: product.Is_published,
			Expiration:   product.Expiration,
			Price:        product.Price,
		}

		// get body
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		// serialize internal.Product
		product = internal.Product{
			ID:           id,
			Name:         reqBody.Name,
			Quantity:     reqBody.Quantity,
			Code_value:   reqBody.Code_value,
			Is_published: reqBody.Is_published,
			Expiration:   reqBody.Expiration,
			Price:        reqBody.Price,
		}

		// update product
		if err := d.sv.Update(&product); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductFieldRequired), errors.Is(err, internal.ErrProductQuality):
				response.Text(w, http.StatusBadRequest, "invalid body")
			case errors.Is(err, internal.ErrProductCodeValueExist):
				response.Text(w, http.StatusConflict, "code_value already exists")
			case errors.Is(err, internal.ErrProductCodeValueAlreadyInUse):
				response.Text(w, http.StatusConflict, "code_value already in use")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// deserialize data
		// deserialize data
		data := ProductJSON{
			ID:           product.ID,
			Name:         product.Name,
			Quantity:     product.Quantity,
			Code_value:   product.Code_value,
			Is_published: product.Is_published,
			Expiration:   product.Expiration,
			Price:        product.Price,
		}
		response.JSON(w, http.StatusOK, map[string]any{"message": "product updated", "data": data})
	}
}

// Delete deletes a product
func (d *DefaultProducts) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// validate token
		if ValidateToken(r) != nil {
			response.Text(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// get id from url
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		// delete product
		if err := d.sv.Delete(id); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product not found")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		response.JSON(w, http.StatusOK, map[string]any{"message": "product deleted"})
	}
}

func ValidateExists(mp map[string]any, keys ...string) (err error) {
	for _, key := range keys {
		if _, ok := mp[key]; !ok {
			return fmt.Errorf("key %s not found", key)
		}
	}
	return
}
func ValidateToken(r *http.Request) (err error) {

	// get token from header
	token := r.Header.Get("Token")
	if token != os.Getenv("TOKEN") {
		return fmt.Errorf("invalid token")
	}
	return
}
