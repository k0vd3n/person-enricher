package service

import (
	"context"
	"errors"
	"person-enricher/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestCreatePerson(t *testing.T) {
	tests := []struct {
		name        string
		req         models.CreatePersonRequest
		mockSetup   func(*MockEnricher, *MockRepository)
		want        models.Person
		expectedErr string
	}{
		{
			name: "successful creation",
			req: models.CreatePersonRequest{
				Name:       "John",
				Surname:    "Doe",
				Patronymic: "Smith",
			},
			mockSetup: func(e *MockEnricher, r *MockRepository) {
				e.On("GetPersonAge", mock.Anything, "John").Return(30, nil)
				e.On("GetPersonGender", mock.Anything, "John").Return("male", nil)
				e.On("GetPersonNationality", mock.Anything, "John").Return("US", nil)
				r.On("Create", mock.Anything, models.Person{
					Name:        "John",
					Surname:     "Doe",
					Patronymic:  "Smith",
					Age:         30,
					Gender:      "male",
					Nationality: "US",
				}).Return(models.Person{ID: "1"}, nil)
			},
			want: models.Person{ID: "1"},
		},
		{
			name: "age enrichment error",
			req: models.CreatePersonRequest{
				Name:    "John",
				Surname: "Doe",
			},
			mockSetup: func(e *MockEnricher, r *MockRepository) {
				e.On("GetPersonAge", mock.Anything, "John").Return(0, errors.New("api error"))
			},
			expectedErr: "could not get person age: api error",
		},
		{
			name: "gender enrichment error",
			req: models.CreatePersonRequest{
				Name:    "John",
				Surname: "Doe",
			},
			mockSetup: func(e *MockEnricher, r *MockRepository) {
				e.On("GetPersonAge", mock.Anything, "John").Return(30, nil)
				e.On("GetPersonGender", mock.Anything, "John").Return("", errors.New("api error"))
			},
			expectedErr: "could not get person gender: api error",
		},
		{
			name: "nationality enrichment error",
			req: models.CreatePersonRequest{
				Name:    "John",
				Surname: "Doe",
			},
			mockSetup: func(e *MockEnricher, r *MockRepository) {
				e.On("GetPersonAge", mock.Anything, "John").Return(30, nil)
				e.On("GetPersonGender", mock.Anything, "John").Return("male", nil)
				e.On("GetPersonNationality", mock.Anything, "John").Return("", errors.New("api error"))
			},
			expectedErr: "could not get person nationality: api error",
		},
		{
			name: "repository error",
			req: models.CreatePersonRequest{
				Name:    "John",
				Surname: "Doe",
			},
			mockSetup: func(e *MockEnricher, r *MockRepository) {
				e.On("GetPersonAge", mock.Anything, "John").Return(30, nil)
				e.On("GetPersonGender", mock.Anything, "John").Return("male", nil)
				e.On("GetPersonNationality", mock.Anything, "John").Return("US", nil)
				r.On("Create", mock.Anything, mock.Anything).Return(models.Person{}, errors.New("db error"))
			},
			expectedErr: "could not create person: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enricher := new(MockEnricher)
			repo := new(MockRepository)
			tt.mockSetup(enricher, repo)

			service := NewPersonService(repo, enricher)
			result, err := service.CreatePerson(context.Background(), tt.req)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}

			enricher.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}

func TestGetPeople(t *testing.T) {
    tests := []struct {
        name        string
        filter      models.PeopleFilter
        mockSetup   func(*MockRepository)
        expected    []models.Person
        expectedErr string
    }{
        {
            name: "successful list",
            filter: models.PeopleFilter{
                Page: 1,
                Size: 10,
            },
            mockSetup: func(r *MockRepository) {
                r.On("List", mock.Anything, models.PeopleFilter{
                    Page: 1,
                    Size: 10,
                }).Return([]models.Person{{ID: "1"}}, nil)
            },
            expected: []models.Person{{ID: "1"}},
        },
        {
            name: "repository list error",
            filter: models.PeopleFilter{
                Page: 1,
                Size: 10,
            },
            mockSetup: func(r *MockRepository) {
                r.On("List", mock.Anything, models.PeopleFilter{
                    Page: 1,
                    Size: 10,
                }).Return([]models.Person{}, errors.New("db error"))
            },
            expectedErr: "could not list people: db error",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(MockRepository)
            enricher := new(MockEnricher)
            tt.mockSetup(repo)

            service := NewPersonService(repo, enricher)
            result, err := service.GetPeople(context.Background(), tt.filter)

            if tt.expectedErr != "" {
                assert.ErrorContains(t, err, tt.expectedErr)
                assert.Nil(t, result) 
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
            repo.AssertExpectations(t)
        })
    }
}

func TestGetPersonByID(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		mockSetup   func(*MockRepository)
		expected    models.Person
		expectedErr string
	}{
		{
			name: "success",
			id:   "1",
			mockSetup: func(r *MockRepository) {
				r.On("GetByID", mock.Anything, "1").Return(models.Person{ID: "1"}, nil)
			},
			expected: models.Person{ID: "1"},
		},
		{
			name: "not found",
			id:   "2",
			mockSetup: func(r *MockRepository) {
				r.On("GetByID", mock.Anything, "2").Return(models.Person{}, nil)
			},
			expected: models.Person{},
		},
		{
			name: "repository error",
			id:   "3",
			mockSetup: func(r *MockRepository) {
				r.On("GetByID", mock.Anything, "3").Return(models.Person{}, errors.New("db error"))
			},
			expectedErr: "could not got person: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRepository)
			tt.mockSetup(repo)

			service := NewPersonService(repo, new(MockEnricher))
			result, err := service.GetPersonByID(context.Background(), tt.id)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestUpdatePerson(t *testing.T) {
    tests := []struct {
        name        string
        id          string
        req         models.UpdatePersonRequest
        mockSetup   func(*MockRepository)
        expected    models.Person
        expectedErr string
    }{
        {
            name: "successful update",
            id:   "1",
            req: models.UpdatePersonRequest{
                Name:        "John",
                Surname:     "Doe",
                Patronymic:  "Smith",
                Age:         30,
                Gender:      "male",
                Nationality: "US",
            },
            mockSetup: func(r *MockRepository) {
                expectedPerson := models.Person{
                    ID:          "1",
                    Name:        "John",
                    Surname:     "Doe",
                    Patronymic:  "Smith",
                    Age:         30,
                    Gender:      "male",
                    Nationality: "US",
                }
                r.On("Update", mock.Anything, expectedPerson).Return(expectedPerson, nil)
            },
            expected: models.Person{
                ID:          "1",
                Name:        "John",
                Surname:     "Doe",
                Patronymic:  "Smith",
                Age:         30,
                Gender:      "male",
                Nationality: "US",
            },
        },
        {
            name: "repository update error",
            id:   "1",
            req: models.UpdatePersonRequest{
                Name:        "John",
                Surname:     "Doe",
                Patronymic:  "Smith",
                Age:         30,
                Gender:      "male",
                Nationality: "US",
            },
            mockSetup: func(r *MockRepository) {
                expectedPerson := models.Person{
                    ID:          "1",
                    Name:        "John",
                    Surname:     "Doe",
                    Patronymic:  "Smith",
                    Age:         30,
                    Gender:      "male",
                    Nationality: "US",
                }
                r.On("Update", mock.Anything, expectedPerson).Return(models.Person{}, errors.New("db error"))
            },
            expectedErr: "could not update person: db error",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(MockRepository)
            enricher := new(MockEnricher)
            tt.mockSetup(repo)

            service := NewPersonService(repo, enricher)
            result, err := service.UpdatePerson(context.Background(), tt.id, tt.req)

            if tt.expectedErr != "" {
                assert.ErrorContains(t, err, tt.expectedErr)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
            repo.AssertExpectations(t)
        })
    }
}

func TestDeletePerson(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		mockSetup   func(*MockRepository)
		expectedErr string
	}{
		{
			name: "success",
			id:   "1",
			mockSetup: func(r *MockRepository) {
				r.On("Delete", mock.Anything, "1").Return(nil)
			},
		},
		{
			name: "error",
			id:   "2",
			mockSetup: func(r *MockRepository) {
				r.On("Delete", mock.Anything, "2").Return(errors.New("db error"))
			},
			expectedErr: "could not delete person: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRepository)
			tt.mockSetup(repo)

			service := NewPersonService(repo, new(MockEnricher))
			err := service.DeletePerson(context.Background(), tt.id)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}