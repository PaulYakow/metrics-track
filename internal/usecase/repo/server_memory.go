package repo

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

// ServerMemoryRepo реализация репозитория сервера (usecase.IServerRepo). Хранение в памяти.
type ServerMemoryRepo struct {
	metrics map[string]*entity.Metric
	sync.Mutex
}

// NewServerMemory создаёт объект ServerMemoryRepo.
func NewServerMemory() *ServerMemoryRepo {
	return &ServerMemoryRepo{
		metrics: make(map[string]*entity.Metric, 30),
	}
}

func (repo *ServerMemoryRepo) Store(metric *entity.Metric) error {
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

func (repo *ServerMemoryRepo) StoreBatch(metrics []entity.Metric) error {
	for _, metric := range metrics {
		metric := metric
		repo.Store(&metric)
	}

	return nil
}

func (repo *ServerMemoryRepo) Read(ctx context.Context, metric entity.Metric) (*entity.Metric, error) {
	repo.Lock()
	local := repo.metrics
	repo.Unlock()

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		if _, ok := local[metric.ID]; !ok {
			return nil, fmt.Errorf("repo - unknown metric: %q", metric.ID)
		}
		return local[metric.ID], nil
	}
}

func (repo *ServerMemoryRepo) ReadAll(ctx context.Context) ([]entity.Metric, error) {
	repo.Lock()
	local := repo.metrics
	repo.Unlock()

	result := make([]entity.Metric, 0, len(local))

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		for _, metric := range local {
			result = append(result, *metric)
		}
		return result, nil
	}
}

var errNoConnection = errors.New("not implement to file storage")

func (repo *ServerMemoryRepo) CheckConnection() error {
	return errNoConnection
}
