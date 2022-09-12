package repo

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
	"sync"
)

type clientMemoryRepo struct {
	sync.Mutex
	metrics map[string]*entity.Metric
}

func NewClientRepo() *clientMemoryRepo {
	return &clientMemoryRepo{
		metrics: make(map[string]*entity.Metric),
	}
}

func (repo *clientMemoryRepo) Store(metrics map[string]*entity.Metric) {
	repo.Lock()
	defer repo.Unlock()

	for name, metric := range metrics {
		if _, ok := repo.metrics[name]; !ok {
			repo.metrics[name] = metric
			continue
		}
		repo.metrics[name].Update(metric)
	}
}

func (repo *clientMemoryRepo) ReadAll() []entity.Metric {
	result := make([]entity.Metric, len(repo.metrics))

	repo.Lock()
	defer repo.Unlock()

	idx := 0
	for _, metric := range repo.metrics {
		result[idx] = *metric
		idx++
	}

	return result
}
