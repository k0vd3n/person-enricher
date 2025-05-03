package service

import (
	"context"
	"fmt"
	"person-enricher/internal/externalapi"
	"person-enricher/internal/models"
	"person-enricher/internal/repository"
)

// PersonService represents a person service.
type PersonService interface {
	CreatePerson(ctx context.Context, req models.CreatePersonRequest) (models.Person, error)
	GetPeople(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error)
	GetPersonByID(ctx context.Context, id string) (models.Person, error)
	UpdatePerson(ctx context.Context, id string, req models.UpdatePersonRequest) (models.Person, error)
	DeletePerson(ctx context.Context, id string) error
}

// personService struct represents a person service
type personService struct {
	repo     repository.PersonRepository
	enricher externalapi.EnrichPersonalData
}

// NewPersonService creates a new instance of personService with the provided repository
// and data enricher. It returns a PersonService interface, which provides methods for
// managing person data, including creating, updating, deleting, and retrieving persons.
func NewPersonService(
	repo repository.PersonRepository,
	enricher externalapi.EnrichPersonalData,
) PersonService {
	return &personService{
		repo:     repo,
		enricher: enricher,
	}
}

// CreatePerson creates a new person with the given name, surname, patronymic, and
// attempts to enrich the person data with age, gender, and nationality using
// the external APIs. It returns the created person and an error, if any.
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

// GetPeople retrieves a list of people based on the provided filter criteria.
// It calls the repository List method with the given context and filter.
// Returns a slice of Person models if successful, otherwise returns an error.
func (s *personService) GetPeople(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error) {
	people, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("could not list people: %w", err)
	}
	return people, nil
}

// GetPersonByID retrieves a person by their unique identifier.
// It calls the repository GetByID method with the provided context and id.
// Returns the person if found, otherwise returns an error.
func (s *personService) GetPersonByID(ctx context.Context, id string) (models.Person, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdatePerson updates a person by their unique identifier.
// It calls the repository Update method with the provided context and the updated person.
// Returns the updated person if the update was successful, otherwise returns an error.
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

// DeletePerson deletes a person by their unique identifier.
// It calls the repository Delete method with the provided context and id.
// Returns an error if the deletion fails.
func (s *personService) DeletePerson(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
