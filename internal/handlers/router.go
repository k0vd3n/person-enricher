package handlers

import (
	"net/http"
	"time"

	_ "person-enricher/docs"
	"person-enricher/internal/metrics"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter returns a new Gorilla Mux router with endpoints:
//
// * GET /people: GetPeople
// * GET /people/{id}: GetPersonByID
// * POST /people: CreatePerson
// * PUT /people/{id}: UpdatePerson
// * DELETE /people/{id}: DeletePerson
func NewRouter(h *Handler) *mux.Router {
	r := mux.NewRouter()

	r.Use(HttpMetricsMiddleware)

	// Swagger
	r.PathPrefix("/v1/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/v1/people", h.GetPeople).Methods(http.MethodGet)
	r.HandleFunc("/v1/people/{id}", h.GetPersonByID).Methods(http.MethodGet)
	r.HandleFunc("/v1/people", h.CreatePerson).Methods(http.MethodPost)
	r.HandleFunc("/v1/people/{id}", h.UpdatePerson).Methods(http.MethodPut)
	r.HandleFunc("/v1/people/{id}", h.DeletePerson).Methods(http.MethodDelete)

	return r
}

// httpMetricsMiddleware измеряет длительность HTTP запросов.
func HttpMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()

		route := mux.CurrentRoute(r)
		var pathTemplate string
		if route != nil {
			pathTemplate, _ = route.GetPathTemplate()
		} else {
			pathTemplate = "not_found"
		}

		metrics.HTTPRequestDuration.WithLabelValues(r.Method, pathTemplate).Observe(duration)
	})
}
