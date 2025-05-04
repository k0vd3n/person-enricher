package handlers

import (
	"net/http"

	_ "person-enricher/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter returns a new Gorilla Mux router with endpoints:
//
// * GET /people: GetPeople
// * GET /people/{id}: GetPersonByID
// * POST /people: CreatePerson
// * PUT /people/{id}: UpdatePerson
// * DELETE /people/{id}: DeletePerson
func NewRouter(h *Handler) *mux.Router {
	r := mux.NewRouter()

	// Swagger
	r.PathPrefix("/v1/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/v1/people", h.GetPeople).Methods(http.MethodGet)
	r.HandleFunc("/v1/people/{id}", h.GetPersonByID).Methods(http.MethodGet)
	r.HandleFunc("/v1/people", h.CreatePerson).Methods(http.MethodPost)
	r.HandleFunc("/v1/people/{id}", h.UpdatePerson).Methods(http.MethodPut)
	r.HandleFunc("/v1/people/{id}", h.DeletePerson).Methods(http.MethodDelete)

	return r
}
