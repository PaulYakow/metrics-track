package usecase

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
)

// Реализация сервера

type Server struct {
	memory IServerMemory
	repo   IServerRepo
	hasher IHasher
}

func NewServerUC(memory IServerMemory, repo IServerRepo, hasher IHasher) *Server {
	return &Server{
		memory: memory,
		repo:   repo,
		hasher: hasher,
	}
}

func (s *Server) Save(metric *entity.Metric) error {
	if err := s.hasher.Check(metric); err != nil {
		return err
	}
	s.hasher.ProcessPointer(metric)
	return s.memory.Store(metric)
}

func (s *Server) Get(metric entity.Metric) (*entity.Metric, error) {
	auxMetric, err := s.memory.Read(metric)
	if err != nil {
		return nil, err
	}

	s.hasher.ProcessPointer(auxMetric)
	return auxMetric, nil
}

func (s *Server) GetAll() ([]entity.Metric, error) {
	return s.memory.ReadAll(), nil
}

func (s *Server) CheckRepo() error {
	return s.repo.CheckConnection()
}
