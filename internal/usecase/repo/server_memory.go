package repo

import (
	"errors"
	"fmt"
	"github.com/PaulYakow/metrics-track/internal/entity"
	"log"
	"sync"
)

var (
	errNotFound = errors.New("not found")
)

type ServerMemory struct {
	sync.Mutex
	metrics map[string]*entity.Metric
}

func NewServerMemory() *ServerMemory {
	return &ServerMemory{
		metrics: make(map[string]*entity.Metric),
	}
}

func (repo *ServerMemory) InitializeMetrics(metrics []*entity.Metric) {
	for _, metric := range metrics {
		if err := repo.Store(metric); err != nil {
			log.Printf("init metrics: %v", err)
		}
	}
}

func (repo *ServerMemory) Store(metric *entity.Metric) error {
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

func (repo *ServerMemory) Read(metric entity.Metric) (*entity.Metric, error) {
	repo.Lock()
	defer repo.Unlock()

	if _, ok := repo.metrics[metric.ID]; !ok {
		return nil, fmt.Errorf("repo - unknown metric: %q", metric.ID)
	}
	return repo.metrics[metric.ID], nil
}

func (repo *ServerMemory) ReadAll() []entity.Metric {
	result := make([]entity.Metric, 0)

	repo.Lock()
	defer repo.Unlock()

	for _, metric := range repo.metrics {
		result = append(result, *metric)
	}
	return result
}
