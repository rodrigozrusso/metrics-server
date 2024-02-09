package metrics

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateHyperTable() error
	AddMetric(state Metric) error
	ListMetrics() ([]string, error)
	GetDataByMetricName(metricName string, granularity Granularity, startDate time.Time, endDate time.Time) ([]MetricAVGResult, error)
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

func (r repository) AddMetric(state Metric) error {
	return r.DB.Create(&state).Error
}

func (r repository) ListMetrics() ([]string, error) {
	var metrics []string
	err := r.DB.Raw("SELECT distinct name FROM metrics ORDER BY name ASC").Scan(&metrics).Error
	return metrics, err
}

// pagination?
func (r repository) GetDataByMetricName(metricName string, granularity Granularity, startDate time.Time, endDate time.Time) ([]MetricAVGResult, error) {
	var metrics []MetricAVGResult
	var timeIntervalStr string
	// 1 minute, 1 hour, 1 day
	timeIntervalStr = "1 " + string(granularity)

	err := r.DB.Raw("SELECT name, time_bucket(?, timestamp) AS time_frame, avg(value) FROM metrics WHERE name = ? AND timestamp BETWEEN ? AND ? GROUP BY time_frame, name ORDER BY time_frame DESC", timeIntervalStr, metricName, startDate, endDate).Scan(&metrics).Error
	return metrics, err
}
