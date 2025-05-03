package service

import (
	"context"
	"person-enricher/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, p models.Person) (models.Person, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(models.Person), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]models.Person), args.Error(1)
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (models.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Person), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, p models.Person) (models.Person, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(models.Person), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockEnricher struct {
	mock.Mock
}

func (m *MockEnricher) GetPersonAge(ctx context.Context, name string) (int, error) {
	args := m.Called(ctx, name)
	return args.Int(0), args.Error(1)
}

func (m *MockEnricher) GetPersonGender(ctx context.Context, name string) (string, error) {
	args := m.Called(ctx, name)
	return args.String(0), args.Error(1)
}

func (m *MockEnricher) GetPersonNationality(ctx context.Context, name string) (string, error) {
	args := m.Called(ctx, name)
	return args.String(0), args.Error(1)
}
