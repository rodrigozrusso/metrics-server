package service

import (
	"time"

	"acme.inc/analytics/internal/domain"
	"acme.inc/analytics/internal/repository"
	"go.uber.org/zap"
)

type Service interface {
	GetDataByMetricName(metricName string, granularity domain.Granularity, startDate time.Time, endDate time.Time) ([]domain.MetricAVGResult, error)
	AddMetric(cmd domain.Metric) error
	ListMetrics() ([]domain.Metric, error)
}

type service struct {
	repository repository.Repository
}

// NewService creates a service
// As the app is simple, it acts as aggregator and projection
// If the service has more features, this should be split
func NewService(r repository.Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) AddMetric(metric domain.Metric) error {

	err := s.repository.AddMetric(metric)
	if err != nil {
		zap.L().Fatal("Error to save Metric",
			zap.String("metric", metric.Name),
			zap.Time("timestamp", metric.Timestamp),
			zap.Error(err))
	}
	return err
}

func (s *service) ListMetrics() ([]domain.Metric, error) {
	return s.repository.ListMetrics()
}

func (s *service) GetDataByMetricName(metricName string, granularity domain.Granularity, startDate time.Time, endDate time.Time) ([]domain.MetricAVGResult, error) {
	return s.repository.GetDataByMetricName(metricName, granularity, startDate, endDate)
}
