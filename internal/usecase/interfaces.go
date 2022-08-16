package usecase

import "github.com/PaulYakow/metrics-track/internal/entity"

// Адаптеры для клиента

type ClientMetric interface {
	Poll()
	Report() []string
}

type ClientMetricRepo interface {
	Store(map[string]*entity.Gauge, map[string]*entity.Counter)
	ReadCurrentMetrics() []string
}

type ClientMetricGather interface {
	Update() (map[string]*entity.Gauge, map[string]*entity.Counter)
}

//Адаптеры для сервера

type ServerMetric interface {
	SaveGauge(name string, value float64)
	SaveCounter(name string, value int)
	GetValueByType(typeName string, name string) (string, error)
	GetAllMetrics() []string
}

type ServerMetricRepo interface {
	Store(typeName string, name string, value any) error
	ReadValue(typeName string, name string) (any, error)
	ReadAll() []string // Прочитать все известные на данный момент значения - "Имя = Значение [Тип]"
}
