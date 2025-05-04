package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
	log.Printf("main: Starting person-enricher")
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, reading environment variables directly")
	}

	// 2) Get environment variables from .env
	log.Printf("main: Loading environment variables")
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
	log.Printf("DB_HOST: %s", dbHost)
	log.Printf("DB_PORT: %d", dbPort)
	log.Printf("DB_USER: %s", dbUser)
	log.Printf("DB_PASSWORD: %s", dbPass)
	log.Printf("DB_NAME: %s", dbName)
	log.Printf("DB_SSLMODE: %s", dbSSL)
	log.Printf("HTTP_PORT: %s", httpPort)
	log.Printf("METRICS_PORT: %s", metricsPort)

	// 3) Connect to DB and initialize repository
	log.Printf("main: Connecting to DB")
	db, err := repository.NewDB(dbHost, dbPort, dbUser, dbPass, dbName, dbSSL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Initialize handlers metrics
	metrics.InitMetrics()

	// 4) Initialize services and handlers
	log.Printf("main: Initializing services")
	repo := repository.NewPersonRepository(db)
	metricsRepo := repository.NewMetricsRepository(repo)

	enricher := externalapi.NewPersonalDataEnricher()
	metricsEnricher := externalapi.NewMetricsEnricher(enricher)

	svc := service.NewPersonService(metricsRepo, metricsEnricher)
	instrumentedSvc := handlers.NewInstrumentedService(svc)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// 5) Initialize router
	log.Printf("main: Initializing router")
	handler := handlers.NewHandler(instrumentedSvc)
	router := handlers.NewRouter(handler)

	// 6) Start HTTP server
	log.Printf("main: Starting HTTP server")
	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: router,
	}
	go func() {
		log.Printf("Server listening on %s …", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	log.Printf("main: Starting metrics server")
	metricsRouter := http.NewServeMux()
	metricsRouter.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    metricsPort,
		Handler: metricsRouter,
	}
	// Metrics server
	go func() {
		log.Printf("Metrics server listening on %s", metricsPort)

		if err := metricsServer.ListenAndServe(); err != nil {
			log.Fatalf("Metrics server failed: %v", err)
		}
	}()

	// 7) Wait for termination signal
	<-stop

	// 8) Gracefully shutdown servers
	log.Printf("main: Stopping http servers")
	ctx := context.Background()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		log.Printf("Metrics server shutdown error: %v", err)
	}

}
