package service

import (
	"context"
	"fmt"
	"log"
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
	log.Printf("service.CreatePerson: creating person")
	// Get person age, gender, nationality
	age, err := s.enricher.GetPersonAge(ctx, req.Name)
	if err != nil {
		log.Printf("service.CreatePerson: could not get person age: %v", err)
		return models.Person{}, fmt.Errorf("could not get person age: %w", err)
	}
	gender, err := s.enricher.GetPersonGender(ctx, req.Name)
	if err != nil {
		log.Printf("service.CreatePerson: could not get person gender: %v", err)
		return models.Person{}, fmt.Errorf("could not get person gender: %w", err)
	}
	nationality, err := s.enricher.GetPersonNationality(ctx, req.Name)
	if err != nil {
		log.Printf("service.CreatePerson: could not get person nationality: %v", err)
		return models.Person{}, fmt.Errorf("could not get person nationality: %w", err)
	}
	// Build a model to save
	log.Printf("service.CreatePerson: building person model")
	person := models.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	// Save the person in the repository
	log.Printf("service.CreatePerson: saving person in repository")
	createdPerson, err := s.repo.Create(ctx, person)
	if err != nil {
		log.Printf("service.CreatePerson: could not create person: %v", err)
		return models.Person{}, fmt.Errorf("could not create person: %w", err)
	}

	log.Printf("service.CreatePerson: person created")
	return createdPerson, nil

}

// GetPeople retrieves a list of people based on the provided filter criteria.
// It calls the repository List method with the given context and filter.
// Returns a slice of Person models if successful, otherwise returns an error.
func (s *personService) GetPeople(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error) {
	log.Printf("service.GetPeople: getting people")
	people, err := s.repo.List(ctx, filter)
	if err != nil {
		log.Printf("service.GetPeople: could not list people: %v", err)
		return nil, fmt.Errorf("could not list people: %w", err)
	}
	log.Printf("service.GetPeople: people listed")
	return people, nil
}

// GetPersonByID retrieves a person by their unique identifier.
// It calls the repository GetByID method with the provided context and id.
// Returns the person if found, otherwise returns an error.
func (s *personService) GetPersonByID(ctx context.Context, id string) (models.Person, error) {
	// Save the person in the repository
	log.Printf("service.GetPersonByID: getting person by id")
	gotPerson, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("service.GetPersonByID: could not get person: %v", err)
		return models.Person{}, fmt.Errorf("could not got person: %w", err)
	}

	log.Printf("service.GetPersonByID: person got %v", gotPerson)
	return gotPerson, nil
}

// UpdatePerson updates a person by their unique identifier.
// It calls the repository Update method with the provided context and the updated person.
// Returns the updated person if the update was successful, otherwise returns an error.
func (s *personService) UpdatePerson(ctx context.Context, id string, req models.UpdatePersonRequest) (models.Person, error) {
	log.Printf("service.UpdatePerson: updating person")
	updatedPerson := models.Person{
		ID:          id,
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         req.Age,
		Gender:      req.Gender,
		Nationality: req.Nationality,
	}
	log.Printf("service.UpdatePerson: updating person in repository")
	updatedPerson, err := s.repo.Update(ctx, updatedPerson)
	if err != nil {
		log.Printf("service.UpdatePerson: could not update person: %v", err)
		return models.Person{}, fmt.Errorf("could not update person: %w", err)
	}
	log.Printf("service.UpdatePerson: person updated")
	return updatedPerson, nil
}

// DeletePerson deletes a person by their unique identifier.
// It calls the repository Delete method with the provided context and id.
// Returns an error if the deletion fails.
func (s *personService) DeletePerson(ctx context.Context, id string) error {
	log.Printf("service.DeletePerson: deleting person")
	if err := s.repo.Delete(ctx, id); err != nil {
		log.Printf("service.DeletePerson: could not delete person: %v", err)
		return fmt.Errorf("could not delete person: %w", err)
	}
	log.Printf("service.DeletePerson: person deleted")
	return nil
}
