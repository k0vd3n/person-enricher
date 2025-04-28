package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter настраивает все HTTP-эндпоинты и возвращает готовый mux.Router
func NewRouter(h *Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/people", h.GetPeople).Methods(http.MethodGet)
	r.HandleFunc("/people/{id}", h.GetPersonByID).Methods(http.MethodGet)
	r.HandleFunc("/people", h.CreatePerson).Methods(http.MethodPost)
	r.HandleFunc("/people/{id}", h.UpdatePerson).Methods(http.MethodPut)
	r.HandleFunc("/people/{id}", h.DeletePerson).Methods(http.MethodDelete)

	return r
}
