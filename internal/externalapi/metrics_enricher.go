package externalapi

import (
	"context"
	"time"

	"person-enricher/internal/metrics"
)

type metricsEnricher struct {
	enricher EnrichPersonalData
}

func NewMetricsEnricher(enricher EnrichPersonalData) EnrichPersonalData {
	return &metricsEnricher{enricher: enricher}
}

func (e *metricsEnricher) GetPersonAge(ctx context.Context, name string) (int, error) {
	start := time.Now()
	defer func() {
		metrics.EnricherRequestDuration.WithLabelValues("age").Observe(time.Since(start).Seconds())
	}()
	return e.enricher.GetPersonAge(ctx, name)
}

func (e *metricsEnricher) GetPersonGender(ctx context.Context, name string) (string, error) {
	start := time.Now()
	defer func() {
		metrics.EnricherRequestDuration.WithLabelValues("gender").Observe(time.Since(start).Seconds())
	}()
	return e.enricher.GetPersonGender(ctx, name)
}

func (e *metricsEnricher) GetPersonNationality(ctx context.Context, name string) (string, error) {
	start := time.Now()
	defer func() {
		metrics.EnricherRequestDuration.WithLabelValues("nationality").Observe(time.Since(start).Seconds())
	}()
	return e.enricher.GetPersonNationality(ctx, name)
}
