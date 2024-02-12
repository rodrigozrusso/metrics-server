package metrics

import (
	"time"

	"go.uber.org/zap"
)

//go:generate ${MOCKGEN} -destination ./mock/mock_service.go -source ./service.go -package mock
type Service interface {
	GetDataByMetricName(metricName string, granularity Granularity, startDate time.Time, endDate time.Time) (AVGMetricResponse, error)
	AddMetric(cmd Metric) error
	ListMetrics() ([]string, error)
}

type service struct {
	repository Repository
}

// NewService creates a service
// As the app is simple, it acts as aggregator and projection
// If the service has more features, this should be split
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) AddMetric(metric Metric) error {

	err := s.repository.AddMetric(metric)
	if err != nil {
		zap.L().Fatal("Error to save Metric",
			zap.String("metric", metric.Name),
			zap.Time("timestamp", metric.Timestamp),
			zap.Error(err))
		return err
	}
	return nil
}

func (s *service) ListMetrics() ([]string, error) {
	return s.repository.ListMetrics()
}

func (s *service) GetDataByMetricName(metricName string, granularity Granularity, startDate time.Time, endDate time.Time) (AVGMetricResponse, error) {
	data, err := s.repository.GetDataByMetricName(metricName, granularity, startDate, endDate)
	var avgMetricResponse AVGMetricResponse
	avgMetricResponse.MetricName = metricName
	avgMetricResponse.Granularity = string(granularity)
	avgMetricResponse.StartTime = startDate
	avgMetricResponse.EndTime = endDate
	avgMetricResponse.Data = data

	return avgMetricResponse, err
}
