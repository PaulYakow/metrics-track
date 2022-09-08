package usecase

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
)

// Реализация сервера

type Server struct {
	repo   IServerRepo
	hasher IHasher
}

func NewServerUC(repo IServerRepo, hasher IHasher) *Server {
	return &Server{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *Server) Save(metric *entity.Metric) error {
	if err := s.hasher.Check(metric); err != nil {
		return err
	}
	s.hasher.ProcessPointer(metric)
	return s.repo.Store(metric)
}

func (s *Server) SaveBatch(metrics []entity.Metric) error {
	return s.repo.StoreBatch(s.hasher.ProcessBatch(metrics))
}

func (s *Server) Get(metric entity.Metric) (*entity.Metric, error) {
	auxMetric, err := s.repo.Read(metric)
	if err != nil {
		return nil, err
	}

	s.hasher.ProcessPointer(auxMetric)
	return auxMetric, nil
}

func (s *Server) GetAll() ([]entity.Metric, error) {
	return s.repo.ReadAll()
}

func (s *Server) CheckRepo() error {
	return s.repo.CheckConnection()
}
