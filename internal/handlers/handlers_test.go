package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetPeople(t *testing.T) {
	router, _ := setupTest()

	tests := []struct {
		name       string
		query        string
		statusCode int
	}{
		{"valid request", "/people", http.StatusOK},
		{"custom page and size", "/people?page=2&size=20", http.StatusOK},
		{"invalid page", "/people?page=abc", http.StatusBadRequest},
		{"invalid size", "/people?size=abc", http.StatusBadRequest},
		{"error case", "/people?filter=error", http.StatusInternalServerError},
		{"empty filter string", "/people?filter=", http.StatusOK},
		{"custom filter", "/people?filter=Test", http.StatusOK},
		{"custom page, size, filter", "/people?page=2&size=20&filter=User", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.query, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}

func TestGetPersonByID(t *testing.T) {
	router, _ := setupTest()

	tests := []struct {
		name       string
		id         string
		statusCode int
	}{
		{"valid request", "valid-id", http.StatusOK},
		{"not found", "notfound-id", http.StatusNotFound},
		{"service error", "error-id", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/people/%s", tt.id)
			req, _ := http.NewRequest("GET", url, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}

func TestCreatePerson(t *testing.T) {
	router, _ := setupTest()

	tests := []struct {
		name       string
		body       string
		statusCode int
	}{
		{"valid request", `{"name":"John","surname":"Doe"}`, http.StatusCreated},
		{"invalid json", `{invalid}`, http.StatusBadRequest},
		{"invalid request with empty name", `{"name":"","surname":"Doe"}`, http.StatusBadRequest},
		{"invalid request with empty surname", `{"name":"John","surname":""}`, http.StatusBadRequest},
		{"invalid request with invalid JSON body", `{"name":"John"}`, http.StatusBadRequest},
		{"service error", `{"name":"error","surname":"Doe"}`, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/people", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}

func TestUpdatePerson(t *testing.T) {
	router, _ := setupTest()

	tests := []struct {
		name       string
		id         string
		body       string
		statusCode int
	}{
		{"valid request", "valid-id", `{"name":"John","surname":"Doe","age":30,"gender":"male","nationality":"US"}`, http.StatusOK},
		{"invalid id", " ", `{"name":"John","surname":"Doe","age":30,"gender":"male","nationality":"US"}`, http.StatusBadRequest},
		{"service error", "error-id", `{"name":"John","surname":"Doe","age":30,"gender":"male","nationality":"US"}`, http.StatusInternalServerError},
		{"invalid request with empty name", "1", `{"name":"","surname":"Doe"}`, http.StatusBadRequest},
		{"invalid request with empty surname", "1", `{"name":"John","surname":""}`, http.StatusBadRequest},
		{"invalid request with empty age", "1", `{"name":"John","surname":"Doe","age":0}`, http.StatusBadRequest},
		{"invalid request with empty gender", "1", `{"name":"John","surname":"Doe","gender":""}`, http.StatusBadRequest},
		{"invalid request with empty nationality", "1", `{"name":"John","surname":"Doe","nationality":""}`, http.StatusBadRequest},
		{"invalid request with invalid JSON body", "1", `{"name":"John"}`, http.StatusBadRequest},
		{"invalid json", "1", `{invalid}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/people/%s", tt.id)
			req, _ := http.NewRequest("PUT", url, bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}

func TestDeletePerson(t *testing.T) {
	router, _ := setupTest()

	tests := []struct {
		name       string
		id         string
		statusCode int
	}{
		{"valid request", "valid-id", http.StatusOK},
		{"empty id", " ", http.StatusBadRequest},
		{"service error", "error-id", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/people/%s", tt.id)
			req, _ := http.NewRequest("DELETE", url, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}

// Вспомогательная функция для создания роутера и обработки запросов
func setupTest() (*mux.Router, *MockPersonService) {
	service := &MockPersonService{}
	handler := NewHandler(service)
	router := NewRouter(handler)
	return router, service
}
