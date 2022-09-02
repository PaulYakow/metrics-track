package usecase

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
)

// Реализация сервера

type Server struct {
	repo   IServerMemory
	hasher IHasher
}

func NewServerUC(repo IServerMemory, hasher IHasher) *Server {
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

func (s *Server) Get(metric entity.Metric) (*entity.Metric, error) {
	auxMetric, err := s.repo.Read(metric)
	if err != nil {
		return nil, err
	}

	s.hasher.ProcessPointer(auxMetric)
	return auxMetric, nil
}

func (s *Server) GetAll() ([]entity.Metric, error) {
	return s.repo.ReadAll(), nil
}
