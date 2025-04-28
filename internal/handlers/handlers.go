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

// NewHandler returns a new Handler with the given PersonService
func NewHandler(s service.PersonService) *Handler {
	return &Handler{service: s}
}

// respondError writes an error response in JSON format.
// It sets the Content-Type header to application/json,
// the status code to the given status, and encodes the error message
// in the request body as a models.ErrorResponse.
func respondError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: msg})
}

// respondJSON writes a JSON response with the given status code and payload.
// It sets the Content-Type header to application/json, the status code to the given status,
// and encodes the payload in the request body as a JSON object.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	// TODO: read filters/pagination, call h.service.GetPeople(...)
	w.WriteHeader(http.StatusNotImplemented)
}

// GetPersonByID responds to GET /people/{id} requests.
// It reads the id path parameter, calls h.service.GetPersonByID with the given id,
// and writes the response as a JSON object with status code 200.
// If the id is empty, it returns a 400 error.
// If the person is not found, it returns a 404 error.
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

// CreatePerson responds to POST /people requests.
// It reads the JSON body, validates the name and surname fields,
// calls h.service.CreatePerson with the given request, and writes the response
// as a JSON object with status code 201.
// If the body is invalid JSON, it returns a 400 error.
// If the name or surname is empty, it returns a 400 error.
// If the person could not be created, it returns a 500 error.
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

// UpdatePerson responds to PUT /people/{id} requests.
// It reads the JSON body, validates the name, surname, age, gender and nationality fields,
// calls h.service.UpdatePerson with the given request, and writes the response
// as a JSON object with status code 200.
// If the body is invalid JSON, it returns a 400 error.
// If the name, surname, age, gender or nationality is empty, it returns a 400 error.
// If the person could not be updated, it returns a 500 error.
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

// DeletePerson responds to DELETE /people/{id} requests.
// It reads the id path parameter, calls h.service.DeletePerson with the given id,
// and writes the response as a JSON object with status code 204.
// If the id is empty, it returns a 400 error.
// If the person could not be deleted, it returns a 500 error.
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

