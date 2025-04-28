package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetPeople_NotImplemented(t *testing.T) {
	h := NewHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/people", nil)
	rr := httptest.NewRecorder()

	h.GetPeople(rr, req)

	if rr.Code != http.StatusNotImplemented {
		t.Fatalf("GET /people: expected %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}

func TestGetPersonByID_NotImplemented(t *testing.T) {
	h := NewHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/people/123", nil)
	rr := httptest.NewRecorder()

	h.GetPersonByID(rr, req)

	if rr.Code != http.StatusNotImplemented {
		t.Fatalf("GET /people/{id}: expected %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}

func TestCreatePerson_NotImplemented(t *testing.T) {
	h := NewHandler(nil)
	body := `{"name":"Dmitriy","surname":"Ushakov"}`
	req := httptest.NewRequest(http.MethodPost, "/people", strings.NewReader(body))
	rr := httptest.NewRecorder()

	h.CreatePerson(rr, req)

	if rr.Code != http.StatusNotImplemented {
		t.Fatalf("POST /people: expected %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}

func TestUpdatePerson_NotImplemented(t *testing.T) {
	h := NewHandler(nil)
	body := `{"name":"Ivan","surname":"Ivanov"}`
	req := httptest.NewRequest(http.MethodPut, "/people/123", strings.NewReader(body))
	rr := httptest.NewRecorder()

	h.UpdatePerson(rr, req)

	if rr.Code != http.StatusNotImplemented {
		t.Fatalf("PUT /people/{id}: expected %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}

func TestDeletePerson_NotImplemented(t *testing.T) {
	h := NewHandler(nil)
	req := httptest.NewRequest(http.MethodDelete, "/people/123", nil)
	rr := httptest.NewRecorder()

	h.DeletePerson(rr, req)

	if rr.Code != http.StatusNotImplemented {
		t.Fatalf("DELETE /people/{id}: expected %d, got %d", http.StatusNotImplemented, rr.Code)
	}
}
