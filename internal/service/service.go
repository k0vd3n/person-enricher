package service

import (
	"context"
	"person-enricher/internal/models"
)

type PersonService interface {
	CreatePerson(ctx context.Context, req models.CreatePersonRequest) (models.Person, error)
	GetPeople(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error)
	GetPersonByID(ctx context.Context, id string) (models.Person, error)
	UpdatePerson(ctx context.Context, id string, req models.UpdatePersonRequest) (models.Person, error)
	DeletePerson(ctx context.Context, id string) error
}
