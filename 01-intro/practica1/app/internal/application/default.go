package application

import (
	"app/internal"
	"app/internal/handlers"
	"app/internal/repository"
	"app/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewDefaultHTTP(address string) *DefaultHTTP {
	//defaultAddress := "localhost:8080"
	//if address != "" {
	//	defaultAddress = address
	//}
	return &DefaultHTTP{
		address: address,
	}
}

type DefaultHTTP struct {
	address string
}

// Run the server
func (h *DefaultHTTP) Run() (err error) {

	// initialize dependencies

	// - repository
	rp := repository.NewProductMap(make(map[int]internal.Product), 0)

	// - service
	sv := service.NewDefaultProduct(rp)
	// - handler

	hd := handlers.NewDefaultProduct(sv)

	rp.LoadJson()

	rp.SaveJson()

	// - router
	r := chi.NewRouter()

	// endpoints
	r.Post("/products", hd.Create())

	r.Get("/products/{id}", hd.GetByID())

	r.Put("/products/{id}", hd.Update())

	r.Patch("/products/{id}", hd.UpdatePartial())

	r.Delete("/products/{id}", hd.Delete())

	err = http.ListenAndServe(h.address, r)
	return
}
