package usecase

import "github.com/PaulYakow/metrics-track/internal/entity"

// Адаптеры для клиента
type (
	// todo: возвращать error во всех функциях (для обработки в контроллере)

	IClient interface {
		Poll()
		GetAll() []entity.Metric
	}

	IClientMemory interface {
		Store(map[string]*entity.Metric)
		ReadAll() []entity.Metric
	}

	IClientGather interface {
		Update() map[string]*entity.Metric
	}
)

//Адаптеры для сервера
type (
	IServer interface {
		Save(metric *entity.Metric) error
		Get(metric entity.Metric) (*entity.Metric, error)
		GetAll() ([]entity.Metric, error)
		CheckRepo() error
	}

	IServerRepo interface {
		Store(metric *entity.Metric) error
		Read(metric entity.Metric) (*entity.Metric, error)
		ReadAll() ([]entity.Metric, error)

		CheckConnection() error
	}
)

// Общие адаптеры
type (
	IHasher interface {
		ProcessBatch([]entity.Metric) []entity.Metric
		ProcessSingle(entity.Metric) entity.Metric
		ProcessPointer(*entity.Metric) *entity.Metric
		Check(*entity.Metric) error
	}
)
