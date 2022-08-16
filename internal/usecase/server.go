package usecase

import "fmt"

type ServerUseCase struct {
	repo ServerMetricRepo
}

func NewServerUC(r ServerMetricRepo) *ServerUseCase {
	return &ServerUseCase{repo: r}
}

func (s *ServerUseCase) SaveGauge(name string, value float64) {
	s.repo.Store("gauge", name, value)
}

func (s *ServerUseCase) SaveCounter(name string, value int) {
	s.repo.Store("counter", name, value)
}

func (s *ServerUseCase) GetValueByType(typeName string, name string) (string, error) {
	value, err := s.repo.ReadValue(typeName, name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", value), nil
}

func (s *ServerUseCase) GetAllMetrics() []string {
	return s.repo.ReadAll()
}
