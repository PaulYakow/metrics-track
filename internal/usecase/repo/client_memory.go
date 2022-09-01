package repo

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
	"sync"
)

type ClientRepo struct {
	sync.Mutex
	metrics map[string]*entity.Metric
}

func NewClientRepo() *ClientRepo {
	return &ClientRepo{
		metrics: make(map[string]*entity.Metric),
	}
}

func (r *ClientRepo) Store(metrics map[string]*entity.Metric) {
	r.Lock()
	defer r.Unlock()

	r.metrics = metrics
}

func (r *ClientRepo) ReadAll() []entity.Metric {
	result := make([]entity.Metric, 0)

	r.Lock()
	defer r.Unlock()

	for _, metric := range r.metrics {
		result = append(result, *metric)
	}

	return result
}
