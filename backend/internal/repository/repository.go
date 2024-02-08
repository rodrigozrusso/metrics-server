package repository

import (
	"time"

	"acme.inc/analytics/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	CreateHyperTable() error
	AddMetric(state domain.Metric) error
	ListMetrics() ([]string, error)
	GetDataByMetricName(metricName string, granularity domain.Granularity, startDate time.Time, endDate time.Time) ([]domain.MetricAVGResult, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) CreateHyperTable() error {
	return r.DB.Exec("SELECT create_hypertable('metrics', by_range('timestamp'), if_not_exists => TRUE)").Error
}

func (r repository) AddMetric(state domain.Metric) error {
	return r.DB.Create(&state).Error
}

func (r repository) ListMetrics() ([]string, error) {
	// var metrics []domain.MetricsResponse
	var metrics []string
	err := r.DB.Raw("SELECT distinct name FROM metrics ORDER BY name ASC").Scan(&metrics).Error
	return metrics, err
}

// pagination?
func (r repository) GetDataByMetricName(metricName string, granularity domain.Granularity, startDate time.Time, endDate time.Time) ([]domain.MetricAVGResult, error) {
	var metrics []domain.MetricAVGResult
	var granularityStr string
	granularityStr = "1 " + string(granularity)
	// err := r.DB.Exec("SELECT * FROM metrics WHERE name = ? AND timestamp BETWEEN ? AND ?", metricName, startDate, endDate).Scan(&metrics).Error
	// select
	// 	name,
	// 	time_bucket('1 hour', timestamp) AS granularity,
	// 	avg(value)
	// from metrics m
	// GROUP BY granularity, name
	// ORDER BY name, granularity DESC;
	err := r.DB.Raw("SELECT name, time_bucket(?, timestamp) AS time_frame, avg(value) FROM metrics WHERE name = ? AND timestamp BETWEEN ? AND ? GROUP BY time_frame, name ORDER BY time_frame DESC", granularityStr, metricName, startDate, endDate).Scan(&metrics).Error

	// err := r.DB.Where("name = ? AND timestamp BETWEEN ? AND ?", metricName, startDate, endDate).Find(&metrics).Error
	return metrics, err
}
