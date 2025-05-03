package repository

import (
	"fmt"
	"person-enricher/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB creates a new GORM database connection with PostgreSQL driver.
//
// It takes the database host, port, user, password, database name, and SSL mode as parameters.
// The connection string is built from these parameters.
// The logger is set to log mode Info.
// The SQL database is set to have a maximum of 10 idle connections and 100 open connections.
// The connection lifetime is set to 1 hour.
// The people table is auto-migrated using the Person struct.
// Returns the *gorm.DB and an error, if any.
func NewDB(
	host string,
	port int,
	user, password, dbname, sslmode string,
) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	// Настройка пула соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migration: creates/updates the people table by struct Person
	if err := db.AutoMigrate(&models.Person{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return db, nil

}
