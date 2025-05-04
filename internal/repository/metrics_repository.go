package repository

import (
	"context"
	"time"

	"person-enricher/internal/metrics"
	"person-enricher/internal/models"
)

type metricsRepository struct {
	repo PersonRepository
}

func NewMetricsRepository(repo PersonRepository) PersonRepository {
	return &metricsRepository{repo: repo}
}

func (r *metricsRepository) Create(ctx context.Context, p models.Person) (models.Person, error) {
	start := time.Now()
	defer func() {
		metrics.RepoMethodDuration.WithLabelValues("Create").Observe(time.Since(start).Seconds())
	}()
	return r.repo.Create(ctx, p)
}

func (r *metricsRepository) List(ctx context.Context, filter models.PeopleFilter) ([]models.Person, error) {
	start := time.Now()
	defer func() {
		metrics.RepoMethodDuration.WithLabelValues("List").Observe(time.Since(start).Seconds())
	}()
	return r.repo.List(ctx, filter)
}

func (r *metricsRepository) GetByID(ctx context.Context, id string) (models.Person, error) {
	start := time.Now()
	defer func() {
		metrics.RepoMethodDuration.WithLabelValues("GetByID").Observe(time.Since(start).Seconds())
	}()
	return r.repo.GetByID(ctx, id)
}

func (r *metricsRepository) Update(ctx context.Context, p models.Person) (models.Person, error) {
	start := time.Now()
	defer func() {
		metrics.RepoMethodDuration.WithLabelValues("Update").Observe(time.Since(start).Seconds())
	}()
	return r.repo.Update(ctx, p)
}

func (r *metricsRepository) Delete(ctx context.Context, id string) error {
	start := time.Now()
	defer func() {
		metrics.RepoMethodDuration.WithLabelValues("Delete").Observe(time.Since(start).Seconds())
	}()
	return r.repo.Delete(ctx, id)
}
