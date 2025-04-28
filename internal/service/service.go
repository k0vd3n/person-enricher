package service

import (
	"context"
	"fmt"
	"person-enricher/internal/externalapi"
	"person-enricher/internal/models"
	"person-enricher/internal/repository"
)

type PersonService interface {
	CreatePerson(ctx context.Context, req models.CreatePersonRequest) (models.Person, error)
	GetPeople(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error)
	GetPersonByID(ctx context.Context, id string) (models.Person, error)
	UpdatePerson(ctx context.Context, id string, req models.UpdatePersonRequest) (models.Person, error)
	DeletePerson(ctx context.Context, id string) error
}

type personService struct {
	repo     repository.PersonRepository
	enricher externalapi.EnrichPersonalData
}

func NewPersonService(
	repo repository.PersonRepository,
	enricher externalapi.EnrichPersonalData,
) PersonService {
	return &personService{
		repo:     repo,
		enricher: enricher,
	}
}

// CreatePerson creates and enrich person via external API
func (s *personService) CreatePerson(ctx context.Context, req models.CreatePersonRequest) (models.Person, error) {
	// Get person age, gender, nationality
	age, err := s.enricher.GetPersonAge(ctx, req.Name)
	if err != nil {
		return models.Person{}, fmt.Errorf("could not get person age: %w", err)
	}
	gender, err := s.enricher.GetPersonGender(ctx, req.Name)
	if err != nil {
		return models.Person{}, fmt.Errorf("could not get person gender: %w", err)
	}
	nationality, err := s.enricher.GetPersonNationality(ctx, req.Name)
	if err != nil {
		return models.Person{}, fmt.Errorf("could not get person nationality: %w", err)
	}
	// Build a model to save
	person := models.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	// Save in repository
	return s.repo.Create(ctx, person)

}

func (s *personService) GetPeople(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error) {
	return s.repo.List(ctx, filter)
}

func (s *personService) GetPersonByID(ctx context.Context, id string) (models.Person, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *personService) UpdatePerson(ctx context.Context, id string, req models.UpdatePersonRequest) (models.Person, error) {
	updated := models.Person{
		ID:          id,
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         req.Age,
		Gender:      req.Gender,
		Nationality: req.Nationality,
	}
	return s.repo.Update(ctx, updated)
}

func (s *personService) DeletePerson(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}