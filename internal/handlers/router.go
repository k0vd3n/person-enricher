package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
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

	r.HandleFunc("/people", h.GetPeople).Methods(http.MethodGet)
	r.HandleFunc("/people/{id}", h.GetPersonByID).Methods(http.MethodGet)
	r.HandleFunc("/people", h.CreatePerson).Methods(http.MethodPost)
	r.HandleFunc("/people/{id}", h.UpdatePerson).Methods(http.MethodPut)
	r.HandleFunc("/people/{id}", h.DeletePerson).Methods(http.MethodDelete)

	return r
}
