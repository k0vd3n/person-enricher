package repository

import (
	"context"
	"person-enricher/internal/models"
)

// PersonRepository represents a repository for managing person data.
type PersonRepository interface {
	Create(ctx context.Context, p models.Person) (models.Person, error)
	List(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error)
	GetByID(ctx context.Context, id string) (models.Person, error)
	Update(ctx context.Context, p models.Person) (models.Person, error)
	Delete(ctx context.Context, id string) error
}
