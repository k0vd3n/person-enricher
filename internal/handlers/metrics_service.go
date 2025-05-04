package handlers

import (
	"context"
	"time"

	"person-enricher/internal/metrics"
	"person-enricher/internal/models"
	"person-enricher/internal/service"
)

// instrumentedService wraps a PersonService and instruments its methods
type instrumentedService struct {
	service service.PersonService
}

// NewInstrumentedService creates a new instance of instrumentedService,
// which wraps an existing PersonService and instruments its methods to
// collect execution metrics such as method duration.
func NewInstrumentedService(s service.PersonService) service.PersonService {
	return &instrumentedService{service: s}
}

// GetPeople instruments the GetPeople method of the underlying PersonService
// and adds execution time metrics to the ServiceMethodDuration metric.
func (s *instrumentedService) GetPeople(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error) {
	start := time.Now()
	defer func() {
		metrics.ServiceMethodDuration.WithLabelValues("GetPeople").Observe(time.Since(start).Seconds())
	}()
	return s.service.GetPeople(ctx, filter)
}

// GetPersonByID instruments the GetPersonByID method of the underlying PersonService
// and adds execution time metrics to the ServiceMethodDuration metric.
func (s *instrumentedService) GetPersonByID(ctx context.Context, id string) (models.Person, error) {
	start := time.Now()
	defer func() {
		metrics.ServiceMethodDuration.WithLabelValues("GetPersonByID").Observe(time.Since(start).Seconds())
	}()
	return s.service.GetPersonByID(ctx, id)
}

// CreatePerson instruments the CreatePerson method of the underlying PersonService
// and adds execution time metrics to the ServiceMethodDuration metric.
func (s *instrumentedService) CreatePerson(ctx context.Context, req models.CreatePersonRequest) (models.Person, error) {
	start := time.Now()
	defer func() {
		metrics.ServiceMethodDuration.WithLabelValues("CreatePerson").Observe(time.Since(start).Seconds())
	}()
	return s.service.CreatePerson(ctx, req)
}

// UpdatePerson instruments the UpdatePerson method of the underlying PersonService
// and adds execution time metrics to the ServiceMethodDuration metric.
func (s *instrumentedService) UpdatePerson(ctx context.Context, id string, req models.UpdatePersonRequest) (models.Person, error) {
	start := time.Now()
	defer func() {
		metrics.ServiceMethodDuration.WithLabelValues("UpdatePerson").Observe(time.Since(start).Seconds())
	}()
	return s.service.UpdatePerson(ctx, id, req)
}

// DeletePerson instruments the DeletePerson method of the underlying PersonService
// and adds execution time metrics to the ServiceMethodDuration metric.
func (s *instrumentedService) DeletePerson(ctx context.Context, id string) error {
	start := time.Now()
	defer func() {
		metrics.ServiceMethodDuration.WithLabelValues("DeletePerson").Observe(time.Since(start).Seconds())
	}()
	return s.service.DeletePerson(ctx, id)
}
