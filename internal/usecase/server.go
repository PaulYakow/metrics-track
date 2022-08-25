package usecase

import (
	"fmt"
)

// Реализация сервера

type Server struct {
	repo IServerMemory
}

func NewServerUC(repo IServerMemory) *Server {
	return &Server{repo: repo}
}

func (s *Server) SaveGauge(name string, value float64) {
	s.repo.Store("gauge", name, value)
}

func (s *Server) SaveCounter(name string, value int) {
	s.repo.Store("counter", name, value)
}

func (s *Server) GetValueByType(mType string, name string) (string, error) {
	value, err := s.repo.ReadValueByType(mType, name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", value), nil
}

func (s *Server) GetAllMetrics() []string {
	result := make([]string, 0)
	gauges, counters := s.repo.ReadAll()

	for name, gauge := range gauges {
		result = append(result, fmt.Sprintf("%s = %v [%s]", name, gauge.GetValue(), gauge.GetType()))
	}

	for name, counter := range counters {
		result = append(result, fmt.Sprintf("%s = %v [%s]", name, counter.GetValue(), counter.GetType()))
	}

	return result
}

func (s *Server) SaveValueByJSON(data []byte) error {
	return s.repo.StoreByJSON(data)
}

func (s *Server) GetValueByJSON(data []byte) ([]byte, error) {
	return s.repo.ReadValueByJSON(data)
}
