package usecase

import (
	"context"
	"github.com/PaulYakow/metrics-track/internal/entity"
)

// Адаптеры для клиента
type (
	// todo: возвращать error во всех функциях (где необходимо, для обработки в контроллере)

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
		Get(ctx context.Context, metric entity.Metric) (*entity.Metric, error)
		SaveBatch(metrics []entity.Metric) error
		GetAll(ctx context.Context) ([]entity.Metric, error)
		CheckRepo() error
	}

	IServerRepo interface {
		Store(metric *entity.Metric) error
		StoreBatch(metrics []entity.Metric) error
		Read(ctx context.Context, metric entity.Metric) (*entity.Metric, error)
		ReadAll(ctx context.Context) ([]entity.Metric, error)

		CheckConnection() error
	}
)

// Общие адаптеры
type (
	IHasher interface {
		ProcessBatch([]entity.Metric) []entity.Metric
		ProcessSingle(entity.Metric) entity.Metric
		ProcessPointer(*entity.Metric)
		Check(*entity.Metric) error
	}
)
