package repository

import (
	"context"
	"errors"
	"person-enricher/internal/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), nil)
	if err != nil {
		panic(err)
	}

	return gormDB, mock
}

func TestGormPersonRepository_Create(t *testing.T) {
	db, mock := NewMockDB()
	repo := NewPersonRepository(db)

	tests := []struct {
		name      string
		input     models.Person
		mock      func()
		want      models.Person
		expectErr bool
	}{
		{
			name: "success",
			input: models.Person{
				Name: "John", Surname: "Doe",
				Age: 30, Gender: "male", Nationality: "US",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "people"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
				mock.ExpectCommit()
			},
			want: models.Person{
				ID: "1", Name: "John", Surname: "Doe",
				Age: 30, Gender: "male", Nationality: "US",
			},
		},
		{
			name: "database error",
			input: models.Person{
				Name: "John", Surname: "Doe",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "people"`)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			result, err := repo.Create(context.Background(), tt.input)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, result.ID)
				assert.Equal(t, tt.input.Name, result.Name)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGormPersonRepository_List(t *testing.T) {
    db, mock := NewMockDB()
    repo := NewPersonRepository(db)

    now := time.Now()
    wantPeople := []models.Person{
        {
            ID:        "1",
            Name:      "John",
            CreatedAt: now,
            UpdatedAt: now, 
        },
        {
            ID:        "2",
            Name:      "Alice",
            CreatedAt: now,
            UpdatedAt: now, 
        },
    }

    tests := []struct {
        name      string
        filter    models.PeopleFilter
        mockSetup func()
        want      []models.Person
        expectErr bool
    }{
        {
            name: "success with filter",
            filter: models.PeopleFilter{
                Filter: "John",
                Page:   1,
                Size:   10,
            },
            mockSetup: func() {
                sql := regexp.QuoteMeta(
                    `SELECT * FROM "people" WHERE ` +
                        `(name ILIKE $1 OR surname ILIKE $2 OR patronymic ILIKE $3) ` +
                        `AND "people"."deleted_at" IS NULL ORDER BY id LIMIT $4`,
                )
                mock.ExpectQuery(sql).
                    WithArgs("%John%", "%John%", "%John%", 10).
                    WillReturnRows(sqlmock.NewRows([]string{
                        "id", "name", "surname", "patronymic",
                        "age", "gender", "nationality",
                        "created_at", "updated_at",
                    }).
                        AddRow("1", "John", "", "", 0, "", "", now, now).
                        AddRow("2", "Alice", "", "", 0, "", "", now, now),
                    )
            },
            want: wantPeople,
        },
        {
            name: "database error",
            filter: models.PeopleFilter{
                Page: 1,
                Size: 10,
            },
            mockSetup: func() {
                sql := regexp.QuoteMeta(
                    `SELECT * FROM "people" WHERE "people"."deleted_at" IS NULL ORDER BY id LIMIT $1`,
                )
                mock.ExpectQuery(sql).
                    WithArgs(10).
                    WillReturnError(errors.New("db error"))
            },
            expectErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := repo.List(context.Background(), tt.filter)

            if tt.expectErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
            assert.NoError(t, mock.ExpectationsWereMet())
        })
    }
}



func TestGormPersonRepository_GetByID(t *testing.T) {
    db, mock := NewMockDB()
    repo := NewPersonRepository(db)
    now := time.Now()

    tests := []struct {
        name      string
        id        string
        mockSetup func()
        wantID    string
        expectErr bool
    }{
        {
            name: "success",
            id:   "1",
            mockSetup: func() {
                sql := regexp.QuoteMeta(
                    `SELECT * FROM "people" WHERE id = $1 AND "people"."deleted_at" IS NULL ORDER BY "people"."id" LIMIT $2`,
                )
                mock.ExpectQuery(sql).
                    WithArgs("1", 1).
                    WillReturnRows(sqlmock.NewRows([]string{
                        "id", "name", "surname", "patronymic",
                        "age", "gender", "nationality",
                        "created_at", "updated_at", "deleted_at",
                    }).
                        AddRow("1", "John", "", "", 0, "", "", now, now, nil),
                    )
            },
            wantID: "1",
        },
        {
            name: "not found",
            id:   "2",
            mockSetup: func() {
                sql := regexp.QuoteMeta(
                    `SELECT * FROM "people" WHERE id = $1 AND "people"."deleted_at" IS NULL ORDER BY "people"."id" LIMIT $2`,
                )
                mock.ExpectQuery(sql).
                    WithArgs("2", 1).
                    WillReturnError(gorm.ErrRecordNotFound)
            },
            expectErr: false,
        },
        {
            name: "database error",
            id:   "3",
            mockSetup: func() {
                sql := regexp.QuoteMeta(
                    `SELECT * FROM "people" WHERE id = $1 AND "people"."deleted_at" IS NULL ORDER BY "people"."id" LIMIT $2`,
                )
                mock.ExpectQuery(sql).
                    WithArgs("3", 1).
                    WillReturnError(errors.New("db error"))
            },
            expectErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := repo.GetByID(context.Background(), tt.id)

            if tt.expectErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantID, got.ID)
            }

            // Проверяем, что все ожидания sqlmock выполнились
            assert.NoError(t, mock.ExpectationsWereMet())
        })
    }
}


func TestGormPersonRepository_Update(t *testing.T) {
    db, mock := NewMockDB()
    repo := NewPersonRepository(db)
    now := time.Now()

    tests := []struct {
        name      string
        input     models.Person
        mockSetup func()
        expectErr bool
    }{
        {
            name: "success",
            input: models.Person{
                ID:   "1",
                Name: "John Updated",
            },
            mockSetup: func() {

                mock.ExpectBegin()
                mock.ExpectExec(regexp.QuoteMeta(
                    `UPDATE "people" SET "id"=$1,"name"=$2,"updated_at"=$3 WHERE id = $4 AND "people"."deleted_at" IS NULL`,
                )).
                    WithArgs("1", "John Updated", sqlmock.AnyArg(), "1").
                    WillReturnResult(sqlmock.NewResult(0, 1))
                mock.ExpectCommit()

                query := regexp.QuoteMeta(
                    `SELECT * FROM "people" WHERE id = $1 AND "people"."deleted_at" IS NULL ORDER BY "people"."id" LIMIT $2`,
                )
                mock.ExpectQuery(query).
                    WithArgs("1", 1).
                    WillReturnRows(sqlmock.NewRows([]string{
                        "id", "name", "surname", "patronymic",
                        "age", "gender", "nationality",
                        "created_at", "updated_at", "deleted_at",
                    }).
                        AddRow("1", "John Updated", "", "", 0, "", "", now, now, nil),
                    )
            },
        },
        {
            name: "database error",
            input: models.Person{
                ID:   "2",
                Name: "Invalid",
            },
            mockSetup: func() {
                mock.ExpectBegin()
                mock.ExpectExec(regexp.QuoteMeta(
                    `UPDATE "people" SET "id"=$1,"name"=$2,"updated_at"=$3 WHERE id = $4 AND "people"."deleted_at" IS NULL`,
                )).
                    WithArgs("2", "Invalid", sqlmock.AnyArg(), "2").
                    WillReturnError(errors.New("db error"))
                mock.ExpectRollback()
            },
            expectErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            result, err := repo.Update(context.Background(), tt.input)

            if tt.expectErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.input.Name, result.Name)
            }
            assert.NoError(t, mock.ExpectationsWereMet())
        })
    }
}


func TestGormPersonRepository_Delete(t *testing.T) {
    db, mock := NewMockDB()
    repo := NewPersonRepository(db)

    tests := []struct {
        name      string
        id        string
        mockSetup func()
        expectErr bool
    }{
        {
            name: "success",
            id:   "1",
            mockSetup: func() {
                mock.ExpectBegin()
                mock.ExpectExec(regexp.QuoteMeta(
                    `UPDATE "people" SET "deleted_at"=$1 WHERE id = $2 AND "people"."deleted_at" IS NULL`,
                )).
                    WithArgs(sqlmock.AnyArg(), "1").
                    WillReturnResult(sqlmock.NewResult(0, 1))
                mock.ExpectCommit()
            },
        },
        {
            name: "not found",
            id:   "2",
            mockSetup: func() {
                mock.ExpectBegin()
                mock.ExpectExec(regexp.QuoteMeta(
                    `UPDATE "people" SET "deleted_at"=$1 WHERE id = $2 AND "people"."deleted_at" IS NULL`,
                )).
                    WithArgs(sqlmock.AnyArg(), "2").
                    WillReturnResult(sqlmock.NewResult(0, 0))
                mock.ExpectCommit()
            },
            // при RowsAffected = 0 ошибок не будет
            expectErr: false,
        },
        {
            name: "database error",
            id:   "3",
            mockSetup: func() {
                mock.ExpectBegin()
                mock.ExpectExec(regexp.QuoteMeta(
                    `UPDATE "people" SET "deleted_at"=$1 WHERE id = $2 AND "people"."deleted_at" IS NULL`,
                )).
                    WithArgs(sqlmock.AnyArg(), "3").
                    WillReturnError(errors.New("db error"))
                mock.ExpectRollback()
            },
            expectErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            err := repo.Delete(context.Background(), tt.id)

            if tt.expectErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
            assert.NoError(t, mock.ExpectationsWereMet())
        })
    }
}

