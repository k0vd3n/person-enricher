package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"person-enricher/internal/handlers"
	"person-enricher/internal/externalapi"
	"person-enricher/internal/repository"
	"person-enricher/internal/service"
)

func main() {
	// 1) Load environment variables from .env file 
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, reading environment variables directly")
	}

	// 2) Get environment variables from .env
	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("invalid DB_PORT: %v", err)
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSL := os.Getenv("DB_SSLMODE")

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// 3) Connect to DB and initialize repository
	db, err := repository.NewDB(dbHost, dbPort, dbUser, dbPass, dbName, dbSSL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// 4) Initialize services and handlers 
	repo := repository.NewPersonRepository(db)
	enricher := externalapi.NewPersonalDataEnricher()
	svc := service.NewPersonService(repo, enricher)

	// 5) Initialize router
	handler := api.NewHandler(svc)
	router := api.NewRouter(handler)

	// 6) Start HTTP server
	addr := fmt.Sprintf(":%s", httpPort)
	log.Printf("Server listening on %s …", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
