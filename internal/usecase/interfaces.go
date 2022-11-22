// Package usecase содержит логику взаимодействия с сервисами/хранилищами (для клиента и сервера).
package usecase

import (
	"context"

	"github.com/PaulYakow/metrics-track/internal/entity"
)

type (
	// IClient абстракция клиента.
	IClient interface {
		// Poll - обновление метрик "сборщиком" и сохранение в хранилище.
		Poll()
		// GetAll - чтение из хранилища всех известных на данный момент метрик.
		GetAll() []entity.Metric
	}

	// IClientMemory абстракция взаимодействия с хранилищем в памяти.
	IClientMemory interface {
		// Store - сохранение переданного словаря с метриками в памяти.
		Store(map[string]*entity.Metric)
		// ReadAll - чтение всех метрик из памяти.
		ReadAll() []entity.Metric
	}

	// IClientGather абстракция "сборщика" метрик.
	IClientGather interface {
		// Update - обновление и возврат текущих значений.
		Update() map[string]*entity.Metric
	}
)

type (
	// IServer абстракция сервера.
	IServer interface {
		// Save - сохранение переданного указателя на метрику в хранилище.
		Save(metric *entity.Metric) error
		// Get - чтение метрики по переданным данным из хранилища.
		Get(ctx context.Context, metric entity.Metric) (*entity.Metric, error)
		// SaveBatch - сохранение переданного пакета метрик в хранилище.
		SaveBatch(metrics []entity.Metric) error
		// GetAll - чтение всех метрик из хранилища.
		GetAll(ctx context.Context) ([]entity.Metric, error)
		// CheckRepo - проверка связи с хранилищем.
		CheckRepo() error
	}

	// IServerRepo абстракция взаимодействия с хранилищем
	IServerRepo interface {
		// Store - сохранение метрики.
		Store(metric *entity.Metric) error
		// StoreBatch - сохранение пакета метрик.
		StoreBatch(metrics []entity.Metric) error
		// Read - чтение метрики.
		Read(ctx context.Context, metric entity.Metric) (*entity.Metric, error)
		// ReadAll - чтение всех известных на данный момент метрик.
		ReadAll(ctx context.Context) ([]entity.Metric, error)
		// CheckConnection - проверка соединения.
		CheckConnection() error
	}
)

// Общие адаптеры
type (
	// IHasher абстракция сервиса обработки хэша.
	IHasher interface {
		// ProcessBatch - обработка пакета метрик.
		ProcessBatch([]entity.Metric) []entity.Metric
		// ProcessSingle - обработка одной метрики.
		ProcessSingle(entity.Metric) entity.Metric
		// ProcessPointer - обработка указателя на метрику.
		ProcessPointer(*entity.Metric)
		// Check - проверка хэша переданной метрики.
		Check(*entity.Metric) error
	}
)
