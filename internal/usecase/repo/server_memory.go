package repo

import (
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"sync"
)

type serverMemoryRepo struct {
	sync.Mutex
	metrics map[string]*entity.Metric
}

func NewServerMemory() *serverMemoryRepo {
	return &serverMemoryRepo{
		metrics: make(map[string]*entity.Metric),
	}
}

func (repo *serverMemoryRepo) Store(metric *entity.Metric) error {
	repo.Lock()
	defer repo.Unlock()

	if _, ok := repo.metrics[metric.ID]; !ok {
		repo.metrics[metric.ID] = metric
		return nil
	}

	if err := repo.metrics[metric.ID].Update(metric); err != nil {
		return err
	}

	return nil
}

func (repo *serverMemoryRepo) Read(metric entity.Metric) (*entity.Metric, error) {
	repo.Lock()
	defer repo.Unlock()

	if _, ok := repo.metrics[metric.ID]; !ok {
		return nil, fmt.Errorf("repo - unknown metric: %q", metric.ID)
	}
	return repo.metrics[metric.ID], nil
}

func (repo *serverMemoryRepo) ReadAll() []entity.Metric {
	result := make([]entity.Metric, 0)

	repo.Lock()
	defer repo.Unlock()

	for _, metric := range repo.metrics {
		result = append(result, *metric)
	}
	return result
}

func (repo *serverMemoryRepo) CheckConnection() error {
	return fmt.Errorf("not implement to file storage")
}
