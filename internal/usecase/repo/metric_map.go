package repo

import (
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
)

type MetricRepo struct {
	gauges   map[string]*entity.Gauge
	counters map[string]*entity.Counter
}

func New() *MetricRepo {
	return &MetricRepo{
		gauges:   make(map[string]*entity.Gauge),
		counters: make(map[string]*entity.Counter),
	}
}

func (r *MetricRepo) Store(g map[string]*entity.Gauge, c map[string]*entity.Counter) {
	r.gauges = g
	r.counters = c
}

func (r *MetricRepo) GetCurrentMetrics() []string {
	result := make([]string, 0)

	for name, gauge := range r.gauges {
		result = append(result, fmt.Sprintf("/%s/%s/%v", gauge.GetType(), name, gauge.GetValue()))
	}

	for name, counter := range r.counters {
		result = append(result, fmt.Sprintf("/%s/%s/%v", counter.GetType(), name, counter.GetValue()))
	}

	return result
}
