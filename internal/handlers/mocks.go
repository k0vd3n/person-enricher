package handlers

import (
	"context"
	"errors"
	"person-enricher/internal/models"
)

type MockPersonService struct{}

func (m *MockPersonService) GetPeople(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error) {
	if filter.Filter == "error" {
		return nil, errors.New("service error")
	}
	return []models.Person{{ID: "1", Name: "Test", Surname: "User"}}, nil
}

func (m *MockPersonService) GetPersonByID(ctx context.Context, id string) (models.Person, error) {
	switch id {
	case "error-id":
		return models.Person{}, errors.New("service error")
	case "notfound-id":
		return models.Person{}, nil
	default:
		return models.Person{ID: id, Name: "Test", Surname: "User"}, nil
	}
}

func (m *MockPersonService) CreatePerson(ctx context.Context, req models.CreatePersonRequest) (models.Person, error) {
	if req.Name == "error" {
		return models.Person{}, errors.New("service error")
	}
	return models.Person{ID: "new-id", Name: req.Name, Surname: req.Surname}, nil
}

func (m *MockPersonService) UpdatePerson(ctx context.Context, id string, req models.UpdatePersonRequest) (models.Person, error) {
	if id == "error-id" {
		return models.Person{}, errors.New("service error")
	}
	return models.Person{
		ID:          id,
		Name:        req.Name,
		Surname:     req.Surname,
		Age:         req.Age,
		Gender:      req.Gender,
		Nationality: req.Nationality,
	}, nil
}

func (m *MockPersonService) DeletePerson(ctx context.Context, id string) error {
	if id == "error-id" {
		return errors.New("service error")
	}
	return nil
}
