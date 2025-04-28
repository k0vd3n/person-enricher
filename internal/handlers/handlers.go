package api

import (
	"encoding/json"
	"net/http"
	"person-enricher/internal/models"
	"person-enricher/internal/service"
	"strings"

	"github.com/gorilla/mux"
)

// Handler stores dependencies (e.g. service) and provides HTTP methods
type Handler struct {
	service service.PersonService
}

// NewHandler constructor, injects service
func NewHandler(s service.PersonService) *Handler {
	return &Handler{service: s}
}

// respondError helper for errors
func respondError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: msg})
}

// respondJSON helper for successful JSON responses
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// GetPeople GET /people
func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	// TODO: read filters/pagination, call h.service.GetPeople(...)
	w.WriteHeader(http.StatusNotImplemented)
}

// GetPersonByID GET /people/{id}
func (h *Handler) GetPersonByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if strings.TrimSpace(id) == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}

	person, err := h.service.GetPersonByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "person not found")
		return
	}

	respondJSON(w, http.StatusOK, person)
}

// CreatePerson POST /people
func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	// validate neccessary fields
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Surname) == "" {
		respondError(w, http.StatusBadRequest, "name and surname are required")
		return
	}

	person, err := h.service.CreatePerson(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not create person")
		return
	}

	respondJSON(w, http.StatusCreated, person)
}

// UpdatePerson PUT /people/{id}
func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var req models.UpdatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	// check id
	vars := mux.Vars(r)
	id := vars["id"]
	if strings.TrimSpace(id) == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}
	// validate neccessary fields
	if strings.TrimSpace(req.Name) == "" ||
		strings.TrimSpace(req.Surname) == "" ||
		req.Age <= 0 ||
		strings.TrimSpace(req.Gender) == "" ||
		strings.TrimSpace(req.Nationality) == "" {
		respondError(w, http.StatusBadRequest, "name, surname, age (>0), gender and nationality are required")
		return
	}

	updated, err := h.service.UpdatePerson(r.Context(), id, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not update person")
		return
	}

	respondJSON(w, http.StatusOK, updated)
}

// DeletePerson DELETE /people/{id}
func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if strings.TrimSpace(id) == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.service.DeletePerson(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "could not delete person")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

