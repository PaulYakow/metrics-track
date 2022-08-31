package usecase

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
)

// Реализация сервера

type Server struct {
	repo IServerMemory
}

func NewServerUC(repo IServerMemory) *Server {
	return &Server{repo: repo}
}

func (s *Server) Save(metric entity.Metric) error {
	return s.repo.Store(metric)
}

func (s *Server) Get(metric entity.Metric) (*entity.Metric, error) {
	return s.repo.Read(metric)
}

func (s *Server) GetAll() ([]entity.Metric, error) {
	return s.repo.ReadAll(), nil
}
