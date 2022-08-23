package usecase

import "github.com/PaulYakow/metrics-track/internal/entity"

// Адаптеры для клиента

type IClient interface {
	Poll()
	UpdateRoutes() []string
	UpdateValues() [][]byte
}

type IClientRepo interface {
	Store(map[string]*entity.Gauge, map[string]*entity.Counter)
	ReadCurrentMetrics() []string
	ReadCurrentValues() [][]byte
}

type IClientGather interface {
	Update() (map[string]*entity.Gauge, map[string]*entity.Counter)
}

//Адаптеры для сервера

type IServer interface {
	SaveGauge(name string, value float64)
	SaveCounter(name string, value int)
	GetValueByType(mType string, name string) (string, error)
	GetAllMetrics() []string // Вернуть все известные на данный момент метрики (в виде "Имя = Значение [Тип]\n")

	SaveValueByJSON(data []byte) error
	GetValueByJSON(data []byte) ([]byte, error)
}

type ISchedule interface {
	RunStoring()
	InitMetrics()
}

type IServerMemory interface {
	Store(mType string, name string, value any) error
	ReadValueByType(mType string, name string) (any, error)
	ReadAll() (map[string]*entity.Gauge, map[string]*entity.Counter) // Считать все известные на данный момент метрики

	StoreByJSON(data []byte) error
	ReadValueByJSON(data []byte) ([]byte, error)

	InitializeMetrics([]entity.Metrics)
}

type IServerFile interface {
	SaveMetrics(metrics []entity.Metrics)
	ReadMetrics() []entity.Metrics
}
