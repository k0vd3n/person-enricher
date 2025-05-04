package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"person-enricher/internal/models"

	"gorm.io/gorm"
)

// PersonRepository represents a repository for managing person data.
type PersonRepository interface {
	Create(ctx context.Context, p models.Person) (models.Person, error)
	List(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error)
	GetByID(ctx context.Context, id string) (models.Person, error)
	Update(ctx context.Context, p models.Person) (models.Person, error)
	Delete(ctx context.Context, id string) error
}

type GormPersonRepository struct {
	db *gorm.DB
}

// NewPersonRepository returns a new PersonRepository using the given *gorm.DB.
//
// It initializes a GormPersonRepository and returns it as a PersonRepository.
func NewPersonRepository(db *gorm.DB) PersonRepository {
	return &GormPersonRepository{db: db}
}

// Create adds a new person to the repository.
// It uses the provided context for request scoping and the person model
// for the data to be stored. It returns the created person and any error
// encountered during the operation.
func (r *GormPersonRepository) Create(ctx context.Context, p models.Person) (models.Person, error) {
	log.Printf("GormPersonRepository.Create: creating person")
	if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
		return models.Person{}, err
	}
	return p, nil
}

// List retrieves a list of people from the repository based on the provided filter criteria.
// It uses the given context for request scoping and applies filtering, pagination, and sorting.
// The filter allows searching by name, surname, and patronymic using a case-insensitive match.
// Pagination is controlled by the Page and Size fields in the filter, and results are ordered by ID.
// Returns a slice of Person models if successful, otherwise returns an error.
func (r *GormPersonRepository) List(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error) {
	log.Printf("GormPersonRepository.List: listing people")
	var people []models.Person
	q := r.db.WithContext(ctx)

	if f := filter.Filter; f != "" {
		like := "%" + f + "%"
		q = q.Where(
			"name ILIKE ? OR surname ILIKE ? OR patronymic ILIKE ?",
			like, like, like,
		)
	}
	log.Printf("GormPersonRepository.List: filter: %v", filter)

	// pagination
	offset := (filter.Page - 1) * filter.Size
	log.Printf("GormPersonRepository.List: offset: %d page: %d, size: %d", offset, filter.Page, filter.Size)
	if err := q.
		Limit(filter.Size).
		Offset(offset).
		Order("id").
		Find(&people).Error; err != nil {
		log.Printf("GormPersonRepository.List: could not list people: %v", err)
		return nil, fmt.Errorf("list people: %w", err)
	}
	return people, nil
}

// GetByID retrieves a person from the repository by their unique identifier.
// It uses the provided context for request scoping and queries the database using the given ID.
// If the person is found, it returns the person model; otherwise, it returns an error.
// If the record is not found, it returns ErrNotFound, otherwise it returns a wrapped error.
func (r *GormPersonRepository) GetByID(ctx context.Context, id string) (models.Person, error) {
	log.Printf("GormPersonRepository.GetByID: getting person by id")
	var p models.Person
	err := r.db.WithContext(ctx).
		First(&p, "id = ?", id).
		Error
	log.Printf("GormPersonRepository.GetByID: err: %v", err)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("GormPersonRepository.GetByID: person not found")
		return models.Person{}, nil
	} else if err != nil {
		log.Printf("GormPersonRepository.GetByID: could not get person by id: %v", err)
		return models.Person{}, fmt.Errorf("get by id: %w", err)
	}
	log.Printf("GormPersonRepository.GetByID: person found")

	return p, nil
}

// Update updates a person by their unique identifier.
// It uses the provided context for request scoping and the person model
// for the data to be stored. It returns the updated person and any error
// encountered during the operation.
func (r *GormPersonRepository) Update(ctx context.Context, p models.Person) (models.Person, error) {
	log.Printf("GormPersonRepository.Update: updating person")
	if err := r.db.WithContext(ctx).
		Model(&models.Person{}).
		Where("id = ?", p.ID).
		Updates(p).
		Error; err != nil {
		log.Printf("GormPersonRepository.Update: could not update person: %v", err)
		return models.Person{}, fmt.Errorf("update person: %w", err)
	}
	log.Printf("GormPersonRepository.Update: person updated")
	return r.GetByID(ctx, p.ID) // Возвращаем обновлённую запись
}

// Delete removes a person from the repository by their unique identifier.
// It uses the provided context for request scoping and deletes the person
// record from the database using the given ID. If the deletion is successful,
// it returns nil; otherwise, it returns a wrapped error indicating the failure.
func (r *GormPersonRepository) Delete(ctx context.Context, id string) error {
	log.Printf("GormPersonRepository.Delete: deleting person")
	if err := r.db.WithContext(ctx).
		Delete(&models.Person{}, "id = ?", id).
		Error; err != nil {
			log.Printf("GormPersonRepository.Delete: could not delete person: %v", err)
		return fmt.Errorf("delete person: %w", err)
	}
	log.Printf("GormPersonRepository.Delete: person deleted")
	return nil
}
