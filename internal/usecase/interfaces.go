package usecase

import "github.com/PaulYakow/metrics-track/internal/entity"

type Metric interface {
	Poll()
	Report() []string
}

type MetricRepo interface {
	Store(map[string]*entity.Gauge, map[string]*entity.Counter)
	GetCurrentMetrics() []string
}

type MetricGather interface {
	Update() (map[string]*entity.Gauge, map[string]*entity.Counter)
}
