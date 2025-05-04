package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"person-enricher/internal/externalapi"
	"person-enricher/internal/handlers"
	"person-enricher/internal/metrics"
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
		httpPort = ":8080"
	}

	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = ":8081"
	}

	// 3) Connect to DB and initialize repository
	db, err := repository.NewDB(dbHost, dbPort, dbUser, dbPass, dbName, dbSSL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Initialize handlers metrics
	metrics.InitMetrics()

	// 4) Initialize services and handlers
	repo := repository.NewPersonRepository(db)
	metricsRepo := repository.NewMetricsRepository(repo)

	enricher := externalapi.NewPersonalDataEnricher()
	metricsEnricher := externalapi.NewMetricsEnricher(enricher)

	svc := service.NewPersonService(metricsRepo, metricsEnricher)
	instrumentedSvc := handlers.NewInstrumentedService(svc)

	go func() {
		metricsRouter := http.NewServeMux()
		metricsRouter.Handle("/metrics", promhttp.Handler())
		log.Printf("Metrics server listening on %s", metricsPort)
		if err := http.ListenAndServe(metricsPort, metricsRouter); err != nil {
			log.Fatalf("Metrics server failed: %v", err)
		}
	}()

	// 5) Initialize router
	handler := handlers.NewHandler(instrumentedSvc)
	router := handlers.NewRouter(handler)

	// 6) Start HTTP server
	log.Printf("Server listening on %s …", httpPort)
	if err := http.ListenAndServe(httpPort, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
