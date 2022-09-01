package usecase

import "github.com/PaulYakow/metrics-track/internal/entity"

// todo: возвращать error во всех функциях (для обработки в контроллере)
// todo: оперировать типом Metric

// Адаптеры для клиента
type (
	IClient interface {
		Poll()
		UpdateRoutes() []string
		UpdateValues() [][]byte
	}

	IClientRepo interface {
		Store(map[string]*entity.Gauge, map[string]*entity.Counter)
		ReadCurrentMetrics() []string
		ReadCurrentValues() [][]byte
	}

	IClientGather interface {
		Update() (map[string]*entity.Gauge, map[string]*entity.Counter)
	}
)

//Адаптеры для сервера
type (
	IServer interface {
		Save(metric *entity.Metric) error
		Get(metric entity.Metric) (*entity.Metric, error)
		GetAll() ([]entity.Metric, error)
	}

	ISchedule interface {
		RunStoring()
		InitMetrics()
	}

	IServerMemory interface {
		Store(metric *entity.Metric) error
		Read(metric entity.Metric) (*entity.Metric, error)
		ReadAll() []entity.Metric

		InitializeMetrics([]entity.Metric)
	}

	IServerFile interface {
		SaveMetrics(metrics []entity.Metric)
		ReadMetrics() []entity.Metric
	}
)
