package hasher

import (
	"github.com/PaulYakow/metrics-track/internal/entity"
)

type hasherImpl struct {
	key string
}

func New(key string) *hasherImpl {
	return &hasherImpl{key: key}
}

func (h *hasherImpl) ProcessBatch(metrics []entity.Metric) []entity.Metric {
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

func (h *hasherImpl) ProcessSingle(metric entity.Metric) entity.Metric {
	if h.key != "" {
		metric.SetHash(h.key)
	}

	return metric
}

func (h *hasherImpl) ProcessPointer(ptr *entity.Metric) *entity.Metric {
	if h.key != "" {
		ptr.SetHash(h.key)
	}

	return ptr
}

func (h *hasherImpl) Check(ptr *entity.Metric) error {
	if h.key != "" {
		return ptr.CheckHash(ptr.GetHash(), h.key)
	}

	return nil
}
