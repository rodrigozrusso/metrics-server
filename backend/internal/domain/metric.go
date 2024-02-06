package domain

import (
	"time"
)

type Metric struct {
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Value     float64   `gorm:"type:float;not null" json:"value"`
}

type Granularity string

const (
	Hour Granularity = "hour"
	Day  Granularity = "day"
	Week Granularity = "week"
)

type MetricAVGResult struct {
	TimeFrame time.Time
	Name      string
	Avg       float64
}
