// Package repo содержит реализации репозиториев для клиента и сервера (usecase.IClientMemory и usecase.IServerRepo).
package repo

import (
	"sync"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

// ClientMemoryRepo реализация репозитория клиента (usecase.IClientMemory). Хранение в памяти.
type ClientMemoryRepo struct {
	metrics map[string]*entity.Metric
	sync.Mutex
}

// NewClientRepo создаёт объект ClientMemoryRepo
func NewClientRepo() *ClientMemoryRepo {
	return &ClientMemoryRepo{
		metrics: make(map[string]*entity.Metric, 30),
	}
}

func (repo *ClientMemoryRepo) Store(metrics map[string]*entity.Metric) {
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

func (repo *ClientMemoryRepo) ReadAll() []entity.Metric {
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
