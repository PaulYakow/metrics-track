package usecase

import (
	"context"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

// ServerUC реализация контроллера сервера (IServer)
type ServerUC struct {
	repo   IServerRepo
	hasher IHasher
}

// NewServerUC создаёт объект ServerUC
func NewServerUC(repo IServerRepo, hasher IHasher) *ServerUC {
	return &ServerUC{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *ServerUC) Save(metric *entity.Metric) error {
	if err := s.hasher.Check(metric); err != nil {
		return err
	}
	s.hasher.ProcessPointer(metric)
	return s.repo.Store(metric)
}

func (s *ServerUC) SaveBatch(metrics []entity.Metric) error {
	return s.repo.StoreBatch(s.hasher.ProcessBatch(metrics))
}

func (s *ServerUC) Get(ctx context.Context, metric entity.Metric) (*entity.Metric, error) {
	auxMetric, err := s.repo.Read(ctx, metric)
	if err != nil {
		return nil, err
	}

	s.hasher.ProcessPointer(auxMetric)
	return auxMetric, nil
}

func (s *ServerUC) GetAll(ctx context.Context) ([]entity.Metric, error) {
	return s.repo.ReadAll(ctx)
}

func (s *ServerUC) CheckRepo() error {
	return s.repo.CheckConnection()
}
