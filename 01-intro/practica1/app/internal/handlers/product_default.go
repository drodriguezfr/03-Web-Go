package handlers

import (
	"app/internal"
	"encoding/json"
	"errors"
	"net/http"
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
		var body BodyRequestProductJSON
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
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
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid body"))
			case errors.Is(err, internal.ErrProductCodeValueExist):
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("code_value already exists"))
			case errors.Is(err, internal.ErrProductCodeValueAlreadyInUse):
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("code_value already exists"))

			default:
				w.WriteHeader(http.StatusInternalServerError)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// pasamos a bytes para enviarlo
		json.NewEncoder(w).Encode(map[string]any{
			"message": "product created successfully",
			"data":    data,
		})
	}
}
