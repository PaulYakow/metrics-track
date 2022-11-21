package repo

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/PaulYakow/metrics-track/internal/entity"
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

func (repo *serverMemoryRepo) StoreBatch(metrics []entity.Metric) error {
	for _, metric := range metrics {
		metric := metric
		repo.Store(&metric)
	}

	return nil
}

func (repo *serverMemoryRepo) Read(ctx context.Context, metric entity.Metric) (*entity.Metric, error) {
	repo.Lock()
	defer repo.Unlock()

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		if _, ok := repo.metrics[metric.ID]; !ok {
			return nil, fmt.Errorf("repo - unknown metric: %q", metric.ID)
		}
		return repo.metrics[metric.ID], nil
	}
}

func (repo *serverMemoryRepo) ReadAll(ctx context.Context) ([]entity.Metric, error) {
	result := make([]entity.Metric, 0, len(repo.metrics))

	repo.Lock()
	defer repo.Unlock()

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		for _, metric := range repo.metrics {
			result = append(result, *metric)
		}
		return result, nil
	}
}

var errNoConnection = errors.New("not implement to file storage")

func (repo *serverMemoryRepo) CheckConnection() error {
	return errNoConnection
}
