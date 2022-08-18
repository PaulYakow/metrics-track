package usecase

import (
	"fmt"
)

//Адаптеры для сервера

type IServer interface {
	SaveGauge(name string, value float64)
	SaveCounter(name string, value int)
	SaveValueByJson(data []byte) error
	GetValueByType(mType string, name string) (string, error)
	GetValueByJson(data []byte) ([]byte, error)
	GetAllMetrics() []string
}

type IServerRepo interface {
	Store(mType string, name string, value any) error
	StoreByJson(data []byte) error
	ReadValueByType(mType string, name string) (any, error)
	ReadValueByJson(data []byte) ([]byte, error)
	ReadAll() []string // Прочитать все известные на данный момент значения - "Имя = Значение [Тип]"
}

// Реализация сервера

type Server struct {
	repo IServerRepo
}

func NewServerUC(r IServerRepo) *Server {
	return &Server{repo: r}
}

func (s *Server) SaveGauge(name string, value float64) {
	s.repo.Store("gauge", name, value)
}

func (s *Server) SaveCounter(name string, value int) {
	s.repo.Store("counter", name, value)
}

func (s *Server) SaveValueByJson(data []byte) error {
	return s.repo.StoreByJson(data)
}

func (s *Server) GetValueByType(mType string, name string) (string, error) {
	value, err := s.repo.ReadValueByType(mType, name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", value), nil
}

func (s *Server) GetAllMetrics() []string {
	return s.repo.ReadAll()
}

func (s *Server) GetValueByJson(data []byte) ([]byte, error) {
	return s.repo.ReadValueByJson(data)
}
