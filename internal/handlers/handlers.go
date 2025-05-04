package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"person-enricher/internal/models"
	"person-enricher/internal/service"
	"strconv"
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
	log.Printf("handlers.respondError: %s", msg)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: msg})
}

func toPersonResponse(p models.Person) models.PersonResponse {
	return models.PersonResponse{
		ID:          p.ID,
		Name:        p.Name,
		Surname:     p.Surname,
		Patronymic:  p.Patronymic,
		Age:         p.Age,
		Gender:      p.Gender,
		Nationality: p.Nationality,
	}
}

// respondJSON writes a JSON response with the given status code and payload.
// It sets the Content-Type header to application/json, the status code to the given status,
// and encodes the payload in the request body as a JSON object.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	log.Printf("handlers.respondJSON: %v", payload)
	json.NewEncoder(w).Encode(payload)
}

// GetPeople responds to GET /people requests.
// It reads the page and size query parameters and calls h.service.GetPeople with the given filter.
// It returns the list of people as a JSON object with status code 200.
// If the page or size is invalid, it returns a 400 error.
// If the people could not be fetched, it returns a 500 error.
func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// filter
	log.Printf("handlers.GetPeople: getting people with filter: %s", q.Get("filter"))
	filterStr := strings.TrimSpace(q.Get("filter"))

	// page
	page := 1
	if p := q.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		} else {
			log.Printf("handlers.GetPeople: invalid page parameter: %s", p)
			respondError(w, http.StatusBadRequest, "invalid page parameter")
			return
		}
	}
	log.Printf("handlers.GetPeople: page: %d", page)

	// size
	size := 10
	if s := q.Get("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			size = v
		} else {
			log.Printf("handlers.GetPeople: invalid size parameter: %s", s)
			respondError(w, http.StatusBadRequest, "invalid size parameter")
			return
		}
	}
	log.Printf("handlers.GetPeople: size: %d", size)

	// filter
	pf := models.PeopleFilter{
		Filter: filterStr,
		Page:   page,
		Size:   size,
	}
	log.Printf("handlers.GetPeople: result people filter: %v", pf)

	// get people
	log.Printf("handlers.GetPeople: getting people")
	log.Printf("handlers.GetPeople: send request to service")
	people, err := h.service.GetPeople(r.Context(), pf)
	if err != nil {
		log.Printf("handlers.GetPeople: could not get people: %v", err)
		respondError(w, http.StatusInternalServerError, "could not fetch people")
		return
	}


	log.Printf("handlers.GetPeople: got %d people", len(people))
	respondJSON(w, http.StatusOK, people)
}

// GetPersonByID responds to GET /people/{id} requests.
// It reads the id path parameter, calls h.service.GetPersonByID with the given id,
// and writes the response as a JSON object with status code 200.
// If the id is empty, it returns a 400 error.
// If the person is not found, it returns a 404 error.
func (h *Handler) GetPersonByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("handlers.GetPersonByID: getting person with id: %s", id)
	log.Printf("handlers.GetPersonByID: send request to service")
	person, err := h.service.GetPersonByID(r.Context(), id)
	if err != nil {
		log.Printf("handlers.GetPersonByID: could not get person: %v", err)
		respondError(w, http.StatusInternalServerError, "internal service error")
		return
	}
	log.Printf("handlers.GetPersonByID: got person: %v", person)
	if person == (models.Person{}) {
		log.Printf("handlers.GetPersonByID: person not found")
		respondError(w, http.StatusNotFound, "person not found")
		return
	}

	respondJSON(w, http.StatusOK, toPersonResponse(person))
}

// CreatePerson responds to POST /people requests.
// It reads the JSON body, validates the name and surname fields,
// calls h.service.CreatePerson with the given request, and writes the response
// as a JSON object with status code 201.
// If the body is invalid JSON, it returns a 400 error.
// If the name or surname is empty, it returns a 400 error.
// If the person could not be created, it returns a 500 error.
func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlers.CreatePerson: creating person")
	var req models.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("handlers.CreatePerson: invalid JSON: %v", err)
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	// validate neccessary fields
	log.Printf("handlers.CreatePerson: validating neccessary fields")
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Surname) == "" {
		respondError(w, http.StatusBadRequest, "name and surname are required")
		return
	}
	log.Printf("handlers.CreatePerson: name: %s, surname: %s", req.Name, req.Surname)

	log.Printf("handlers.CreatePerson: send request to service")
	person, err := h.service.CreatePerson(r.Context(), req)
	if err != nil {
		log.Printf("handlers.CreatePerson: could not create person: %v", err)
		respondError(w, http.StatusInternalServerError, "could not create person")
		return
	}

	log.Printf("handlers.CreatePerson: person created: %v", person)
	respondJSON(w, http.StatusCreated, toPersonResponse(person))
}

// UpdatePerson responds to PUT /people/{id} requests.
// It reads the JSON body, validates the name, surname, age, gender and nationality fields,
// calls h.service.UpdatePerson with the given request, and writes the response
// as a JSON object with status code 200.
// If the body is invalid JSON, it returns a 400 error.
// If the name, surname, age, gender or nationality is empty, it returns a 400 error.
// If the person could not be updated, it returns a 500 error.
func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlers.UpdatePerson: updating person")
	var req models.UpdatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("handlers.UpdatePerson: invalid JSON: %v", err)
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	
	// check id
	log.Printf("handlers.UpdatePerson: checking id")
	vars := mux.Vars(r)
	id := vars["id"]
	if strings.TrimSpace(id) == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}

	// validate neccessary fields
	log.Printf("handlers.UpdatePerson: validating neccessary fields")
	if strings.TrimSpace(req.Name) == "" ||
		strings.TrimSpace(req.Surname) == "" ||
		req.Age <= 0 ||
		strings.TrimSpace(req.Gender) == "" ||
		strings.TrimSpace(req.Nationality) == "" {
		respondError(w, http.StatusBadRequest, "name, surname, age (>0), gender and nationality are required")
		return
	}
	log.Printf("handlers.UpdatePerson: name: %s, surname: %s, age: %d, gender: %s, nationality: %s", req.Name, req.Surname, req.Age, req.Gender, req.Nationality)

	log.Printf("handlers.UpdatePerson: send update request to service")
	updated, err := h.service.UpdatePerson(r.Context(), id, req)
	if err != nil {
		log.Printf("handlers.UpdatePerson: could not update person: %v", err)
		respondError(w, http.StatusInternalServerError, "could not update person")
		return
	}
	log.Printf("handlers.UpdatePerson: person updated: %v", updated)

	respondJSON(w, http.StatusOK, toPersonResponse(updated))
}

// DeletePerson responds to DELETE /people/{id} requests.
// It reads the id path parameter, calls h.service.DeletePerson with the given id,
// and writes the response as a JSON object with status code 204.
// If the id is empty, it returns a 400 error.
// If the person could not be deleted, it returns a 500 error.
func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	log.Printf("handlers.DeletePerson: deleting person")
	vars := mux.Vars(r)
	id := vars["id"]
	if strings.TrimSpace(id) == "" {
		log.Printf("handlers.DeletePerson: id is required")
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}

	log.Printf("handlers.DeletePerson: send delete request to service")
	if err := h.service.DeletePerson(r.Context(), id); err != nil {
		log.Printf("handlers.DeletePerson: could not delete person: %v", err)
		respondError(w, http.StatusInternalServerError, "could not delete person")
		return
	}

	log.Printf("handlers.DeletePerson: person deleted")
	// вернуть пользователю информацию, что запись успешно удалена
	respondJSON(w, http.StatusOK, "the record was successfully deleted")
}
