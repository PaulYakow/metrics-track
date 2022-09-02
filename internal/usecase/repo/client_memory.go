package repo

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
	"sync"
)

type clientRepo struct {
	sync.Mutex
	metrics map[string]*entity.Metric
}

func NewClientRepo() *clientRepo {
	return &clientRepo{
		metrics: make(map[string]*entity.Metric),
	}
}

func (r *clientRepo) Store(metrics map[string]*entity.Metric) {
	r.Lock()
	defer r.Unlock()

	r.metrics = metrics
}

func (r *clientRepo) ReadAll() []entity.Metric {
	result := make([]entity.Metric, len(r.metrics))

	r.Lock()
	defer r.Unlock()

	idx := 0
	for _, metric := range r.metrics {
		result[idx] = *metric
		idx++
	}

	return result
}
