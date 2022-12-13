// Package hasher содержит функционал обработки и проверки хэша по заданному ключу.
package hasher

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
)

// HasherImpl реализация сервиса обработки хэша (usecase.IHasher).
type HasherImpl struct {
	key string
}

// New создаёт объект HasherImpl
func New(key string) *HasherImpl {
	return &HasherImpl{key: key}
}

func (h *HasherImpl) ProcessBatch(metrics []entity.Metric) []entity.Metric {
	if h.key != "" {
		result := make([]entity.Metric, len(metrics))
		for idx, metric := range metrics {
			metric.SetHash(h.key)
			result[idx] = metric
		}
		return result
	}
	return metrics
}

func (h *HasherImpl) ProcessSingle(metric entity.Metric) entity.Metric {
	if h.key != "" {
		metric.SetHash(h.key)
	}

	return metric
}

func (h *HasherImpl) ProcessPointer(ptr *entity.Metric) {
	if h.key != "" {
		ptr.SetHash(h.key)
	}
}

func (h *HasherImpl) Check(ptr *entity.Metric) error {
	if h.key != "" {
		return ptr.CheckHash(ptr.GetHash(), h.key)
	}

	return nil
}
